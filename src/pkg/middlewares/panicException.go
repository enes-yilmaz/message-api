package middlewares

import (
	"MESSAGEAPI/src/pkg/errors"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

func PanicExceptionHandling() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if err := recover(); err != nil {
					//err = fmt.Sprintf("Err: %v , StackTrace: %s", err, string(debug.Stack()))
					switch v := err.(type) {
					case *errors.Error:
						go v.Log()
						c.JSON(v.StatusCode, v.Public)
					case errors.Error:
						go v.Log()
						c.JSON(v.StatusCode, v.Public)
					default:
						go logrus.Error(err)
						c.NoContent(http.StatusInternalServerError)
					}
				}
			}()
			return next(c)
		}
	}
}
