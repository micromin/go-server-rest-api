package models

type User struct {
	Name               string `json:"name"`
	Email              string `json:"email"`
	Password           string `json:"_"`
	IP                 string `json:"ip"`
	LastLogin          int    `json:"lastLogin"`
	FailedLoginAttempt int    `json:"failedLoginAttempt"`
	DateCreated        int    `json:"dateCreated"`
	DateUpdated        int    `json:"dateUpdated"`
}
