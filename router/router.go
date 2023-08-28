package router

import (
	"BanjirEWS/admin"
	"BanjirEWS/carrousel"
	"BanjirEWS/faq"
	"BanjirEWS/feedback"
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

func Init(riverHandler *river.Handler, userHandler *user.Handler, newsHandler *news.Handler, reportHandler *report.Handler, feedbackHandler *feedback.Handler, galleryHandler *gallery.Handler, faqHandler *faq.Handler, historyHandler *history.Handler, adminHandler *admin.Handler, carroueslHandler *carrousel.Handler) {
	r = gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "X-Requested-With", "Content-Length"}
	config.AllowCredentials = true
	r.Use(cors.New(config))

	adminRoutes := r.Group("/admin")
	adminRoutes.GET("", adminHandler.GetAdmin)
	adminRoutes.GET("/:id", adminHandler.ViewUserById)
	adminRoutes.POST("/login", adminHandler.Login)
	adminRoutes.POST("/logout", adminHandler.Logout)
	adminRoutes.GET("/user", adminHandler.ViewUser)
	adminRoutes.POST("/user", adminHandler.AddUser)
	adminRoutes.DELETE("/user/:id", adminHandler.RemoveUser)
	adminRoutes.PUT("/user/:id", adminHandler.ChangeUser)
	adminRoutes.POST("/upload", adminHandler.UploadImage)

	riverRoutes := r.Group("/river")
	riverRoutes.GET("", riverHandler.DisplayRiver)
	riverRoutes.GET("/sort", riverHandler.SortRiversHandler)
	riverRoutes.GET("/search", riverHandler.SearchHandler)
	riverRoutes.GET("/update", riverHandler.UpdateRiver)
	riverRoutes.GET("/ws", riverHandler.WebSocket)
	riverRoutes.POST("", riverHandler.AddRiver)
	riverRoutes.GET("/status", riverHandler.DisplayRiverByStatus)
	riverRoutes.GET("/:id", riverHandler.DisplayRiverById)
	riverRoutes.GET("/total", riverHandler.DisplayRiverCount)
	riverRoutes.PUT("/:id", riverHandler.ChangeRiver)
	riverRoutes.DELETE("/:id", riverHandler.RemoveRiver)

	userRoutes := r.Group("/user")
	userRoutes.POST("/login", userHandler.Login)
	userRoutes.POST("/signup", userHandler.CreateUser)
	userRoutes.POST("/logout", userHandler.Logout)
	userRoutes.PUT("/:id", userHandler.ChangeProfile)
	userRoutes.GET("", userHandler.GetUser)
	userRoutes.PUT("/password", userHandler.ChangePassword)

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

	feedbackRoutes := r.Group("/feedback")
	feedbackRoutes.POST("", feedbackHandler.AddFeedback)
	feedbackRoutes.GET("", feedbackHandler.DisplayFeedback)
	feedbackRoutes.GET("/:id", feedbackHandler.DisplayFeedbackById)
	feedbackRoutes.GET("/report/:id", feedbackHandler.DisplayFeedbackByReportId)

	galleryRoutes := r.Group("/gallery")
	galleryRoutes.GET("", galleryHandler.ViewGallery)
	galleryRoutes.GET("/:id", galleryHandler.ViewGalleryById)
	galleryRoutes.GET("/admin/:id", galleryHandler.ViewGalleryByIdAdmin)
	galleryRoutes.POST("", galleryHandler.AddGallery)
	galleryRoutes.PUT("/:id", galleryHandler.ChangeGallery)
	galleryRoutes.DELETE("/:id", galleryHandler.RemoveGallery)

	faqRoutes := r.Group("/faq")
	faqRoutes.GET("", faqHandler.ViewFaq)
	faqRoutes.POST("", faqHandler.AddFaq)
	faqRoutes.PATCH("updateFaq/:id", faqHandler.ChangeFaq)
	faqRoutes.DELETE("/:id", faqHandler.RemvoeFaq)
	faqRoutes.GET("/search", faqHandler.SearchHandler)
	faqRoutes.GET("/category", faqHandler.ViewCategory)
	faqRoutes.GET("/qa", faqHandler.ViewQa)

	historyRoutes := r.Group("/history")
	historyRoutes.GET("/river/:id", historyHandler.DisplayHistoryByRiverId)
	historyRoutes.GET("/:id", historyHandler.DisplayHistoryById)
	historyRoutes.GET("/time/:id", historyHandler.DisplayHistoryByRiverIdByTime)
	historyRoutes.GET("/count", historyHandler.DisplayHistoryCount)
	historyRoutes.DELETE("", historyHandler.RemoveAllHistory)
	historyRoutes.DELETE("time", historyHandler.RemoveAllHistoryByTime)
	historyRoutes.DELETE("/river/:id", historyHandler.RemoveHistoryByRiverId)
	historyRoutes.DELETE("/:id", historyHandler.RemoveHistoryById)
	historyRoutes.DELETE("time/:id", historyHandler.RemoveHistoryByRiverIdByTime)
}

func Start(addr string) error {
	return r.Run(addr)
}
