package models

import "time"

type Status string

const (
	StatusPending  Status = "pending"
	StatusApproved Status = "approved"
	StatusRejected Status = "rejected"
	StatusInfo            = "info"
)

type Report struct {
	Id                      int    `db:"id" json:"id"`
	PurchaseRequisition     bool   `db:"purchase_requisition" json:"purchase_requisition"`
	IsQuality               bool   `db:"is_quality" json:"is_quality"`
	IsQuantity              bool   `db:"is_quantity" json:"is_quantity"`
	IsDamaged               bool   `db:"is_damaged" json:"is_damaged"`
	IsSampleSame            bool   `db:"is_sample_same" json:"is_sample_same"`
	IsAnyCertification      bool   `db:"is_any_certification;omitempty" json:"is_any_certification"`
	IsDocument              bool   `db:"is_document;omitempty" json:"is_document"`
	IsMaterialCertification bool   `db:"is_material_certification;omitempty" json:"is_material_certification"`
	IsMillCertification     bool   `db:"is_mill_certification;omitempty" json:"is_mill_certification"`
	IsAppliedFinish         bool   `db:"is_applied_finish;omitempty" json:"is_applied_finish"`
	IsTestResult            bool   `db:"is_test_result;omitempty" json:"is_test_result"`
	IsDataSheet             bool   `db:"is_data_sheet;omitempty" json:"is_data_sheet"`
	IsOther                 bool   `db:"is_other;omitempty" json:"is_other"`
	IsSpareDelivery         bool   `db:"is_spare_delivery;omitempty" json:"is_spare_delivery"`
	IsMaterialComply        bool   `db:"is_material_comply;omitempty" json:"is_material_comply"`
	Comment                 string `db:"comment;omitempty" json:"comment"`
	Name                    string `db:"name" json:"name"`
	Signature               string `db:"signature;omitempty" json:"sign"`
	Status                  string `db:"status;omitempty" json:"status"`
	Remark                  string `db:"remark;omitempty" json:"remark"`
	CreatedAt               string `json:"created_at"`
	OrderId                 int    `db:"order_id;omitempty" json:"order_id"`
	Order                   Order  `json:"order"`
	UserId                  int    `db:"user_id" json:"user_id"`
	User                    User   `json:"user"`
	ReportType              string `db:"report_type;omitempty" json:"report_type"`
	InspEngSign             string `db:"insp_eng_sign;omitempty" json:"insp_eng_sign"`
	File1                   string `db:"file1;omitempty" json:"file1"`
	File2                   string `db:"file2;omitempty" json:"file2"`
	File3                   string `db:"file3;omitempty" json:"file3"`
	File4                   string `db:"file4;omitempty" json:"file4"`
}

type ReportResponse struct {
	Report       Report       `json:"report"`
	Order        Order        `json:"order"`
	Job          Job          `json:"job"`
	ClientReport ClientReport `json:"client_report"`
}

type Comment struct {
	Id             int    `db:"id" json:"id"`
	ApproveComment string `db:"approve_comment" json:"approve_comment"`
	RejectComment  string `db:"reject_comment" json:"reject_comment"`
	ReportId       int    `db:"report_id" json:"report_id"`
}

type ClientReport struct {
	Id              int       `db:"id" json:"id"`
	IsSpecification bool      `db:"is_specification" json:"is_specification"`
	Comment         string    `db:"comment" json:"comment"`
	Name            string    `db:"name" json:"name"`
	Signature       string    `db:"signature" json:"signature"`
	SigningDate     time.Time `db:"signing_date" json:"signing_date"`
	ReportId        int       `db:"report_id" json:"report_id"`
	Report          Report
	ClientName      string    `db:"client_name" json:"client_name"`
	ClientEngSign   string    `db:"client_eng_sign;omitempty" json:"client_eng_sign"`
	ClientSignDate  time.Time `db:"client_sign_date" json:"client_sign_date"`
	ClientComment   string    `db:"client_comment;omitempty" json:"client_comment"`
}
