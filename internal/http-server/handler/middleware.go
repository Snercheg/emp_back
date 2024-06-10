package handler

import (
	"EMP_Back/internal/lib/api/response"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtxKey          = "userId"
)

func (h *Handler) userIdentify(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		response.NewErrorResponse(c, http.StatusUnauthorized, "empty authorization header")
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		response.NewErrorResponse(c, http.StatusUnauthorized, "invalid authorization header")
		return
	}

	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		response.NewErrorResponse(c, http.StatusUnauthorized, "invalid authorization header")
		return
	}
	c.Set(userCtxKey, userId)
}

func (h *Handler) getUserID(c *gin.Context) (int64, error) {
	id, ok := c.Get(userCtxKey)
	if !ok {
		response.NewErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}
	idInt, ok := id.(int64)
	if !ok {
		response.NewErrorResponse(c, http.StatusInternalServerError, "user id is invalid")
		return 0, errors.New("user id is not found")
	}

	return idInt, nil
}

func (h *Handler) checkIfUserIsAdmin(c *gin.Context) {
	id, err := h.getUserID(c)
	if err != nil {
		response.NewErrorResponse(c, http.StatusUnauthorized, "user is not found")
		return
	}

	isAdmin, err := h.services.Authorization.IsAdmin(id)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, "error checking user status")
		return
	}

	if !isAdmin {
		response.NewErrorResponse(c, http.StatusUnauthorized, "user is not allowed")
		return
	}

}
