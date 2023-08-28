package main

import (
	"BanjirEWS/admin"
	"BanjirEWS/carrousel"
	"BanjirEWS/database"
	"BanjirEWS/faq"
	"BanjirEWS/feedback"
	"BanjirEWS/gallery"
	"BanjirEWS/history"
	"BanjirEWS/news"
	"BanjirEWS/report"
	"BanjirEWS/river"
	"BanjirEWS/router"
	"BanjirEWS/user"
)

func main() {
	db := database.OpenDB()
	router.Init(
		river.NewHandler(
			river.NewService(
				river.NewRepository(db.GetDB()),
			),
		),
		user.NewHandler(
			user.NewService(
				user.NewRepository(db.GetDB()),
			),
		),
		news.NewHandler(
			news.NewService(
				news.NewRepository(db.GetDB()),
			),
		),
		report.NewHandler(
			report.NewService(
				report.NewRepository(db.GetDB()),
			),
		),
		feedback.NewHandler(
			feedback.NewService(
				feedback.NewRepository(db.GetDB()),
			),
		),
		gallery.NewHandler(
			gallery.NewService(
				gallery.NewRepository(db.GetDB()),
			),
		),
		faq.NewHandler(
			faq.NewService(
				faq.NewRepository(db.GetDB()),
			),
		),
		history.NewHandler(
			history.NewService(
				history.NewRepository(db.GetDB()),
			),
		),
		admin.NewHandler(
			admin.NewService(
				admin.NewRepository(db.GetDB()),
			),
		),
		carrousel.NewHandler(
			carrousel.NewService(
				carrousel.NewRepository(db.GetDB()),
			),
		),
	)

	router.Start(":8080")

}
