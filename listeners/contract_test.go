package listeners

import "testing"

func TestExistingAndLegacyCommandEnvelopes(t *testing.T) {
	for _, body := range []string{`{"user_id":42,"arguments":["101"]}`, `{"chat_id":42,"arguments":["101"]}`, `{"user_id":42,"chat_id":42,"arguments":["101"]}`} {
		id, args, err := parseCommand([]byte(body))
		if err != nil || id != 42 || len(args) != 1 || args[0] != "101" {
			t.Errorf("contract: id=%d args=%v err=%v", id, args, err)
		}
	}
	for _, body := range []string{`{"user_id":42,"chat_id":43}`, `{}`, `broken`} {
		if _, _, err := parseCommand([]byte(body)); err == nil {
			t.Errorf("accepted %s", body)
		}
	}
}
