package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/justincletus/mct_backend_golang/handler"
	"github.com/justincletus/mct_backend_golang/middelware"
)

func Setup(app *fiber.App) {
	api := app.Group("/api", logger.New())

	// user login and registration
	api.Post("/register", handler.Register)
	api.Post("/login", handler.Login)
	api.Post("/logout", handler.Logout)

	// all user operation
	api.Get("/user", handler.GetUser)
	api.Get("/user/alluser", middelware.Protected(), handler.GetAllUsers)
	api.Post("/user/create_team", middelware.Protected(), handler.CreateTeam)
	api.Get("/user/teams", middelware.Protected(), handler.GetTeams)
	api.Get("/teams/:id", middelware.Protected(), handler.GetTeamById)
	api.Put("/teams/:id", middelware.Protected(), handler.UpdateTeam)
	api.Get("/user/:code", handler.EmailVerify)
	api.Get("/user/:id/edit", middelware.Protected(), handler.GetUserById)
	api.Post("/user/:id/edit", middelware.Protected(), handler.UpdateUserById)
	api.Delete("/user/:id/delete", middelware.Protected(), handler.DeleteUser)
	api.Delete("/team/:id", middelware.Protected(), handler.DeleteTeam)

	// reports
	api.Get("/allreport", middelware.Protected(), handler.GetAllReports)
	// api.Post("/report", handler.CreateReport)
	api.Post("/add_report", middelware.Protected(), handler.AddReport)
	api.Post("/upload", middelware.Protected(), handler.UploadFile)
	api.Get("/reports/:id/report", middelware.Protected(), handler.GetReportById)
	api.Get("/reports/:username", middelware.Protected(), handler.GetReportByUsername)
	api.Delete("/reports/:id/delete", middelware.Protected(), handler.DeleteReport)
	api.Put("/reports/:id/edit", middelware.Protected(), handler.UpdateReport)
	api.Put("/reports/:id/update", middelware.Protected(), handler.UpdateReportMgr)

	api.Get("/client_reports", middelware.Protected(), handler.GetClientReports)
	api.Delete("/client_reports/:id", middelware.Protected(), handler.DeleteClientReport)
	api.Get("/comments", middelware.Protected(), handler.GetComments)
	api.Delete("/comments/:id", middelware.Protected(), handler.DeleteComment)

	api.Post("/job", middelware.Protected(), handler.CreateJob)
	api.Get("/job", middelware.Protected(), handler.GetJobs)
	api.Get("/job", middelware.Protected(), handler.GetProjects)
	api.Get("/job/:id", middelware.Protected(), handler.GetProjectByID)
	api.Delete("/job/:id", middelware.Protected(), handler.DeleteProject)
	api.Put("/job/:id", middelware.Protected(), handler.UpdateProject)

}
