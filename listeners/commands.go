package listeners

import (
	"encoding/json"
	"fmt"
)

// Older Kino Cat text services use chat_id; the pinned core uses user_id.
// Do not silently route a conflicting envelope to either recipient.
func parseCommand(data []byte) (int64, []string, error) {
	var msg struct {
		UserID    int64    `json:"user_id"`
		ChatID    int64    `json:"chat_id"`
		Arguments []string `json:"arguments"`
	}
	if err := json.Unmarshal(data, &msg); err != nil {
		return 0, nil, err
	}
	if msg.UserID != 0 && msg.ChatID != 0 && msg.UserID != msg.ChatID {
		return 0, nil, fmt.Errorf("conflicting command identities")
	}
	id := msg.UserID
	if id == 0 {
		id = msg.ChatID
	}
	if id == 0 {
		return 0, nil, fmt.Errorf("command recipient missing")
	}
	return id, msg.Arguments, nil
}
