package middlewares

import (
	"TartaLette/gh"
	"github.com/hashicorp/go-memdb"
	"github.com/labstack/echo/v4"
)

func GitHub(client *gh.Client) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("ghClient", client)
			return next(c)
		}
	}
}

func DbConnection(dbConnection *memdb.MemDB) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Set("dbConnection", dbConnection)
			return next(c)
		}
	}
}
