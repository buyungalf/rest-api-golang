package controllers

import (
	"rapidtech/shopping-rest-api/database"
	"rapidtech/shopping-rest-api/models"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductController struct {
	Db *gorm.DB
}

// GET ALL PRODUCTS

// @BasePath /api

// IndexProducts godoc
// @Summary IndexProducts example
// @Schemes
// @Description IndexProducts
// @Tags training
// @Accept json
// @Produce json
// @Success 200 {json} product
// @Router /products [get]
func (controller *ProductController) IndexProducts(c *fiber.Ctx) error {
	var products []models.Product
	err := models.ReadProducts(controller.Db, &products)

	if err != nil {
		return c.SendStatus(500)
	}

	return c.JSON(products)
}

// DETAIL PRODUCT

// @BasePath /api

// DetailProduct godoc
// @Summary DetailProduct example
// @Schemes
// @Description DetailProduct
// @Param id path int true "Product Id" minimum(1)
// @Tags training
// @Accept json
// @Produce json
// @Success 200 {json} product
// @Router /products/{id} [get]
func (controller *ProductController) DetailProduct(c *fiber.Ctx) error {
	var product models.Product
	var id, err = c.ParamsInt("id")
	
	if err != nil {
		return err
	}

	err2 := models.ReadProductById(controller.Db, &product, id)

	if err2 != nil {
		return c.SendStatus(500)
	}

	return c.JSON(product)
}

// CREATE PRODUCT

// @BasePath /api

// CreateProduct godoc
// @Summary CreateProduct example
// @Schemes
// @Description CreateProduct
// @Param product body models.Product true "product"
// @Tags training
// @Accept json
// @Produce json
// @Success 200 {json} product
// @Router /products [post]
func (controller *ProductController) CreateProduct(c *fiber.Ctx) error {
	// data := new(models.Product)
	var data models.Product

		if err := c.BodyParser(&data); err != nil {
			return c.SendStatus(400)
		}

		err := models.CreateProduct(controller.Db, &data)

		if err != nil {
			return c.SendStatus(500)
		}

		return c.JSON(data)
}

// EDIT PRODUCT

// @BasePath /api

// EditProduct godoc
// @Summary EditProduct example
// @Schemes
// @Description EditProduct
// @Param id path int true "Product Id" minimum(1)
// @Param product body models.Product true "product"
// @Tags training
// @Accept json
// @Produce json
// @Success 200 {json} product
// @Router /products/{id} [put]
func (controller *ProductController) EditProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	idn,_ := strconv.Atoi(id)

	var product models.Product
	err := models.ReadProductById(controller.Db, &product, idn)
	if err!=nil {
		return c.SendStatus(500) // http 500 internal server error
	}
	var updateProduct models.Product

	if err := c.BodyParser(&updateProduct); err != nil {
		return c.SendStatus(400)
	}

	product.Name = updateProduct.Name
	product.Quantity = updateProduct.Quantity
	product.Price = updateProduct.Price

	// save product
	models.UpdateProduct(controller.Db, &product)
	
	return c.JSON(product)
}

// DELETE PRODUCT

// @BasePath /api

// DeleteProduct godoc
// @Summary DeleteProduct example
// @Schemes
// @Description DeleteProduct
// @Param id path int true "Product Id" minimum(1)
// @Tags training
// @Accept json
// @Produce json
// @Success 200 {json} message: "Delete success!"
// @Router /products/{id} [delete]
func (controller *ProductController) DeleteProduct(c *fiber.Ctx) error {
	var product models.Product
	var id, err = c.ParamsInt("id")
	
	if err != nil {
		return err
	}

	err2 := models.DeleteProduct(controller.Db, &product, id)

	if err2 != nil {
		return c.SendStatus(500)
	}

	return c.JSON(fiber.Map{
		"message": "Delete success!",
	})
}

func InitProductController() *ProductController {
	db := database.InitDb()
	db.AutoMigrate(&models.Product{})

	return &ProductController{Db: db}
}