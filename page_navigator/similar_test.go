package page_navigator

import "testing"

func TestSimilarIDValidation(t *testing.T) {
	for _, id := range []string{"", "2.1/../../login", "2.0", "bad", "https://other.test"} {
		if _, err := GetSimilarPageByID(id); err == nil {
			t.Errorf("accepted %q", id)
		}
	}
}
