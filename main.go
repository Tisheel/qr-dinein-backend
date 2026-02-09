package main

import (
	"log"
	"os"
	"qr-dinein-backend/auth"
	"qr-dinein-backend/handler"
	"qr-dinein-backend/migrations"
	"qr-dinein-backend/service"
	"qr-dinein-backend/store"
	"qr-dinein-backend/strategy"

	"gofr.dev/pkg/gofr"
)

func main() {
	app := gofr.New()

	// Run migrations
	app.Migrate(migrations.All())

	// Initialize JWT manager
	jwtManager, err := auth.NewJWTManager()
	if err != nil {
		log.Fatalf("Failed to initialize JWT manager: %v", err)
	}

	// Initialize auth middleware
	authMiddleware := auth.NewMiddleware(jwtManager)

	// Configure public paths (no auth required)
	authMiddleware.AddPublicPath("POST", "/auth/login")
	authMiddleware.AddPublicPath("GET", "app/restaurants/{id}")
	authMiddleware.AddPublicPath("GET", "/restaurants/{restaurantId}/settings")
	authMiddleware.AddPublicPath("GET", "/restaurants/slug/{slug}")
	authMiddleware.AddPublicPath("GET", "/restaurants/{restaurantId}/categories")
	authMiddleware.AddPublicPath("GET", "/restaurants/{restaurantId}/products")
	authMiddleware.AddPublicPath("POST", "/restaurants/{restaurantId}/orders")
	authMiddleware.AddPublicPath("GET", "/restaurants/{restaurantId}/orders/{id}")
	authMiddleware.AddPublicPath("POST", "/customer/send-otp")
	authMiddleware.AddPublicPath("POST", "/customer/verify-otp")
	authMiddleware.AddPublicPath("GET", "/customer/session")
	authMiddleware.AddPublicPath("GET", "/restaurants/{restaurantId}/customer/orders")
	authMiddleware.AddPublicPath("POST", "/restaurants/{restaurantId}/orders/{orderId}/rating")
	authMiddleware.AddPublicPath("GET", "/restaurants/{restaurantId}/orders/{orderId}/rating")
	authMiddleware.AddPublicPath("POST", "/superuser/login")

	// Apply auth middleware (with authorization)
	app.UseMiddleware(authMiddleware.HandlerWithAuth)

	// --- Store layer ---
	restaurantStore := store.NewRestaurant()
	categoryStore := store.NewCategory()
	productStore := store.NewProduct()
	orderStore := store.NewOrder()
	staffStore := store.NewStaff()
	settingsStore := store.NewSettings()
	ratingStore := store.NewRating()

	// --- Service layer ---
	restaurantSvc := service.NewRestaurant(restaurantStore)
	categorySvc := service.NewCategory(categoryStore)
	productSvc := service.NewProduct(productStore)
	staffSvc := service.NewStaff(staffStore)
	settingsSvc := service.NewSettings(settingsStore)
	superuserUsername := os.Getenv("SUPERUSER_USERNAME")
	superuserPassword := os.Getenv("SUPERUSER_PASSWORD")
	authSvc := service.NewAuth(staffStore, jwtManager, superuserUsername, superuserPassword)
	smsSvc := service.NewSMSService()
	customerSvc := service.NewCustomer(smsSvc)
	chefResolver := strategy.NewResolver(settingsStore, staffStore, orderStore)
	orderSvc := service.NewOrder(orderStore, productStore, settingsStore, customerSvc, chefResolver)
	ratingSvc := service.NewRating(ratingStore, orderStore)

	// --- Handler layer ---
	restaurantH := handler.NewRestaurant(restaurantSvc)
	categoryH := handler.NewCategory(categorySvc)
	productH := handler.NewProduct(productSvc)
	orderH := handler.NewOrder(orderSvc)
	staffH := handler.NewStaff(staffSvc)
	settingsH := handler.NewSettings(settingsSvc)
	authH := handler.NewAuth(authSvc)
	customerH := handler.NewCustomer(customerSvc)
	ratingH := handler.NewRating(ratingSvc)

	// ==================== Routes ====================

	// --- Auth ---
	app.POST("/auth/login", authH.Login)
	app.GET("/auth/me", authH.Me)

	// --- Superuser ---
	app.POST("/superuser/login", authH.SuperuserLogin)

	// --- Restaurants ---
	app.GET("/restaurants", restaurantH.GetAll)
	app.POST("/restaurants", restaurantH.Create)
	app.GET("/restaurants/{id}", restaurantH.GetByID)
	app.GET("/restaurants/slug/{slug}", restaurantH.GetBySlug)
	app.PUT("/restaurants/{id}", restaurantH.Update)
	app.DELETE("/restaurants/{id}", restaurantH.Delete)

	// --- Categories (scoped to restaurant) ---
	app.GET("/restaurants/{restaurantId}/categories", categoryH.GetAll)
	app.POST("/restaurants/{restaurantId}/categories", categoryH.Create)
	app.GET("/restaurants/{restaurantId}/categories/{id}", categoryH.GetByID)
	app.PUT("/restaurants/{restaurantId}/categories/{id}", categoryH.Update)
	app.DELETE("/restaurants/{restaurantId}/categories/{id}", categoryH.Delete)

	// --- Products (scoped to restaurant) ---
	app.GET("/restaurants/{restaurantId}/products", productH.GetAll)
	app.POST("/restaurants/{restaurantId}/products", productH.Create)
	app.GET("/restaurants/{restaurantId}/products/{id}", productH.GetByID)
	app.PUT("/restaurants/{restaurantId}/products/{id}", productH.Update)
	app.DELETE("/restaurants/{restaurantId}/products/{id}", productH.Delete)

	// --- Orders (scoped to restaurant) ---
	app.GET("/restaurants/{restaurantId}/orders", orderH.GetAll)
	app.POST("/restaurants/{restaurantId}/orders", orderH.Create)
	app.GET("/restaurants/{restaurantId}/orders/{id}", orderH.GetByID)
	app.PUT("/restaurants/{restaurantId}/orders/{id}", orderH.Update)
	app.DELETE("/restaurants/{restaurantId}/orders/{id}", orderH.Delete)

	// --- Customer Orders (public, filtered by phone) ---
	app.GET("/restaurants/{restaurantId}/customer/orders", orderH.GetByPhone)

	// --- Order Ratings (scoped to restaurant + order) ---
	app.POST("/restaurants/{restaurantId}/orders/{orderId}/rating", ratingH.Create)
	app.GET("/restaurants/{restaurantId}/orders/{orderId}/rating", ratingH.GetByOrderID)
	app.GET("/restaurants/{restaurantId}/ratings", ratingH.GetAllByRestaurant)

	// --- Staff (scoped to restaurant) ---
	app.GET("/restaurants/{restaurantId}/staff", staffH.GetAll)
	app.POST("/restaurants/{restaurantId}/staff", staffH.Create)
	app.GET("/restaurants/{restaurantId}/staff/{id}", staffH.GetByID)
	app.PUT("/restaurants/{restaurantId}/staff/{id}", staffH.Update)
	app.DELETE("/restaurants/{restaurantId}/staff/{id}", staffH.Delete)

	// --- Settings (scoped to restaurant) ---
	app.GET("/restaurants/{restaurantId}/settings", settingsH.GetAll)
	app.PUT("/restaurants/{restaurantId}/settings", settingsH.BulkUpsert)
	app.GET("/restaurants/{restaurantId}/settings/{key}", settingsH.GetByKey)
	app.PUT("/restaurants/{restaurantId}/settings/{key}", settingsH.Upsert)
	app.DELETE("/restaurants/{restaurantId}/settings/{key}", settingsH.Delete)

	// --- Customer OTP (public endpoints) ---
	app.POST("/customer/send-otp", customerH.SendOTP)
	app.POST("/customer/verify-otp", customerH.VerifyOTP)
	app.GET("/customer/session", customerH.GetSession)

	app.Run()
}
