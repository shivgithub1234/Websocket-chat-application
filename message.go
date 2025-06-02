// package main
// type Message struct {
// 	Sender string `json:"sender,omitempty"`
// 	Recipient string `json:"recipient,omitempty"`
// 	Content string `json:"content,omitempty"`	
// }

package main

// Message represents a chat message
type Message struct {
    Type      string `json:"type,omitempty"`
    Sender    string `json:"sender,omitempty"`
    Recipient string `json:"recipient,omitempty"`
    Content   string `json:"content,omitempty"`
    IsOwnMessage bool `json:"isOwnMessage,omitempty"`
}