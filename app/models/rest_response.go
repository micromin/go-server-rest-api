package models

// Response is a generic rest response
type Response struct {
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type Paged struct {
	Items  interface{} `json:"items,omitempty"`
	LastId int64       `json:"lastId,omitempty"`
	Limit  int         `json:"limit,omitempty"`
}
