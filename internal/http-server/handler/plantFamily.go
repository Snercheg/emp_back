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

// @Summary Get all plant families
// @Description Get all plant families
// @Tags PlantFamily
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllPlantFamiliesResponse
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/plant-family [get]
// @Security ApiKeyAuth
func (h *Handler) GetPlantFamilies(c *gin.Context) {
	plantFamilies, err := h.services.GetPlantFamilies()
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllPlantFamiliesResponse{plantFamilies})
}

// @Summary Create a new plant family
// @Description Create a new plant family
// @Tags PlantFamily
// @Accept  json
// @Produce  json
// @Param input body models.PlantFamily true "Plant Family"
// @Success 201 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/plant-family [post]
// @Security ApiKeyAuth
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

// @Summary Get plant family by id
// @Description Get plant family by id
// @Tags PlantFamily
// @Accept  json
// @Produce  json
// @Param id path int true "Plant Family ID"
// @Success 200 {object} models.PlantFamily
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/plant-family/{id} [get]
// @Security ApiKeyAuth
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

// @Summary Update plant family by id
// @Description Update plant family by id
// @Tags PlantFamily
// @Accept  json
// @Produce  json
// @Param id path int true "Plant Family ID"
// @Param input body models.PlantFamily true "Plant Family"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/plant-family/{id} [put]
// @Security ApiKeyAuth
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

// @Summary Delete plant family by id
// @Description Delete plant family by id
// @Tags PlantFamily
// @Accept  json
// @Produce  json
// @Param id path int true "Plant Family ID"
// @Success 200 {object} response.Response
// @Failure 400,404 {object} response.Response
// @Failure 500 {object} response.Response
// @Failure default {object} response.Response
// @Router /api/plant-family/{id} [delete]
// @Security ApiKeyAuth
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
