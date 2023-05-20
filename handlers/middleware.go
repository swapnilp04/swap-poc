package handlers

import (
	"crypto/subtle"
	"fmt"
	"net/http"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"encoding/json"
)

type CustomContext struct {
	echo.Context
	session *models.Session
}

func IsLoggedIn(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		fmt.Println("IsLoggedIn...")
		if c.Request().Header.Get("token") == "" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"errors": swapErr.ErrForbidden.Error()})
		}

		session := &models.Session{
			ID: c.Request().Header.Get("token"),
		}

		if err := session.Load(); err != nil {
			fmt.Println("session.Load()", err)
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"errors": swapErr.ErrForbidden.Error()})
		}

		if ok := session.Valid(); !ok {
			fmt.Println("session.Valid()", "expired")
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{"errors": swapErr.ErrForbidden.Error()})
		}
		
		cc := CustomContext{c, session}
		return next(cc)
	}
}

func OnlySwapnil() echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		// Be careful to use constant time comparison to prevent timing attacks
		if subtle.ConstantTimeCompare([]byte(username), []byte("swapnil")) == 1 &&
			subtle.ConstantTimeCompare([]byte(password), []byte("mastahey")) == 1 {
			return true, nil
		}
		return false, nil
	})
}

func MarshalFormError(err error) map[string][]string {
	val, _ := json.Marshal(err)
  m := make(map[string][]string)
  json.Unmarshal(val, &m)
  return m
}

