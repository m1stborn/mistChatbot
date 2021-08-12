package user

import "fmt"

type User struct {
	Account         string `json:"account"`
	Type            string `json:"type,omitempty"`
	Email           string `json:"email"`
	Line            string `json:"line"`
	LineAccessToken string `json:"lineAccessToken"`

	Subscribes []string `json:"Subscribes"`
}

var TestUser = User{
	Subscribes: []string{"Twitch"},
}

var ListUsers []string

func init() {
	fmt.Println(TestUser, ListUsers)
}
