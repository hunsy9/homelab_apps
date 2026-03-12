package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"homelab-dashboard/internal/handler"
	"homelab-dashboard/internal/k8s"
	"homelab-dashboard/internal/service"
)

func main() {
	client, err := k8s.NewClient()
	if err != nil {
		log.Fatalf("Failed to create k8s client: %v", err)
	}

	clusterSvc := service.NewClusterService(client)
	clusterHandler := handler.NewClusterHandler(clusterSvc)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
	}))

	// API routes
	api := e.Group("/api/v1")
	api.GET("/cluster/nodes", clusterHandler.GetNodes)
	api.GET("/cluster/namespaces", clusterHandler.GetNamespaces)
	api.GET("/namespaces/:ns/pods", clusterHandler.GetPods)
	api.GET("/namespaces/:ns/services", clusterHandler.GetServices)
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"status": "ok"})
	})

	// Serve frontend static files
	e.Static("/", "static")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	e.Logger.Fatal(e.Start(":" + port))
}
