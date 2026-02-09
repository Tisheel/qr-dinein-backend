package auth

import (
	"context"
	"net/http"
	"strings"
)

type contextKey string

const ClaimsContextKey contextKey = "claims"

type Middleware struct {
	jwtManager  *JWTManager
	publicPaths map[string]map[string]bool // path -> method -> bool
}

func NewMiddleware(jwtManager *JWTManager) *Middleware {
	return &Middleware{
		jwtManager:  jwtManager,
		publicPaths: make(map[string]map[string]bool),
	}
}

func (m *Middleware) AddPublicPath(method, path string) {
	if m.publicPaths[path] == nil {
		m.publicPaths[path] = make(map[string]bool)
	}
	m.publicPaths[path][method] = true
}

func (m *Middleware) isPublicPath(method, path string) bool {
	// Check exact match first
	if methods, ok := m.publicPaths[path]; ok {
		if methods[method] || methods["*"] {
			return true
		}
	}

	// Check pattern matches for paths with IDs
	for pattern, methods := range m.publicPaths {
		if matchPath(pattern, path) && (methods[method] || methods["*"]) {
			return true
		}
	}

	return false
}

// matchPath checks if a path matches a pattern with {param} placeholders
func matchPath(pattern, path string) bool {
	patternParts := strings.Split(pattern, "/")
	pathParts := strings.Split(path, "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i, part := range patternParts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			continue // wildcard match
		}
		if part != pathParts[i] {
			return false
		}
	}

	return true
}

// Helper functions used by both middleware handlers
func splitAuthHeader(header string) []string {
	parts := strings.SplitN(header, " ", 2)
	if len(parts) == 2 {
		parts[0] = strings.ToLower(parts[0])
	}
	return parts
}

func containsExpired(s string) bool {
	return strings.Contains(s, "expired")
}

func withClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, ClaimsContextKey, claims)
}

// Handler provides authentication only (no authorization checks)
func (m *Middleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip auth for public paths
		if m.isPublicPath(r.Method, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error":"missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		parts := splitAuthHeader(authHeader)
		if len(parts) != 2 || parts[0] != "bearer" {
			http.Error(w, `{"error":"invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Validate token
		claims, err := m.jwtManager.ValidateToken(tokenString)
		if err != nil {
			if containsExpired(err.Error()) {
				http.Error(w, `{"error":"token expired"}`, http.StatusUnauthorized)
				return
			}
			http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
			return
		}

		// Attach claims to context
		ctx := withClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetClaimsFromContext retrieves the JWT claims from the request context
func GetClaimsFromContext(ctx context.Context) *Claims {
	claims, ok := ctx.Value(ClaimsContextKey).(*Claims)
	if !ok {
		return nil
	}
	return claims
}
