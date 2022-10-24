package controllers

import (
	"rapidtech/shopping-rest-api/database"
	"rapidtech/shopping-rest-api/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	Db *gorm.DB
	store *session.Store
}

func (controller *AuthController) Login(c *fiber.Ctx) error {
	sess, err := controller.store.Get(c)

	if err!=nil {
		panic(err)
	}

	var myform models.User
	var data models.User

	if err := c.BodyParser(&myform); err != nil {
		return c.JSON(fiber.Map{"error": err})
	}

	username := myform.Username
	plainPassword := myform.Password

	err2 := models.ReadOneUser(controller.Db, &data, username)

	if err2 != nil {
		return c.Redirect("/login")
	}
	
	hashPassword := data.Password

	check := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(plainPassword))

	status := check == nil

	if status {
		sess.Set("username", username)
		sess.Set("id", data.Id)
		sess.Save()
		return c.JSON(fiber.Map{
			"message": "Login Success!",
			"id": data.Id,
			"username": username,
		})
	} else {
		return c.JSON(fiber.Map{
			"message": "Login error!",
		})
	}

	
}

func (controller *AuthController) Register(c *fiber.Ctx) error {
	var register models.User

		if err := c.BodyParser(&register); err != nil {
			return c.JSON(fiber.Map{
				"message": "Please fill the username & password correctly!",
			})
		}

		bytes, _ := bcrypt.GenerateFromPassword([]byte(register.Password), 8)
		sHash := string(bytes)
		
		register.Password = sHash

		err := models.Register(controller.Db, &register)

		if err != nil {
			return c.JSON(fiber.Map{
				"message": "Register failed!",
			})
		}
		
		return c.JSON(fiber.Map{
			"message": "Register success!",
			"username": register.Username,
		})
}

func (controller *AuthController) CheckSession(c *fiber.Ctx) error {
	sess,_ := controller.store.Get(c)
	user := sess.Get("username")
	id := sess.Get("id")

	return c.JSON(fiber.Map{
		"Title": "Users",
		"Id": id,
		"Username": user,
	})
}

// /logout
func (controller *AuthController) Logout(c *fiber.Ctx) error {
	sess, err := controller.store.Get(c)

	if err != nil {
		panic(err)
	}

	sess.Destroy()

	return c.JSON(fiber.Map{
		"message": "Logout success!",
	})
}

func InitAuthController(s *session.Store) *AuthController {
	db := database.InitDb()

	db.AutoMigrate(&models.User{})
	
	return &AuthController{Db: db, store: s}
}