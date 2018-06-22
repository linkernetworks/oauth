package entity

// NotificationPreference contains flags indicating whether to send SMS / Email to the user,
// and determine what kind of messages will be sent.
type NotificationPreference struct {
	EnableSMS   bool `json:"enableSMS" bson:"enableSMS"`
	EnableEmail bool `json:"enableEmail" bson:"enableEmail"`

	WhenJobStart   bool `json:"whenJobStart" bson:"whenJobStart"`
	WhenJobSuccess bool `json:"whenJobSuccess" bson:"whenJobSuccess"`
	WhenJobFail    bool `json:"whenJobFail" bson:"whenJobFail"`
	WhenJobStop    bool `json:"whenJobStop" bson:"whenJobStop"`
	WhenJobDelete  bool `json:"whenJobDelete" bson:"whenJobDelete"`
}
