package handler

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/lib/api/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type getAllModulesResponse struct {
	Modules []models.Module `json:"modules"`
}

func (h *Handler) CreateModule(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		return
	}

	var input models.Module
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.SaveModule(userID, input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, map[string]interface{}{"id": id})
}

func (h *Handler) GetModules(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		return
	}
	var input models.Module
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	modules, err := h.services.GetModules(userID)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllModulesResponse{Modules: modules})
}

func (h *Handler) GetModule(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		return
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	module, err := h.services.GetModule(userID, int64(id))
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, module)
}

func (h *Handler) UpdateModule(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	var input models.UpdateModuleInput
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.UpdateModule(userID, int64(id), input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response.StatusSuccess)
}

func (h *Handler) DeleteModule(c *gin.Context) {
	userID, err := h.getUserID(c)
	if err != nil {
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	err = h.services.DeleteModule(userID, int64(id))
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response.StatusSuccess)
}
