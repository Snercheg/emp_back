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

// @Summary create module
// @Tags module
// @Accept  json
// @Produce  json
// @Param input body models.Module true "module"
// @Success 201 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/module [post]
// @Security ApiKeyAuth
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

// @Summary get all modules
// @Tags module
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllModulesResponse
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/module [get]
// @Security ApiKeyAuth
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

// @Summary get module by id
// @Tags module
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Success 200 {object} models.Module
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/module/{id} [get]
// @Security ApiKeyAuth
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

// @Summary update module
// @Tags module
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Param input body models.UpdateModuleInput true "module"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/module/{id} [put]
// @Security ApiKeyAuth
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

// @Summary delete module
// @Tags module
// @Accept  json
// @Produce  json
// @Param id path int true "id"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/module/{id} [delete]
// @Security ApiKeyAuth
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
