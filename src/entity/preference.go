package entity

// Preference saves settings of a user
type Preference struct {
	Notification NotificationPreference `json:"notification" bson:"notification"`
}
