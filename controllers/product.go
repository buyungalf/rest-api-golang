package controllers

import (
	"fmt"
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
	if form, err := c.MultipartForm(); err == nil {
		files := form.File["image"]
		
		for _, file := range files {
			var data models.Product
			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			
			if err := c.BodyParser(&data); err != nil {
				return c.JSON(fiber.Map{
					"message": "Please fill the form correctly!",
				})
			}
			
			if err := c.SaveFile(file, fmt.Sprintf("./upload/%s", file.Filename)); err != nil {
				return err
			}

			data.Image = file.Filename
		
			err := models.CreateProduct(controller.Db, &data)
		
			if err != nil {
				return c.JSON(data)
			}

			c.JSON(data)
		}
		return c.JSON(fiber.Map{
			"message": "Product inserted successfully!",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Formdata error",
	})
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
	if form, err := c.MultipartForm(); err == nil {
		files := form.File["image"]
		
		for _, file := range files {
			var data models.Product
			
			id := c.Params("id")
			idn,_ := strconv.Atoi(id)


			err := models.ReadProductById(controller.Db, &data, idn)
			if err!=nil {
				return c.SendStatus(500) // http 500 internal server error
			}
			
			var myform models.Product

			if err := c.BodyParser(&myform); err != nil {
				return c.JSON(fiber.Map{
					"message": "Please complete the form",
				})
			}

			data.Name = myform.Name
			data.Quantity = myform.Quantity
			data.Price = myform.Price
			// save product
			

			fmt.Println(file.Filename, file.Size, file.Header["Content-Type"][0])
			
			if err := c.SaveFile(file, fmt.Sprintf("./upload/%s", file.Filename)); err != nil {
				return err
			}

			data.Image = file.Filename
		
			err = models.UpdateProduct(controller.Db, &data)
		
			if err != nil {
				return c.JSON(data)
			}

			c.JSON(data)
		}
		return c.JSON(fiber.Map{
			"message": "Product updated successfully!",
		})
	}

	return c.JSON(fiber.Map{
		"message": "error",
	})
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