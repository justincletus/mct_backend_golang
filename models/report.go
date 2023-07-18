package models

import "gorm.io/gorm"

type Report struct {
	gorm.Model

	Id           uint   `db:"id" json:"id" gorm:"primaryKey; autoIncrement:true"`
	Title1       string `db:"title1" json:"title1"`
	File1        string `db:"file1" json:"file1"`
	Title2       string `db:"title2" json:"title2"`
	File2        string `db:"file2" json:"file2"`
	Title3       string `db:"title3" json:"title3"`
	File3        string `db:"file3" json:"file3"`
	Title4       string `db:"title4" json:"title4"`
	File4        string `db:"file4" json:"file4"`
	Status       string `db:"status" json:"status" gorm:"column:status;type:enum('approved','reject','pending');default:'pending'"`
	Uid          uint   `db:"uid" json:"uid" gorm:"index"`
	UserFullName string `db:"user_fullname" json:"user_fullname" gorm:"default:null"`
	User         User   `json:"user" gorm:"foreignKey:Uid"`
}

type MRI_Report struct {
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
	JobId                   int    `db:"job_id;omitempty" json:"job_id"`
	Job                     Job    `json:"job"`
}

type MRI_ReportResponse struct {
	MRI_Report MRI_Report `json:"mri_report"`
	Job        Job        `json:"job"`
}
