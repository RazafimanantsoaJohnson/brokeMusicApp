package youtube

import "testing"

func TestGetStreamUrl(t *testing.T) {
	err := GetStreamUrl("https://www.youtube.com/watch?v=PVRbKHXwM58")

	// err = GetStreamUrl("https://www.youtube.com/watch?v=HE6hpNS2i6Y")

	if err != nil {
		t.Error(err)
	}
	t.Errorf("i")
}
