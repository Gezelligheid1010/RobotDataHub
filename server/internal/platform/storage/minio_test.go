package storage

import "testing"

func TestUnsafeKey(t *testing.T) {
	for _, s := range []string{"../abc", "/root", "C:/test"} {
		if _, err := SafeObjectKey("id", s); err == nil {
			t.Fatalf("accepted %s", s)
		}
	}
}
