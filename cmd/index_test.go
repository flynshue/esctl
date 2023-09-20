package cmd

import (
	"testing"
)

func TestIndex_ListIndexVersion(t *testing.T) {
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

func TestIndex_ShowIdxSizes(t *testing.T) {
	if err := showIdxSizes(); err != nil {
		t.Error(err)
	}
}

func TestIndex_ListIndexTemplates(t *testing.T) {
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

func TestIndex_GetIndexTemplate(t *testing.T) {
	if err := getIndexTemplate("logs"); err != nil {
		t.Error(err)
	}
}

func TestIndex_ListIndexTemplatesLegacy(t *testing.T) {
	if err := listIndexTemplatesLegacy("*"); err != nil {
		t.Error(err)
	}
}

func TestIndex_ListIndexDate(t *testing.T) {
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

func TestIndex_GetIndexSettings(t *testing.T) {
	testCases := []struct {
		name       string
		idxPattern string
	}{
		{"all", "*"},
		{"indexPattern", ".fleet-*"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := listIndexSettingsSummary(tc.idxPattern); err != nil {
				t.Error(err)
			}
		})

	}
}

func TestIndex_SetIndexReplicas(t *testing.T) {
	index := "test-idx-replicas"
	if err := console("put", index, nil); err != nil {
		t.Error(err)
	}
	defer func() {
		console("delete", index, nil)
	}()
	setIndexReplicas(index, 3)
	if err := listIndexSettingsSummary(index); err != nil {
		t.Error()
	}
}
