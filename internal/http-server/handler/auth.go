package handler

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/lib/api/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// @Summary SignUp
// @Description create user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body models.User true "User"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/signup [post]

func (h *Handler) SignUp(c *gin.Context) {
	var input models.User
	if err := c.Bind(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.Authorization.SaveUser(input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

type signInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary SignIn
// @Description login user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body signInInput true "User"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /auth/signin [post]

func (h *Handler) SignIn(c *gin.Context) {
	var input = signInInput{}
	if err := c.Bind(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	token, err := h.services.Authorization.GenerateToken(input.Email, input.Password)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
