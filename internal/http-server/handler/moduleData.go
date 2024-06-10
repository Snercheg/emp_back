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
