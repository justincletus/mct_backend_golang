package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"time"

	"cloud.google.com/go/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/justincletus/cms/config"
	"github.com/justincletus/cms/database"
	"github.com/justincletus/cms/models"
	"github.com/justincletus/cms/utils"
)

type ClientUploader struct {
	Cl         *storage.Client
	ProjectId  string
	BucketName string
	UploadPath string
}

// type ReportStruct struct {
// 	Id string
// 	Status string
// }

var uploader *ClientUploader

func init() {
	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", "sa.json")
	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client %v", err)
	}

	env, err := config.Config()
	if err != nil {
		log.Fatalf("error occurred in getting env values %v", err)
	}

	uploader = &ClientUploader{
		Cl:         client,
		BucketName: env["gcp_bucket_name"],
		ProjectId:  env["google_project_id"],
		UploadPath: "images/",
	}
}

func (c *ClientUploader) UploadFile(file multipart.File, object string) (string, error) {
	filePath := "https://storage.cloud.google.com/"
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	wc := c.Cl.Bucket(c.BucketName).Object(c.UploadPath + object).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("io.Copy: %v", err)

	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("writer closed %v", err)
	}

	filePath += c.BucketName + "/" + c.UploadPath + object

	return filePath, nil
}

func UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{
			"message": "not able to process the file",
		})
	}

	blobFile, err := file.Open()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "not able read the file",
		})
	}

	filePath, err := uploader.UploadFile(blobFile, file.Filename)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "upload file is failed",
		})
	}

	// fmt.Println(filePath)

	return c.Status(200).JSON(fiber.Map{
		"file_name": filePath,
	})
}

func convertStringBool(flag string) bool {
	if flag == "yes" {
		return true
	} else if flag == "no" {
		return false
	} else {
		return false
	}
}

func AddReport(c *fiber.Ctx) error {

	var data map[string]interface{}

	err := c.BodyParser(&data)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "not able get the post data",
		})
	}

	fmt.Println(data)

	//fmt.Println(data)
	var job models.Job

	database.DB.Where("id=?", 1000).First(&job)
	fmt.Println(job)

	report := models.MRI_Report{
		PurchaseRequisition:     data["purchase_requisition"].(bool),
		IsQuality:               data["is_quality"].(bool),
		IsQuantity:              data["is_quantity"].(bool),
		IsDamaged:               data["is_damaged"].(bool),
		IsSampleSame:            data["is_sample_same"].(bool),
		IsAnyCertification:      data["is_any_certification"].(bool),
		IsDocument:              data["is_document"].(bool),
		IsMaterialCertification: data["is_material_certification"].(bool),
		IsMillCertification:     data["is_mill_certification"].(bool),
		IsAppliedFinish:         data["is_applied_finish"].(bool),
		IsTestResult:            data["is_test_result"].(bool),
		IsDataSheet:             data["is_data_sheet"].(bool),
		IsOther:                 data["is_other"].(bool),
		IsSpareDelivery:         data["is_spare_delivery"].(bool),
		IsMaterialComply:        data["is_material_comply"].(bool),
		Comment:                 data["comment"].(string),
		Name:                    data["name"].(string),
		Signature:               data["signature"].(string),
		JobId:                   job.Id,
		Job:                     job,
	}

	fmt.Println(report)

	txtDB := database.DB.Create(&report)
	rId := txtDB.RowsAffected
	if rId == 0 {
		fmt.Println(txtDB.Error)
		return c.Status(400).JSON(fiber.Map{
			"message": fmt.Sprintf("%v", txtDB.Error),
		})

	}

	// fmt.Println(report)

	return c.Status(201).JSON(fiber.Map{
		"data": report,
	})
}

func CreateReport(ctx *fiber.Ctx) error {

	title1 := ctx.FormValue("data")
	username := ctx.FormValue("username")

	file1, err := ctx.FormFile("file")
	if err != nil {
		fmt.Println(err)
		return ctx.Status(400).JSON(fiber.Map{
			"message": "file not received",
		})
	}

	blobFile1, err := file1.Open()
	if err != nil {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "file is not able open",
		})
	}

	filePath1, err := uploader.UploadFile(blobFile1, file1.Filename)
	if err != nil {
		return ctx.Status(500).JSON(fiber.Map{
			"message": "file upload fail",
		})
	}

	var user models.User

	database.DB.Where("username", username).First(&user)
	if user.Id == 0 {
		return ctx.Status(404).JSON(fiber.Map{
			"message": "user not found",
		})
	}

	//fmt.Println(manager.Email)

	report := models.Report{
		Title1: title1,
		File1:  filePath1,
		Uid:    user.Id,
		User:   user,
		Status: "pending",
	}

	database.DB.Create(&report)
	if report.Id == 0 {
		return ctx.Status(400).JSON(fiber.Map{
			"message": "report creation is failed",
		})
	}

	id := strconv.Itoa(int(report.Id))

	eMailStruct := utils.EmailBody{
		Id:     id,
		Status: report.Status,
	}

	//fmt.Println(user.Email)

	err = eMailStruct.SendEmail(string(user.Email), "Inspection Report", "report_create.html")
	if err != nil {
		fmt.Println(err.Error())
	}

	var team models.TeamMem
	database.DB.Where("members", user.Email).First(&team)

	if team.Id != 0 {
		//fmt.Println(team.UserId)
		var manager models.User
		database.DB.Where("id=?", team.UserId).First(&manager)

		remtHost := utils.GetRemoteHostAddress(ctx)

		eMailStructMgr := utils.EmailBody{
			Fullname: manager.Fullname,
			Email:    user.Email,
			Status:   report.Status,
			Id:       id,
			Url:      remtHost,
		}

		err = eMailStructMgr.SendEmail(string(manager.Email), "Inspection Report", "create-rpt-mgr.html")
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	return ctx.Status(201).JSON(fiber.Map{
		"message": "report saved",
	})

}

func GetAllReports(c *fiber.Ctx) error {

	var reports []models.MRI_Report

	database.DB.Order("created_at desc").Find(&reports)

	var reps []models.MRI_ReportResponse
	var rep models.MRI_ReportResponse

	for _, item := range reports {
		tmp_project := getProject(item.JobId)
		rep.MRI_Report = item
		rep.Job = tmp_project
		reps = append(reps, rep)
	}

	return c.Status(200).JSON(fiber.Map{
		"data": reps,
	})
}

func getProject(id int) models.Job {
	var project models.Job

	database.DB.Where("id=?", id).First(&project)
	return project
}

func GetReports(ctx *fiber.Ctx) error {

	var output []models.Report
	database.DB.Preload("Report").Order("created_at desc").Find(&output)

	return ctx.Status(200).JSON(fiber.Map{
		"data": output,
	})

}

func GetReportById(c *fiber.Ctx) error {
	id := c.Params("id")
	var report models.Report
	database.DB.Where("id=?", id).First(&report)
	if report.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Report not found",
		})
	}

	return c.JSON(fiber.Map{
		"data": report,
	})
}

func GetReportByUsername(c *fiber.Ctx) error {
	username := c.Params("username")
	var user models.User
	database.DB.Where("username", username).First(&user)
	var reports []models.Report

	if user.Role == "manager" {
		var team models.TeamMem
		fmt.Println(user.Id)

		database.DB.Where("user_id=?", user.Id).First(&team)
		fmt.Println(team.Members)

		var repUser models.User
		database.DB.Where("email=?", team.Members).First(&repUser)

		database.DB.Where("uid=?", repUser.Id).Find(&reports)

		return c.Status(200).JSON(fiber.Map{
			"data": reports,
		})

	}

	database.DB.Where("uid=?", user.Id).Find(&reports)

	return c.Status(200).JSON(fiber.Map{
		"data": reports,
	})
}

func DeleteReport(c *fiber.Ctx) error {
	id := c.Params("id")
	var report models.Report

	database.DB.Where("id=?", id).First(&report)
	if report.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "report not found",
		})
	}

	database.DB.Unscoped().Delete(report, id)

	return c.Status(200).JSON(fiber.Map{
		"data": "report deleted",
	})

}

func UpdateReport(c *fiber.Ctx) error {
	id := c.Params("id")
	var report models.Report
	database.DB.Where("id=?", id).First(&report)
	if report.Id == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "report not found",
		})
	}

	title1 := c.FormValue("data")
	file1, err := c.FormFile("file")
	if err != nil {
		//fmt.Println(err)
		if err.Error() == "there is no uploaded file associated with the given key" {
			report.Title1 = title1
			database.DB.Save(&report)
			return c.Status(200).JSON(fiber.Map{
				"message": "report updated",
			})
		} else {
			return c.Status(400).JSON(fiber.Map{
				"message": "file not received",
			})
		}
	}

	blobFile1, err := file1.Open()
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "file is able to open",
		})
	}

	filePath1, err := uploader.UploadFile(blobFile1, file1.Filename)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "file upload failed",
		})
	}

	report.Title1 = title1
	report.File1 = filePath1

	database.DB.Save(&report)

	return c.Status(200).JSON(fiber.Map{
		"data": report,
	})

}

func UpdateReportMgr(c *fiber.Ctx) error {
	id := c.Params("id")
	var report models.Report
	database.DB.Where("id=?", id).First(&report)
	var body map[string]string

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "not received proper request",
		})
	}

	report.Status = body["status"]
	database.DB.Save(&report)

	var user models.User

	database.DB.Where("id=?", report.Uid).First(&user)

	emailBody := utils.EmailBody{
		Id:     strconv.FormatInt(int64(report.Id), 10),
		Status: report.Status,
	}

	if report.Status == "approved" {
		err = emailBody.SendEmail(string(user.Email), "Report Approved", "report_approval.html")
		if err != nil {
			fmt.Println(err)
		}

	} else if report.Status == "reject" {
		err = emailBody.SendEmail(string(user.Email), "Report Rejected", "report_reject.html")
		if err != nil {
			fmt.Println(err)
		}
	}

	return c.Status(200).JSON(fiber.Map{
		"data": report,
	})

}
