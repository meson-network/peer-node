package api

import (
	"github.com/labstack/echo/v4"
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
