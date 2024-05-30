package main

import (
	"ngc8/config"
	"ngc8/handler"
	"ngc8/middleware"
	"ngc8/model"
	"ngc8/repo"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.Use(gin.Recovery())

	db := config.ConnectPostgre()

	Repo := &repo.PostgreRepo{DB: db}

	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.ProductDB{})
	db.AutoMigrate(&model.Store{})

	h := &handler.ProductHandler{Repo: Repo}
	userhandler := handler.UserHandler{Repo: Repo}

	r.POST("/users/register", userhandler.Register)
	r.POST("/users/login", userhandler.Login)

	product := r.Group("/")
	product.Use(middleware.Auth())
	{
		product.GET("/products", h.GetProducts)
		product.GET("/product/:id", h.GetProductById)
		product.POST("/product", h.CreateProduct)
		product.PUT("/product/:id", h.UpdateProduct)
		product.DELETE("/product/:id", h.DeleteProduct)
	}

	r.Run(":8080")
}

// // Insert example data
// exampleStores := []*model.Store{
// 	{StoreName: "Store One", StorePwd: "password123", StoreEmail: "storeone@example.com", StoreType: "Retail"},
// 	{StoreName: "Store Two", StorePwd: "password456", StoreEmail: "storetwo@example.com", StoreType: "Wholesale"},
// 	{StoreName: "Store Three", StorePwd: "password789", StoreEmail: "storethree@example.com", StoreType: "Online"},
// }
// db.Create(exampleStores)
