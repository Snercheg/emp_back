package handler

import (
	"EMP_Back/internal/domain/models"
	"EMP_Back/internal/lib/api/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type getAllRecommendationsResponse struct {
	Recommendations []models.Recommendation `json:"recommendations"`
}

// @Summary Create a new recommendation
// @Description Create a new recommendation
// @Tags Recommendation
// @Accept  json
// @Produce  json
// @Param input body models.Recommendation true "Create a new recommendation"
// @Success 201 {object} response.StatusSuccess
// @Failure 400,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router api/recommendation [post]
// @Security ApiKeyAuth

func (h *Handler) CreateRecommendation(c *gin.Context) {
	var input models.Recommendation
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	id, err := h.services.SaveRecommendation(input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, map[string]interface{}{"id": id})

}

// @Summary Get a recommendation
// @Description Get a recommendation
// @Tags Recommendation
// @Accept  json
// @Produce  json
// @Param id path int true "Recommendation id"
// @Success 200 {object} models.Recommendation
// @Failure 400,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router api/recommendation/{id} [get]
// @Security ApiKeyAuth

func (h *Handler) GetRecommendation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}
	recommendation, err := h.services.GetRecommendation(int64(id))
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, recommendation)
}

// @Summary Update a recommendation
// @Description Update a recommendation
// @Tags Recommendation
// @Accept  json
// @Produce  json
// @Param id path int true "Recommendation id"
// @Param input body models.Recommendation true "Update a recommendation"
// @Success 200 {object} response.StatusSuccess
// @Failure 400,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router api/recommendation/{id} [put]
// @Security ApiKeyAuth

func (h *Handler) UpdateRecommendation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}

	var input models.Recommendation
	if err := c.BindJSON(&input); err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	err = h.services.UpdateRecommendation(int64(id), input)
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response.StatusSuccess)
}

// @Summary Delete a recommendation
// @Description Delete a recommendation
// @Tags Recommendation
// @Accept  json
// @Produce  json
// @Param id path int true "Recommendation id"
// @Success 200 {object} response.StatusSuccess
// @Failure 400,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router api/recommendation/{id} [delete]
// @Security ApiKeyAuth

func (h *Handler) DeleteRecommendation(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.NewErrorResponse(c, http.StatusBadRequest, "Invalid id parameter")
		return
	}
	err = h.services.DeleteRecommendation(int64(id))
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, response.StatusSuccess)

}

// @Summary Get all recommendations
// @Description Get all recommendations
// @Tags Recommendation
// @Accept  json
// @Produce  json
// @Success 200 {object} getAllRecommendationsResponse
// @Failure 400,404 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Failure default {object} response.ErrorResponse
// @Router api/recommendation [get]
// @Security ApiKeyAuth

func (h *Handler) GetRecommendations(c *gin.Context) {
	recommendations, err := h.services.GetRecommendations()
	if err != nil {
		response.NewErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, getAllRecommendationsResponse{recommendations})
}
