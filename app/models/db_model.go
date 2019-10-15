package models

// swagger:model User
// User model
type User struct {
	UserId             int64  `json:"userId"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Password           string `json:"password"`
	LastLogin          int    `json:"lastLogin,omitempty"`
	FailedLoginAttempt int    `json:"failedLoginAttempt,omitempty"`
	DateCreated        int64  `json:"dateCreated,omitempty"`
	DateUpdated        int64  `json:"dateUpdated,omitempty"`
}

// swagger:model Task
// Task model
type Task struct {
	TaskId        int64  `json:"taskId"`
	UserId        int64  `json:"userId"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	DateCreated   int64  `json:"dateCreated,omitempty"`
	DateUpdated   int64  `json:"dateUpdated,omitempty"`
	DueDate       int64  `json:"dueDate,omitempty"`
	DateCompleted int64  `json:"dateCompleted,omitempty"`
}
