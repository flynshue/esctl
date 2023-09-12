package esapi

import (
	"testing"
)

func TestRenderEndpoint(t *testing.T) {
	testCases := []struct {
		name     string
		baseUrl  string
		endpoint string
		want     string
	}{
		{"SlashAfterBase", "https://localhost:9200/", "search", "https://localhost:9200/search"},
		{"SlashBeforeEndpoint", "https://localhost:9200", "/search", "https://localhost:9200/search"},
		{"SlashOnBoth", "https://localhost:9200/", "/search", "https://localhost:9200/search"},
		{"NoSlash", "https://localhost:9200", "search", "https://localhost:9200/search"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := renderEndpoint(tc.baseUrl, tc.endpoint)
			if got != tc.want {
				t.Errorf("got = %s, want = %s", got, tc.want)
			}
		})
	}
}

func TestClientDo(t *testing.T) {
	client := NewClient("https://httpbin.org")
	client.SkipTLS()
	auth := BasicAuth{Username: "foo", Password: "bar"}
	client.SetAuth(auth)
	testCases := []struct {
		name     string
		method   string
		endpoint string
		body     []byte
	}{
		{"basic-auth", "GET", "basic-auth/foo/bar", nil},
		{"data", "post", "anything", []byte(`{"anything":"fooBar"}`)},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, err := client.Do(tc.method, tc.endpoint, tc.body)
			if err != nil {
				t.Error(err)
			}
			t.Log(string(b))
		})
	}
}
