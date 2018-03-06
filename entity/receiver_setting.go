package entity

// ReceiverSetting contains flags indicating whether to send SMS / Email to the user,
// and determine what kind of messages will be sent.
type ReceiverSetting struct {
	EnableSMS   bool `json:"enable_sms" bson:"enable_sms"`
	EnableEmail bool `json:"enable_email" bson:"enable_email"`

	WhenJobStart   bool `json:"when_job_start" bson:"when_job_start"`
	WhenJobSuccess bool `json:"when_job_success" bson:"when_job_success"`
	WhenJobFail    bool `json:"when_job_fail" bson:"when_job_fail"`
	WhenJobStop    bool `json:"when_job_stop" bson:"when_job_stop"`
	WhenJobDelete  bool `json:"when_job_delete" bson:"when_job_delete"`
}
