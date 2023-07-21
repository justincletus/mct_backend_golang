package models

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
}

type ReportResponse struct {
	Report Report `json:"report"`
	Order  Order  `json:"order"`
	Job    Job    `json:"job"`
}

type Comment struct {
	Id             int    `db:"id" json:"id"`
	ApproveComment string `db:"approve_comment" json:"approve_comment"`
	RejectComment  string `db:"reject_comment" json:"reject_comment"`
	MriReportId    int    `db:"mri_report_id" json:"mri_report_id"`
}
