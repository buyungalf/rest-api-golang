package controllers

import (
	"rapidtech/shopping-rest-api/database"
	"rapidtech/shopping-rest-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"gorm.io/gorm"
)


type CartController struct {
	Db *gorm.DB
	store *session.Store
}

//GET
func (controller *CartController) GetCart(c *fiber.Ctx) error {
	userId := c.Query("userid")
	user_id, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(fiber.Map{
			"message": "Login first",
		})
	}

	var cart []models.Cart
	err = models.ViewCart(controller.Db, &cart, user_id)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(cart)
}

func (controller *CartController) AddtoCart(c *fiber.Ctx) error {
	userId := c.Query("userid")
	productId := c.Query("productid")

	product_id,_ := strconv.Atoi(productId)
	user_id, err := strconv.Atoi(userId)
	
	if err != nil {
		c.JSON(fiber.Map{
			"message": "Login first",
		})
	}
  
	var cart models.Cart
	var product models.Product
	var find models.Cart

	err2 := models.ReadProductById(controller.Db, &product, product_id)

	if err2 != nil {
		c.JSON(fiber.Map{
			"message": "Product Error",
			"status": 500,
		})
	}

	err4 := models.FindCart(controller.Db, &find, product_id, user_id)
	if err4 != nil {
		c.JSON(fiber.Map{
			"message": "Internal Server Error",
			"status": 500,
		})
	}

	if find.Id != 0 {
		find.Quantity = find.Quantity + 1
		find.Total = find.Total + product.Price

		models.UpdateCart(controller.Db, &find)

		return c.JSON(find)		
	} else {
		cart.ProductId = product_id
		cart.UserId = user_id
		cart.Quantity = 1
		cart.Total = float64(cart.Quantity)*product.Price

		err3 := models.AddtoCart(controller.Db, &cart)

		if err3 != nil {
			c.JSON(fiber.Map{
				"message": "Internal Server Error",
				"status": 500,
			})
		}

		return c.JSON(cart)
	}
	
}

func InitCartController(s *session.Store) *CartController {
	db := database.InitDb()

	db.AutoMigrate(&models.Cart{})

	return &CartController{Db: db, store: s}
}