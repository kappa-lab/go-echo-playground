package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	//e.Use(middleware.Logger())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper: func(c echo.Context) bool {
			log.Println("skipper")
			return c.Request().Method != "POST"
		},
		KeyLookup: "header:x-api-key",
		Validator: func(key string, c echo.Context) (bool, error) {
			log.Println("validator")
			return key == "enjoy", nil
		},
	}))

	e.POST("/users", createUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	e.Logger.Fatal(e.Start(":1323"))
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func createUser(c echo.Context) error {
	log.Println("createUser")
	u := &User{}
	if err := c.Bind(u); err != nil {
		return err
	}
	u.ID = uuid.NewString()
	return c.JSON(http.StatusCreated, u)
}

func getUser(c echo.Context) error {
	id := c.Param("id")
	u := User{
		ID:    id,
		Name:  "Smith",
		Email: "smith@test.com",
	}
	return c.JSON(http.StatusOK, u)
}

func updateUser(c echo.Context) error {
	id := c.Param("id")
	u := &User{}
	if err := c.Bind(u); err != nil {
		return err
	}
	u.ID = id
	return c.JSON(http.StatusOK, u)
}

func deleteUser(c echo.Context) error {
	return c.String(http.StatusNoContent, "ok")
}
