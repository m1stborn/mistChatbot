package user

import "fmt"

type User struct {
	Subscribes []string         `json:"Subscribes"`
	Profile    `json:"Profile"` // not clear
}

type Profile struct {
	Account         string `json:"account"`
	Type            string `json:"type,omitempty"`
	Email           string `json:"email"`
	Line            string `json:"line"`
	LineAccessToken string `json:"lineAccessToken"`
}

var u = User{Subscribes: []string{"1", "2", "3"}}

func init() {
	fmt.Println(u)
}
