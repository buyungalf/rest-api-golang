package main

import (
	"rapidtech/shopping-rest-api/controllers"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
)

func main() {
	store := session.New()
	app := fiber.New()

	productAPI := controllers.InitProductController()
	auth := controllers.InitAuthController(store)
	cart := controllers.InitCartController()
	transaction := controllers.InitTransactionController()

	p := app.Group("/api")
	p.Get("/products", productAPI.IndexProducts)
	p.Get("/products/:id", productAPI.DetailProduct)
	p.Post("/products", productAPI.CreateProduct)
	p.Put("/products/:id", productAPI.EditProduct)
	p.Delete("/products/:id", productAPI.DeleteProduct)

	app.Post("/login", auth.Login)
	app.Post("/register", auth.Register)

	app.Get("/cart", cart.GetCart)
	app.Post("/cart", cart.AddtoCart)

	app.Get("/transaction", transaction.GetTransaction)
	app.Post("/transaction", transaction.AddtoTransaction)
	app.Get("/trarnsaction", transaction.FinishTransaction)

	app.Listen(":3000")
}