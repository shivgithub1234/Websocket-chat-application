package main
type Message struct {
	Sender string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content string `json:"content,omitempty"`

}