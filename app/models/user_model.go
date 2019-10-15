package models

/*
CREATE TABLE GOS_USER (
user_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
name VARCHAR(100) NOT NULL,
email VARCHAR(200),
password VARCHAR(100),
last_login int(10),
failed_login_attempt int(8),
date_created int(10),
date_updated int(10)) ENGINE=InnoDB;
*/
type User struct {
	UserId             int64  `json:"userId"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Password           string `json:"_"`
	LastLogin          int    `json:"lastLogin"`
	FailedLoginAttempt int    `json:"failedLoginAttempt"`
	DateCreated        int64    `json:"dateCreated"`
	DateUpdated        int64    `json:"dateUpdated"`
}