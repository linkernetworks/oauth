package mongo

// MongoConfig
type MongoConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	TimeOut  int64  `json:"timeout"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}
