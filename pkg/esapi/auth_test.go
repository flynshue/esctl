package esapi

import "testing"

func TestBasicAuth(t *testing.T) {
	want := "Basic Zm9vQmFyOmZha2VQYXNz"
	auth := BasicAuth{Username: "fooBar", Password: "fakePass"}
	got := auth.AuthorizationHeader()
	if got != want {
		t.Errorf("got = %s, want = %s", got, want)
	}
}
