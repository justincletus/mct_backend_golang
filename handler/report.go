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

func AddReport(c *fiber.Ctx) error {

	var data map[string]interface{}

	err := c.BodyParser(&data)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "not able get the post data",
		})
	}

	uId, err := GetUserId(c)
	if err != nil {
		return fiber.NewError(404, "user not found")
	}

	var job models.Job
	var job_id = 0

	if data["job_no"] != "" {
		job_no := data["job_no"].(string)
		job_id, _ = strconv.Atoi(job_no)
	}

	database.DB.Where("id=?", job_id).First(&job)

	// fmt.Println(job)

	var order models.Order
	order.Project = data["project"].(string)
	order.Description = data["description"].(string)
	order.RequisitionNo = data["requisition_no"].(string)
	order.PurchaseOrderNo = data["purchase_order_no"].(string)
	order.DeliveryNoteNo = data["delivery_note_no"].(string)
	order.DateOFDelivery = data["date_of_delivery"].(string)
	order.JobId = job.Id
	order.Job = job

	//fmt.Println(order)

	txt := database.DB.Create(&order)
	if txt.RowsAffected == 0 {
		fmt.Println(txt.Error)
		return c.Status(400).JSON(fiber.Map{
			"message": "order details not saved",
		})
	}

	report := models.Report{
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
		Status:                  "pending",
		CreatedAt:               time.Now().Format("2006-01-02 15:04:05"),
		OrderId:                 order.Id,
		Order:                   order,
		UserId:                  uId,
		ReportType:              data["report_type"].(string),
	}

	txtDB := database.DB.Create(&report)
	rId := txtDB.RowsAffected
	if rId == 0 {
		fmt.Println(txtDB.Error)
		return c.Status(400).JSON(fiber.Map{
			"message": fmt.Sprintf("%v\n", txtDB.Error),
		})
	}

	var user models.User

	database.DB.Where("id=?", uId).First(&user)

	if user.Id != 0 {
		var team models.TeamMem
		teamTxt := database.DB.Where("sub_contractor=?", user.Email).First(&team)
		teamID := teamTxt.RowsAffected

		if teamID != 0 {
			id := strconv.Itoa(report.Id)
			eMailStruct := utils.EmailBody{
				Id:     id,
				Status: report.Status,
			}

			err = eMailStruct.SendEmail(string(user.Email), "Inspection Report", "report_create.html")
			if err != nil {
				fmt.Println(err.Error())
			}

			remtHost := utils.GetRemoteHostAddress(c)

			eMailStructMgr := utils.EmailBody{
				Email:  user.Email,
				Status: report.Status,
				Id:     id,
				Url:    remtHost,
			}

			err = eMailStructMgr.SendEmail(string(team.ContractorEmail), "Inspection Report", "create-rpt-mgr.html")
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}

	return c.Status(201).JSON(fiber.Map{
		"data": report,
	})
}

func GetAllReports(c *fiber.Ctx) error {

	var reports []models.Report

	database.DB.Order("created_at desc").Find(&reports)

	var reps []models.ReportResponse
	var rep models.ReportResponse

	for _, item := range reports {
		tmp_project := getOrder(item.OrderId)
		job := getJob(tmp_project.JobId)
		rep.Report = item
		rep.Order = tmp_project
		rep.Job = job
		reps = append(reps, rep)
	}

	return c.Status(200).JSON(fiber.Map{
		"data": reps,
	})
}

func getOrder(id int) models.Order {
	var order models.Order

	database.DB.Where("id=?", id).First(&order)
	return order
}

func getJob(id int) models.Job {
	var job models.Job
	database.DB.Where("id=?", id).First(&job)

	return job
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
	var ReportResponse models.ReportResponse
	var report models.Report
	database.DB.Where("id=?", id).First(&report)
	if report.Id == 0 {

		return c.Status(404).JSON(fiber.Map{
			"message": "Report not found",
		})
	}

	var order models.Order
	database.DB.Where("id=?", report.OrderId).First(&order)

	var job models.Job
	database.DB.Where("id=?", order.JobId).First(&job)

	var clientReport models.ClientReport
	database.DB.Where("report_id=?", id).First(&clientReport)

	ReportResponse.Report = report
	ReportResponse.Order = order
	ReportResponse.Job = job
	ReportResponse.ClientReport = clientReport

	return c.JSON(fiber.Map{
		"data": ReportResponse,
	})
}

func GetReportByUsername(c *fiber.Ctx) error {
	username := c.Params("username")

	var user models.User
	database.DB.Where("username", username).First(&user)

	var reports []models.Report

	if user.Role == "contractor" {
		var team models.TeamMem
		//fmt.Println(user.Id)
		database.DB.Where("contractor_email=?", user.Email).First(&team)
		fmt.Println("--------")
		fmt.Println(team)
		fmt.Println("-------")
		var subContractor models.User
		database.DB.Where("email=?", team.SubContractor).First(&subContractor)

		database.DB.Where("user_id=?", subContractor.Id).Find(&reports)

		return c.Status(200).JSON(fiber.Map{
			"data": reports,
		})

	} else if user.Role == "client" {
		var team models.TeamMem
		database.DB.Where("client_email=?", user.Email).First(&team)
		var subContractor models.User
		database.DB.Where("email=?", team.SubContractor).First(&subContractor)
		database.DB.Where("user_id=?", subContractor.Id).Find(&reports)
		return c.Status(200).JSON(fiber.Map{
			"data": reports,
		})
	}

	database.DB.Where("user_id=?", user.Id).Find(&reports)

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
			report.Name = title1
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

	report.Name = title1
	report.Signature = filePath1

	database.DB.Save(&report)

	return c.Status(200).JSON(fiber.Map{
		"data": report,
	})

}

func UpdateReportMgr(c *fiber.Ctx) error {
	id := c.Params("id")
	var report models.Report
	database.DB.Where("id=?", id).First(&report)
	var body map[string]interface{}

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"message": "not received proper request",
		})
	}

	fmt.Println(body)

	var user models.User
	database.DB.Where("id=?", report.UserId).First(&user)

	var team models.TeamMem
	database.DB.Where("sub_contractor=?", user.Email).First(&team)

	var client models.User

	uid, _ := GetUserId(c)
	database.DB.Where("id=?", uid).First(&client)

	if client.Role == "client" {
		var clientReport models.ClientReport
		if body["is_specification"] == "yes" {
			clientReport.IsSpecification = true
		} else {
			clientReport.IsSpecification = false
		}
		tmpSignData := body["signing_date"].(string)
		signingDate, _ := time.Parse("2006-01-02", tmpSignData)
		clientReport.Name = body["name"].(string)
		clientReport.Signature = body["signature"].(string)
		clientReport.Comment = body["comment"].(string)
		clientReport.SigningDate = signingDate
		clientReport.ReportId = report.Id
		transAct := database.DB.Create(&clientReport)
		if transAct.RowsAffected == 0 {
			return fiber.NewError(500, "client form update failed")
		}

		emailBody := utils.EmailBody{
			Id:      strconv.FormatInt(int64(report.Id), 10),
			Status:  report.Status,
			Message: "This reported is accepted by client",
		}

		fmt.Println(team.ContractorEmail)

		err = emailBody.SendEmail(team.SubContractor, "Report Approved", "report_approval.html", team.ContractorEmail)
		if err != nil {
			fmt.Errorf("something went wrong in sending report %s", err.Error())
		}

	} else if client.Role == "contractor" {

		report.Status = body["status"].(string)
		database.DB.Save(&report)

		var comment models.Comment

		if body["status"] == "approved" {
			comment.ApproveComment = body["comment"].(string)
		} else if body["status"] == "reject" {
			comment.RejectComment = body["comment"].(string)
		}

		if body["comment"] != "" {
			comment.ReportId = report.Id
			database.DB.Create(&comment)
		}

		emailBody := utils.EmailBody{
			Id:     strconv.FormatInt(int64(report.Id), 10),
			Status: report.Status,
			Email:  team.ContractorEmail,
		}

		if report.Status == "approved" {
			remtHost := utils.GetRemoteHostAddress(c)
			emailBody.Url = remtHost

			err = emailBody.SendEmail(string(team.ClientEmail), "Inspection Report", "create-client-rpt.html")
			if err != nil {
				fmt.Errorf("Error %s\n", err.Error())
			}

		} else if report.Status == "rejected" {
			emailBody.Message = body["comment"].(string)
			err = emailBody.SendEmail(string(user.Email), "Report Rejected", "report_reject.html")
			if err != nil {
				fmt.Println(err)
			}

		}

	}

	return c.Status(200).JSON(fiber.Map{
		"data": report,
	})

}
