package handler

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/lib/api/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type getAllPlantFamiliesResponse struct {
	PlantFamilies []models.PlantFamily `json:"plant_family"`
}

func (h *Handler) GetPlantFamilies(c *gin.Context) {
	plantFamilies, err := h.services.GetPlantFamilies()
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllPlantFamiliesResponse{plantFamilies})
}

func (h *Handler) CreatePlantFamily(c *gin.Context) {
	var input models.PlantFamily
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.SavePlantFamily(input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, map[string]interface{}{"id": id})
}

func (h *Handler) GetPlantFamily(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}
	plantFamily, err := h.services.GetPlantFamily(int64(id))
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, plantFamily)

}

func (h *Handler) UpdatePlantFamily(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	var input models.PlantFamily
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.UpdatePlantFamily(int64(id), input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response.StatusSuccess)

}

func (h *Handler) DeletePlantFamily(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}
	err = h.services.DeletePlantFamily(int64(id))
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response.StatusSuccess)
}
