package main

import (
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/kappa-lab/go-echo-playground/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

func main() {
	e := createEcho()
	e.Logger.Fatal(e.Start(":1323"))
}
func createEcho() *echo.Echo {
	zapConf := zap.NewDevelopmentConfig()
	zapConf.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	zap, err := zapConf.Build(
		zap.AddStacktrace(zap.ErrorLevel),
	)

	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	e.Use(logger.LoggerMiddleware(zap))

	e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		Skipper: func(c echo.Context) bool {
			logger.FromContext(c.Request().Context()).Debug("skipper")
			return c.Request().Method != "POST"
		},
		KeyLookup: "header:x-api-key",
		Validator: func(key string, c echo.Context) (bool, error) {
			logger.FromContext(c.Request().Context()).Debug("validator")
			return key == "enjoy", nil
		},
	}))

	e.POST("/users", createUser)
	e.GET("/users/:id", getUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)

	return e
}

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func createUser(c echo.Context) error {
	logger.FromContext(c.Request().Context()).Debug("createUser")
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
