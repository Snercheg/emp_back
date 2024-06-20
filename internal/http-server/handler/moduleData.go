package handler

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/lib/api/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type getModuleDataRequest struct {
	startTime time.Time
	endTime   time.Time
}

type getAllModuleDataResponse struct {
	ModuleData []models.ModuleData `json:"module_data"`
}

// @Summary Create module data
// @Security ApiKeyAuth
// @Description Create module data
// @Tags ModuleData
// @Accept  json
// @Produce  json
// @Param input body models.ModuleData true "Create module data"
// @Success 201 {object} response.Response
// @Router /api/module/data [post]
func (h *Handler) CreateModuleData(c *gin.Context) {
	var input models.ModuleData
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err := h.services.SaveModuleData(input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, response.StatusSuccess)

}

// @Summary Get all module data
// @Security ApiKeyAuth
// @Description Get all module data
// @Tags ModuleData
// @Accept  json
// @Produce  json
// @Param input body getModuleDataRequest true "Get all module data by date"
// @Success 200 {object} getAllModuleDataResponse
// @Router /api/module/data [get]
func (h *Handler) GetModuleData(c *gin.Context) {
	moduleId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	var input getModuleDataRequest
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	data, err := h.services.GetModuleData(int64(moduleId), input.startTime, input.endTime)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllModuleDataResponse{data})

}
