package model

type Channel string

const (
	EMAIL Channel = "email"
	SLACK Channel = "slack"
	INAPP Channel = "inapp"
)

type Content struct {
	Subject         string   `json:"subject"`
	SubPlaceHolder  []string `json:"subplaceholder"`
	Body            string   `json:"body"`
	BodyPlaceHolder []string `json:"bodyplaceholder"`
}

type Notification struct {
	ID            string  `json:"id"`
	To            string  `json:"to"`
	From          string  `json:"from"`
	Template      *string `json:"template"`
	SendTimeInSec *int    `json:"time"`
	Message       Content `json:"message"`
	Channel       Channel `json:"channel"`
}

type Template struct {
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

type APIResponse struct {
	Code    int
	Message string
}
