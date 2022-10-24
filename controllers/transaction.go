package controllers

import (
	"rapidtech/shopping-rest-api/database"
	"rapidtech/shopping-rest-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)


type TransactionController struct {
	Db *gorm.DB
}

//GET
func (controller *TransactionController) GetTransaction(c *fiber.Ctx) error {
	
	var transaction []models.Transaction
	err := models.ViewTransaction(controller.Db, &transaction)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(transaction)
}

func (controller *TransactionController) AddtoTransaction(c *fiber.Ctx) error {
	userId := c.Query("userid")

	user_id, err := strconv.Atoi(userId)
	
	if err != nil {
		c.Redirect("/login")
	}
  
	var cart []models.Cart
	var transaction models.Transaction

	models.ViewCart(controller.Db, &cart, user_id)

	for i := 0; i < len(cart); i++ {
		transaction.CartId = int(cart[i].Id)
		transaction.UserId = user_id
		transaction.Total = cart[i].Total

		cart[i].Status = "done"

		err := models.UpdateCart(controller.Db, &cart[i])
		if err != nil {
			c.JSON(err)
		}
		
		err = models.AddtoTransaction(controller.Db, &transaction)
	
		if err != nil {
			c.JSON(err)
		}
	}

	return c.JSON(transaction)
}

func (controller *TransactionController) FinishTransaction(c *fiber.Ctx) error {
	transid := c.Query("id")

	trans_id,_ := strconv.Atoi(transid)

	var transaction models.Transaction

	err := models.FinishTransaction(controller.Db, &transaction, trans_id)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(fiber.Map{
		"message": "Selected transaction is set to finish!",
	})
}

func InitTransactionController() *TransactionController {
	db := database.InitDb()

	db.AutoMigrate(&models.Transaction{})

	return &TransactionController{Db: db}
}