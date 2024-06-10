package handler

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/lib/api/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

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
