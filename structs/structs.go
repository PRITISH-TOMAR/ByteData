package structs

import(
	"encoding/json"
)

type Communicators struct {
	Enc *json.Encoder
	Dec *json.Decoder
}


type Message struct {
	Type     string   `json:"type"`
	Field    string   `json:"field,omitempty"`
	Username string   `json:"username,omitempty"`
	Password string   `json:"password,omitempty"`
	Command  string   `json:"command,omitempty"`
	Message  string   `json:"message,omitempty"`
	Data     []string `json:"data,omitempty"`
}


