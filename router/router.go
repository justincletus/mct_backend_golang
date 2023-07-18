package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/justincletus/cms/handler"
	"github.com/justincletus/cms/middelware"
)

func Setup(app *fiber.App) {
	api := app.Group("/api", logger.New())

	// user login and registration
	api.Post("/register", handler.Register)
	api.Post("/login", handler.Login)
	api.Post("/logout", handler.Logout)

	// all user operation
	api.Get("/user", handler.GetUser)
	api.Get("/user/alluser", handler.GetAllUsers)
	api.Post("/user/create_team", handler.CreateTeam)
	api.Get("/user/teams", handler.GetTeams)
	api.Get("/user/:code", handler.EmailVerify)
	api.Get("/user/:id/edit", handler.GetUserById)
	api.Post("/user/:id/edit", handler.UpdateUserById)
	api.Delete("/user/:id/delete", handler.DeleteUser)
	api.Delete("/team/:id", handler.DeleteTeam)

	// reports
	api.Get("/allreport", handler.GetAllReports)
	api.Post("/report", handler.CreateReport)
	api.Post("/add_report", middelware.Protected(), handler.AddReport)
	api.Post("/upload", middelware.Protected(), handler.UploadFile)
	api.Get("/reports/:id/report", handler.GetReportById)
	api.Get("/reports/:username", handler.GetReportByUsername)
	api.Delete("/reports/:id/delete", handler.DeleteReport)
	api.Put("/reports/:id/edit", handler.UpdateReport)
	api.Put("/reports/:id/update", handler.UpdateReportMgr)

	api.Post("/job", middelware.Protected(), handler.CreateJob)
	api.Get("/project", middelware.Protected(), handler.GetProjects)
	api.Get("/project/:id", middelware.Protected(), handler.GetProjectByID)
	api.Delete("/project/:id", middelware.Protected(), handler.DeleteProject)
	api.Put("/project/:id", middelware.Protected(), handler.UpdateProject)

}
