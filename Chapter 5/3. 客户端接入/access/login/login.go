package main

import (
	"ProjectX/access/login/internal/logic"
	"ProjectX/access/login/internal/types"
	"ProjectX/library/log"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/api/connect", connect, middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
	}))

	// Start server
	e.Logger.Fatal(e.Start(":6900"))
}

// connect handler
func connect(c echo.Context) error {
	token := c.FormValue("token")
	platform := c.FormValue("platform")

	if token == "" || platform == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Missing parameters",
		})
	}

	log.Debug("token: %s, platform: %s", token, platform)
	l := logic.NewConnectLogic(c.Request().Context())
	resp, err := l.ConnectLogic(&types.ConnectLogicRequest{
		Token:    token,
		Platform: platform,
	})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Internal server error, " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, resp)
}
