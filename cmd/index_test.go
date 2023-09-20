package cmd

import (
	"testing"
)

func TestListIndexVersion(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
	}{
		{"all", "*"},
		{"kibana", ".kibana*"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := listIndexVersion(tc.pattern); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestShowIdxSizes(t *testing.T) {
	if err := showIdxSizes(); err != nil {
		t.Error(err)
	}
}

func TestListIndexTemplates(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
	}{
		{"all", "*"},
		{"monitoring", ".monitoring-*"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := listIndexTemplates(tc.pattern); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestGetIndexTemplate(t *testing.T) {
	if err := getIndexTemplate("logs"); err != nil {
		t.Error(err)
	}
}

func TestListIndexTemplatesLegacy(t *testing.T) {
	if err := listIndexTemplatesLegacy("*"); err != nil {
		t.Error(err)
	}
}

func TestListIndexDate(t *testing.T) {
	testCases := []struct {
		name       string
		local      bool
		idxPattern []string
	}{
		{"UTC", false, []string{"*"}},
		{"LocalTime", true, []string{"*"}},
		{"IndexPattern", false, []string{".fleet*"}},
		{"MultiIndexPattern", false, []string{".fleet*", "cust*"}},
	}
	for _, tc := range testCases {
		localTime = tc.local
		t.Run(tc.name, func(t *testing.T) {
			if err := listIndexDate(tc.idxPattern); err != nil {
				t.Error(err)
			}
		})
	}
}
