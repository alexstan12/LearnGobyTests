package selectTest

import "testing"

func TestRacer(t *testing.T){
	slowURL := "http://www.facebook.com"
	fastUrl := "http://www.quii.co.uk"

	want:= fastUrl
	got:= Racer(slowURL, fastUrl)
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
