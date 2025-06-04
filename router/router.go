package router

import (
	"BanjirEWS/admin"
	"BanjirEWS/carrousel"
	"BanjirEWS/gallery"
	"BanjirEWS/history"
	"BanjirEWS/news"
	"BanjirEWS/report"
	"BanjirEWS/river"
	"BanjirEWS/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var r *gin.Engine

func Init(riverHandler *river.Handler, userHandler *user.Handler, newsHandler *news.Handler, reportHandler *report.Handler, galleryHandler *gallery.Handler, historyHandler *history.Handler, adminHandler *admin.Handler, carroueslHandler *carrousel.Handler) {
	r = gin.Default()

	r.Static("/uploads", "./uploads")

	config := cors.Config{
		AllowOrigins: []string{"http://localhost:5173", "http://localhost:3000", "https://gobanjiradmin.netlify.app", "https://gobanjirclient.netlify.app"},
		AllowOriginFunc: func(origin string) bool {
			return origin == "file://" || origin == "http://localhost:5173" || origin == "http://localhost:3000" || origin == "https://gobanjirclient.netlify.app" || origin == "https://gobanjiradmin.netlify.app"
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "X-Requested-With", "Content-Length"},
		AllowCredentials: true,
	}
	r.Use(cors.New(config))

	adminRoutes := r.Group("/admin")
	adminRoutes.POST("/signup", adminHandler.CreateAdmin)
	adminRoutes.GET("", adminHandler.GetAdmin)
	adminRoutes.GET("/admin", adminHandler.ViewAdmin)
	adminRoutes.DELETE("/:id", adminHandler.RemoveAdmin)
	adminRoutes.GET("/:id", adminHandler.ViewAdminById)
	adminRoutes.POST("/login", adminHandler.Login)
	adminRoutes.POST("/logout", adminHandler.Logout)
	adminRoutes.POST("/upload", adminHandler.UploadImage)

	riverRoutes := r.Group("/river")
	riverRoutes.GET("", riverHandler.DisplayRiver)
	riverRoutes.GET("/sort", riverHandler.SortRiversHandler)
	riverRoutes.GET("/search", riverHandler.SearchHandler)
	riverRoutes.GET("/connect/:id", riverHandler.ConnectSensor)
	riverRoutes.GET("/status/:id", riverHandler.GetSensorStatus)
	riverRoutes.GET("/ws", riverHandler.SendSensorData)
	riverRoutes.GET("/sensor/:id", riverHandler.ReceiveSensorData)
	riverRoutes.POST("", riverHandler.AddRiver)
	riverRoutes.GET("/status", riverHandler.DisplayRiverByStatus)
	riverRoutes.GET("/:id", riverHandler.DisplayRiverById)
	riverRoutes.GET("/total", riverHandler.DisplayRiverCount)
	riverRoutes.PUT("/:id", riverHandler.ChangeRiver)
	riverRoutes.DELETE("/:id", riverHandler.RemoveRiver)

	userRoutes := r.Group("/user")
	userRoutes.POST("/login", userHandler.Login)
	userRoutes.POST("/signup", userHandler.RegisterUser)
	userRoutes.POST("/logout", userHandler.Logout)
	userRoutes.PUT("/:id", userHandler.ChangeProfile)
	userRoutes.GET("", userHandler.GetUser)
	userRoutes.PUT("/password", userHandler.ChangePassword)
	userRoutes.POST("/sendEmail", userHandler.SendEmail)
	userRoutes.POST("/resetPassword", userHandler.CreatePassword)
	userRoutes.GET("/userList", userHandler.ViewUser)
	userRoutes.POST("", userHandler.AddUser)
	userRoutes.DELETE("/:id", userHandler.RemoveUser)
	userRoutes.GET("/:id", userHandler.ViewUserById)

	carrouselRoutes := r.Group("/carrousel")
	carrouselRoutes.GET("/:id", carroueslHandler.DisplayCarrouselById)
	carrouselRoutes.GET("/admin/:id", carroueslHandler.DisplayCarrouselByIdAdmin)
	carrouselRoutes.POST("", carroueslHandler.AddCarrousel)
	carrouselRoutes.PUT("/:id", carroueslHandler.ChangeCarrousel)
	carrouselRoutes.DELETE("/:id", carroueslHandler.RemoveCarrousel)
	carrouselRoutes.GET("", carroueslHandler.DisplayCarrousel)

	newsRoutes := r.Group("/news")
	newsRoutes.GET("/category", newsHandler.DisplayNews)
	newsRoutes.GET("/:id", newsHandler.DisplayNewsById)
	newsRoutes.DELETE("/:id", newsHandler.RemoveNews)
	newsRoutes.POST("", newsHandler.AddNews)
	newsRoutes.PUT("/:id", newsHandler.ChangeNews)

	reportRoutes := r.Group("/report")
	reportRoutes.POST("", reportHandler.AddReport)
	reportRoutes.GET("", reportHandler.DisplayReport)
	reportRoutes.GET("/:id", reportHandler.DisplayReportById)
	reportRoutes.GET("/user/:userid", reportHandler.DisplayReportByUserId)
	reportRoutes.GET("/user/:userid/:id", reportHandler.DisplayReportByUserIdbyId)
	reportRoutes.DELETE("/:id", reportHandler.RemoveReport)

	galleryRoutes := r.Group("/gallery")
	galleryRoutes.GET("", galleryHandler.ViewGallery)
	galleryRoutes.GET("/:id", galleryHandler.ViewGalleryById)
	galleryRoutes.GET("/admin/:id", galleryHandler.ViewGalleryByIdAdmin)
	galleryRoutes.POST("", galleryHandler.AddGallery)
	galleryRoutes.PUT("/:id", galleryHandler.ChangeGallery)
	galleryRoutes.DELETE("/:id", galleryHandler.RemoveGallery)

	historyRoutes := r.Group("/history")
	historyRoutes.GET("/river/:id", historyHandler.DisplayHistoryByRiverId)
	historyRoutes.GET("/:id", historyHandler.DisplayHistoryById)
	historyRoutes.GET("/time/:id", historyHandler.DisplayHistoryByRiverIdByTime)
	historyRoutes.GET("/count/:id", historyHandler.DisplayHistoryCount)
	historyRoutes.DELETE("", historyHandler.RemoveAllHistory)
	historyRoutes.DELETE("/time", historyHandler.RemoveAllHistoryByTime)
	historyRoutes.DELETE("/river/:id", historyHandler.RemoveHistoryByRiverId)
	historyRoutes.DELETE("/:id", historyHandler.RemoveHistoryById)
	historyRoutes.DELETE("/time/:id", historyHandler.RemoveHistoryByRiverIdByTime)
}

func Start(addr string) error {
	return r.Run(addr)
}
