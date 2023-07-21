package models

type Order struct {
	Id              int    `db:"id" json:"id"`
	Project         string `db:"project" json:"project"`
	RequisitionNo   string `db:"requisition_no" json:"requisition_no"`
	PurchaseOrderNo string `db:"purchase_order_no" json:"purchase_order_no"`
	DeliveryNoteNo  string `db:"delivery_note_no;omitempty" json:"delivery_note_no"`
	DateOFDelivery  string `db:"date_of_delivery;omitempty" json:"date_of_delivery"`
	Description     string `db:"description;omitempty" json:"description"`
	JobId           int    `db:"job_id" json:"job_id"`
	Job             Job
}
