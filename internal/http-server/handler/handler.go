package handler

import (
	"EMP_Back/internal/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.SignUp)
		auth.POST("/sign-in", h.SignIn)
	}

	api := router.Group("/api")
	{
		plantFamily := api.Group("/plant-family")
		{
			plantFamily.GET("/", h.GetPlantFamily)
			plantFamily.POST("/", h.CreatePlantFamily)
			plantFamily.GET("/:id", h.GetPlantFamily)
			plantFamily.PUT("/:id", h.UpdatePlantFamily)
			plantFamily.DELETE("/:id", h.DeletePlantFamily)

			recommendation := plantFamily.Group("/recommendation")
			{
				recommendation.POST("/", h.CreatePlantFamilyRecommendation)
				recommendation.GET("/:recommendation_id", h.GetPlantFamilyRecommendation)
				recommendation.PUT("/:recommendation_id", h.UpdatePlantFamilyRecommendation)
				recommendation.DELETE("/:recommendation_id", h.DeletePlantFamilyRecommendation)
			}
		}
		module := api.Group("/module")
		{
			module.GET("/", h.GetModules)
			module.POST("/", h.CreateModule)
			module.GET("/:id", h.GetModule)
			module.PUT("/:id", h.UpdateModule)
			module.DELETE("/:id", h.DeleteModule)

			data := module.Group("/data")
			{
				data.POST("/", h.CreateModuleData)
				data.GET("/", h.GetModuleData)
			}
			setting := module.Group("/setting")
			{
				setting.GET("/", h.GetModuleSetting)
				setting.POST("/", h.CreateModuleSetting)
				setting.PUT("/:setting_id", h.UpdateModuleSetting)
			}
		}
	}
	return router
}
