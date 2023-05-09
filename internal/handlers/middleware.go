package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

const IS_ADMIN_CTX_KEY = "admin"

func (h *Handler) AuthMW(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		tokenCookie, err := ctx.Request().Cookie(SESSION_COOKIE_NAME)

		if err != nil {
			return ctx.HTML(http.StatusUnauthorized, "auth required")
		}

		isAdmin, err := h.u.IsAdmin(ctx.Request().Context(), tokenCookie.Value)

		if err != nil {
			return ctx.HTML(http.StatusInternalServerError, err.Error())
		}

		if !isAdmin {
			return ctx.HTML(http.StatusUnauthorized, "auth required")
		}

		ctx.Set(IS_ADMIN_CTX_KEY, isAdmin)
		return next(ctx)
	}
}

func IsAdmin(ctx echo.Context) bool {
	isAdmin, ok := ctx.Get(IS_ADMIN_CTX_KEY).(bool)
	if !ok {
		return false
	}
	return isAdmin
}
