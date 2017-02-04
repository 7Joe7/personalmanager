package exporter

type exportConfig struct {
	SmtpAddress        string `json:"smtpAddress"`
	SmtpPort           int    `json:"smtpPort"`
	AdminEmailAddress  string `json:"adminEmailAddress"`
	AdminEmailPassword string `json:"adminEmailPassword"`
}
