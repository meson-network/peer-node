package api

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/meson-network/peer-node/src/remote/client"
	"github.com/meson-network/peer-node/tools/http"
)

func MidToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Set("token", http.GetBearToken(c.Request().Header))
		//continue
		if err := next(c); err != nil {
			c.Error(err)
		}
		return nil
	}
}

func CheckToken(c echo.Context) error {
	v := c.Get("token")
	if v == nil {
		return errors.New("token not exist")
	}
	token := v.(string)
	if token != client.Token {
		return errors.New("token error, no auth")
	}
	return nil
}
