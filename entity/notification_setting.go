package entity

import (
	"bitbucket.org/linkernetworks/aurora/src/config"
)

// NotificationSetting contains flags indicating whether to send SMS / Email to the user,
// and determine what kind of messages will be sent.
type NotificationSetting struct {
	EnableSMS   bool `json:"enable_sms" bson:"enable_sms"`
	EnableEmail bool `json:"enable_email" bson:"enable_email"`

	OnJobStart   bool `json:"on_job_start" bson:"on_job_start"`
	OnJobSuccess bool `json:"on_job_success" bson:"on_job_success"`
	OnJobFail    bool `json:"on_job_fail" bson:"on_job_fail"`
	OnJobStop    bool `json:"on_job_stop" bson:"on_job_stop"`
	OnJobDelete  bool `json:"on_job_delete" bson:"on_job_delete"`
}

// NewNotificationSettingFrom creates a NotificationSetting from NotificationConfig
func NewNotificationSettingFrom(c *config.NotificationConfig) NotificationSetting {
	return NotificationSetting{
		EnableSMS:    c.EnableSMS,
		EnableEmail:  c.EnableEmail,
		OnJobStart:   c.OnJobStart,
		OnJobSuccess: c.OnJobSuccess,
		OnJobFail:    c.OnJobFail,
		OnJobStop:    c.OnJobStop,
		OnJobDelete:  c.OnJobDelete,
	}
}
