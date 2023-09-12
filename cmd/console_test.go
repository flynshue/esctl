package cmd

import "testing"

func TestConsole(t *testing.T) {
	testCases := []struct {
		name     string
		method   string
		endpoint string
		data     []byte
	}{
		{"catIndex", "get", "/_cat/indices", nil},
		{"createIndex", "put", "/test",
			[]byte(`{
"settings": {
	"number_of_shards": 1
},
"mappings": {
	"properties": {
	"field1": { "type": "text" }
	}
}
}`)},
		{"getIndex", "get", "/test/_settings", nil},
		{"deleteIndex", "delete", "/test", nil},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := console(tc.method, tc.endpoint, tc.data); err != nil {
				t.Error(err)
			}
		})
	}
}
