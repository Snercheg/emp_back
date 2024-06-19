package handler

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/lib/api/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// @Summary Get module setting
// @Description Get module setting
// @Tags Setting
// @Accept  json
// @Produce  json
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.ErrorResponse
// @Router api/module/setting [get]
// @Security ApiKeyAuth

func (h *Handler) GetModuleSetting(c *gin.Context) {
	moduleId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	var input models.Setting
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	setting, err := h.services.GetSetting(int64(moduleId))
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, setting)

}

// @Summary Create module setting
// @Description Create module setting
// @Tags Setting
// @Accept  json
// @Produce  json
// @Param input body models.Setting true "Create module setting"
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.ErrorResponse
// @Router api/module/setting [post]
// @Security ApiKeyAuth

func (h *Handler) CreateModuleSetting(c *gin.Context) {
	var input models.Setting
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.SaveSetting(int64(input.ModuleId), input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, map[string]interface{}{"id": id})
}

// @Summary Update module setting
// @Description Update module setting
// @Tags Setting
// @Accept  json
// @Produce  json
// @Param input body models.Setting true "Update module setting"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.ErrorResponse
// @Router api/module/setting [put]
// @Security ApiKeyAuth

func (h *Handler) UpdateModuleSetting(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	var input models.Setting
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.UpdateSetting(int64(id), input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response.StatusSuccess)

}
