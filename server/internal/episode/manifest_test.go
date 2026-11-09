package episode

import "testing"

func TestManifest(t *testing.T) {
	p := []byte(`{"fps":50,"action_dim":14,"episodes":[{"episode_index":0,"length":400,"from_index":0,"to_index":400}]}`)
	eps, total, _, _, err := ParseManifest(p)
	if err != nil || total != 400 || len(eps) != 1 {
		t.Fatalf("unexpected result: %v", err)
	}
}
