package models

import "gorm.io/gorm"

type BaseModel struct {
	gorm.Model
	CreatedBy *uint   `json:"created_by" gorm:"not null"`
	UpdatedBy *uint   `json:"updated_by" gorm:"not null"`
	RequestID *string `json:"request_id" gorm:"type:varchar(36);uniqueIndex;default:null"`
}

func setAuditFields(db *gorm.DB) {
	if db.Statement.Context == nil {
		return
	}

	// Extract values from context
	userID, _ := db.Statement.Context.Value("user_id").(uint)
	requestID, _ := db.Statement.Context.Value("request_id").(string)

	// Set CreatedBy if field exists
	if db.Statement.Schema.LookUpField("CreatedBy") != nil {
		db.Statement.SetColumn("CreatedBy", &userID)
	}

	// Set UpdatedBy if field exists
	if db.Statement.Schema.LookUpField("UpdatedBy") != nil {
		db.Statement.SetColumn("UpdatedBy", &userID)
	}

	// Set RequestID if field exists
	if requestID != "" && db.Statement.Schema.LookUpField("RequestID") != nil {
		db.Statement.SetColumn("RequestID", &requestID)
	}
}

func setUpdateAuditFields(db *gorm.DB) {
	if db.Statement.Context == nil {
		return
	}

	userID, _ := db.Statement.Context.Value("user_id").(uint)
	requestID, _ := db.Statement.Context.Value("request_id").(string)

	if db.Statement.Schema.LookUpField("UpdatedBy") != nil {
		db.Statement.SetColumn("UpdatedBy", &userID)
	}
	if requestID != "" && db.Statement.Schema.LookUpField("RequestID") != nil {
		db.Statement.SetColumn("RequestID", &requestID)
	}
}

func RegisterCallbacks(db *gorm.DB) {
	db.Callback().Create().Before("gorm:create").Register("set_audit_fields", setAuditFields)
	db.Callback().Update().Before("gorm:update").Register("set_update_fields", setUpdateAuditFields)
}
