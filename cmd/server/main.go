package main

import (
	"net/http"

	"github.com/JMKobayashi/Basic-API-Server/configs"
	"github.com/JMKobayashi/Basic-API-Server/internal/entity"
	"github.com/JMKobayashi/Basic-API-Server/internal/infra/database"
	"github.com/JMKobayashi/Basic-API-Server/internal/infra/webserver/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(sqlite.Open("teste.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&entity.Product{}, &entity.User{})
	productDB := database.NewProduct(db)
	productHandler := handlers.NewProductHandler(productDB)

	userDB := database.NewUser(db)
	userHandler := handlers.NewUserHandler(userDB, configs.TokenAuth, configs.JwtExperesIn)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/product", func(r chi.Router) {
		r.Use(jwtauth.Verifier(configs.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Post("/", productHandler.CreateProduct)
		r.Get("/{id}", productHandler.GetProduct)
		r.Put("/{id}", productHandler.UpdateProduct)
		r.Delete("/{id}", productHandler.DeleteProduct)
		r.Get("/", productHandler.GetProducts)
	})

	r.Post("/user", userHandler.CreateUser)
	r.Post("/user/token", userHandler.GetJWT)
	http.ListenAndServe(":8000", r)
}
