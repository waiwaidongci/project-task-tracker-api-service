package router

import (
	"github.com/gin-gonic/gin"

	"todo-api/internal/handler"
	"todo-api/internal/middleware"
)

func New(projects *handler.ProjectHandler, tasks *handler.TaskHandler) *gin.Engine {
	r := gin.New()
	r.Use(middleware.RequestID(), middleware.RequestLogger(), middleware.Recovery())

	api := r.Group("/api/v1")
	{
		projectsGroup := api.Group("/projects")
		{
			projectsGroup.POST("", projects.Create)
			projectsGroup.GET("", projects.List)
			projectsGroup.GET("/stats", projects.Stats)
			projectsGroup.GET("/:id", projects.Get)
			projectsGroup.PUT("/:id", projects.Update)
			projectsGroup.DELETE("/:id", projects.Delete)
		}

		tasksGroup := api.Group("/tasks")
		{
			tasksGroup.POST("", tasks.Create)
			tasksGroup.GET("", tasks.List)
			tasksGroup.GET("/today", tasks.Today)
			tasksGroup.GET("/:id", tasks.Get)
			tasksGroup.PUT("/:id", tasks.Update)
			tasksGroup.PATCH("/:id/status", tasks.ChangeStatus)
			tasksGroup.DELETE("/:id", tasks.Delete)
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	return r
}
