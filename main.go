package main

import (
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

	e.POST("/users", createUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	e.Use(middleware.Logger())

	e.Logger.Fatal(e.Start(":1323"))
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func createUser(c echo.Context) error {
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
