package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "hyperlocal/docs"
	"hyperlocal/internal/db/postgres"
	"hyperlocal/internal/handlers"
	"hyperlocal/internal/models"
	"hyperlocal/internal/services"
	"hyperlocal/internal/web/rest"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {

	if os.Getenv("ENV") == "" {
		if err := godotenv.Load(); err != nil {
			log.Fatalln("Error loading env file", err)
		}
	}

	v := validator.New()

	db := postgres.Connect()

	model := models.New(db)
	fmt.Println("Model layer initialized")

	service := services.New(model)
	fmt.Println("Service layer initialized")

	handler := handlers.New(service, v)
	fmt.Println("Handler layer initialized")

	r := rest.NewRouter(handler)
	fmt.Println("Routers loaded")

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))
	fmt.Println("Swagger documentation available at /swagger/index.html")

	allowedOrigins := []string{"http://localhost:3000"}
	if os.Getenv("ENV") == "prod" {
		allowedOrigins = []string{"https://hyperlocal-frontend.vercel.app"}
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
	})

	corsHandler := c.Handler(r)

	fmt.Println("Server listening on port 8080...")
	http.ListenAndServe(":8080", corsHandler)

}
