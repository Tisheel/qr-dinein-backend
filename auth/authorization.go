package auth

import (
	"net/http"
	"regexp"
	"strconv"
)

type Permission struct {
	Methods map[string]bool
}

type RolePermissions struct {
	Resources map[string]Permission
}

var rolePermissions = map[string]RolePermissions{
	"admin": {
		Resources: map[string]Permission{
			"restaurants": {Methods: map[string]bool{"GET": true, "PUT": true}},
			"categories":  {Methods: map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}},
			"products":    {Methods: map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}},
			"orders":      {Methods: map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}},
			"staff":       {Methods: map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}},
			"settings":    {Methods: map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}},
			"ratings":     {Methods: map[string]bool{"GET": true}},
		},
	},
	"chef": {
		Resources: map[string]Permission{
			"orders": {Methods: map[string]bool{"GET": true, "PUT": true}},
		},
	},
	"superuser": {
		Resources: map[string]Permission{
			"restaurants-global": {Methods: map[string]bool{"GET": true, "POST": true, "PUT": true}},
			"restaurants":        {Methods: map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}},
			"categories":         {Methods: map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}},
			"products":           {Methods: map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}},
			"orders":             {Methods: map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}},
			"staff":              {Methods: map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}},
			"settings":           {Methods: map[string]bool{"GET": true, "POST": true, "PUT": true, "DELETE": true}},
			"ratings":            {Methods: map[string]bool{"GET": true}},
		},
	},
}

// Resource patterns to extract resource type and restaurant ID from path
var resourcePatterns = []struct {
	pattern   *regexp.Regexp
	resource  string
	restIDIdx int // capture group index for restaurant ID (0 means no restaurant scope)
}{
	{regexp.MustCompile(`^/restaurants/(\d+)/categories(?:/\d+)?$`), "categories", 1},
	{regexp.MustCompile(`^/restaurants/(\d+)/products(?:/\d+)?$`), "products", 1},
	{regexp.MustCompile(`^/restaurants/(\d+)/orders/\d+/rating$`), "ratings", 1},
	{regexp.MustCompile(`^/restaurants/(\d+)/orders(?:/\d+)?$`), "orders", 1},
	{regexp.MustCompile(`^/restaurants/(\d+)/ratings$`), "ratings", 1},
	{regexp.MustCompile(`^/restaurants/(\d+)/staff(?:/\d+)?$`), "staff", 1},
	{regexp.MustCompile(`^/restaurants/(\d+)/settings(?:/[^/]+)?$`), "settings", 1},
	{regexp.MustCompile(`^/restaurants/(\d+)$`), "restaurants", 1},
	{regexp.MustCompile(`^/restaurants$`), "restaurants-global", 0},
	{regexp.MustCompile(`^/auth/`), "auth", 0},
	{regexp.MustCompile(`^/superuser/`), "superuser-auth", 0},
}

type ResourceInfo struct {
	Name         string
	RestaurantID int
}

func extractResourceInfo(path string) *ResourceInfo {
	for _, rp := range resourcePatterns {
		matches := rp.pattern.FindStringSubmatch(path)
		if matches != nil {
			info := &ResourceInfo{Name: rp.resource}
			if rp.restIDIdx > 0 && len(matches) > rp.restIDIdx {
				info.RestaurantID, _ = strconv.Atoi(matches[rp.restIDIdx])
			}
			return info
		}
	}
	return nil
}

func checkAuthorization(claims *Claims, method, path string) (bool, string) {
	// Extract resource info from path
	resourceInfo := extractResourceInfo(path)

	if resourceInfo == nil {
		return false, "unknown resource"
	}

	// Auth endpoints are handled separately (public or authenticated)
	if resourceInfo.Name == "auth" || resourceInfo.Name == "superuser-auth" {
		return true, ""
	}

	// Get role permissions
	rolePerm, exists := rolePermissions[claims.Role]
	if !exists {
		return false, "unknown role"
	}

	// Check if role has access to this resource
	resourcePerm, hasResource := rolePerm.Resources[resourceInfo.Name]
	if !hasResource {
		return false, "access denied: no permission for this resource"
	}

	// Check if method is allowed
	if !resourcePerm.Methods[method] {
		return false, "access denied: method not allowed"
	}

	// For restaurant-scoped resources, verify restaurant ownership (superuser bypasses this)
	if claims.Role != "superuser" && resourceInfo.RestaurantID > 0 && resourceInfo.RestaurantID != claims.RestaurantID {
		return false, "access denied: restaurant mismatch"
	}

	return true, ""
}

func (m *Middleware) HandlerWithAuth(next http.Handler) http.Handler {
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

		// Check authorization
		allowed, reason := checkAuthorization(claims, r.Method, r.URL.Path)
		if !allowed {
			http.Error(w, `{"error":"`+reason+`"}`, http.StatusForbidden)
			return
		}

		// Attach claims to context
		ctx := withClaims(r.Context(), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
