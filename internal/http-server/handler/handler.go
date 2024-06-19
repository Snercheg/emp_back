package handler

import (
	"EMP_Back/internal/service"
	"github.com/gin-gonic/gin"

	//_ "emp_back/docs"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
			plantFamily.POST("/", h.CreatePlantFamily, h.checkIfUserIsAdmin)
			plantFamily.GET("/:id", h.GetPlantFamily)
			plantFamily.PUT("/:id", h.UpdatePlantFamily, h.checkIfUserIsAdmin)
			plantFamily.DELETE("/:id", h.DeletePlantFamily, h.checkIfUserIsAdmin)

		}
		recommendation := api.Group("/recommendation")
		{
			recommendation.POST("/", h.CreateRecommendation, h.checkIfUserIsAdmin)

			recommendation.GET("/:recommendation_id", h.GetRecommendation)
			recommendation.GET("/", h.GetRecommendations)
			recommendation.PUT("/:recommendation_id", h.UpdateRecommendation, h.checkIfUserIsAdmin)
			recommendation.DELETE("/:recommendation_id", h.DeleteRecommendation, h.checkIfUserIsAdmin)
		}
		module := api.Group("/module", h.userIdentify)
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
				setting.PUT("/", h.UpdateModuleSetting)
			}
		}
	}
	return router
}
