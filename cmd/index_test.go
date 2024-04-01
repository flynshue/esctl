package cmd

import (
	"fmt"
	"testing"
)

func TestIndex_ListIndexVersion(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
	}{
		{"all", "*"},
		{"indexPattern", ".kibana*"},
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
	testCases := []struct {
		name       string
		idxPattern []string
	}{
		{"allIndexes", []string{"*"}},
		{"indexPattern", []string{".kibana*"}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := showIdxSizes(tc.idxPattern); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestIndex_ListIndexTemplates(t *testing.T) {
	testCases := []struct {
		name    string
		pattern string
	}{
		{"all", "*"},
		{"indexPattern", ".monitoring-*"},
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
		{"MultiIndexPattern", false, []string{".fleet*", ".kibana*"}},
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

func TestIndex_GetIndexReadonly(t *testing.T) {
	index := "test-read-only"
	if err := escConsole("put", index, nil); err != nil {
		t.Error(err)
	}
	defer func() {
		setIndexSettings(index, `{"index.blocks.read_only": "false"}`)
		deleteIndex([]string{index})
	}()
	if err := escConsole("put", fmt.Sprintf("%s/_block/read_only", index), nil); err != nil {
		t.Error(err)
	}
	if err := getIndexReadonly(); err != nil {
		t.Error(err)
	}
}

func TestIndex_SetIndexReplicas(t *testing.T) {
	index := "test-idx-replicas"
	if err := escConsole("put", index, nil); err != nil {
		t.Error(err)
	}
	defer func() {
		escConsole("delete", index, nil)
	}()
	setIndexReplicas(index, 3)
	if err := listIndexSettingsSummary(index); err != nil {
		t.Error()
	}
}

func TestIndex_SetIndexAutoExpand(t *testing.T) {
	index := "test-auto-expand-0001"
	if err := escConsole("put", index, nil); err != nil {
		t.Error(err)
	}
	defer func() {
		escConsole("delete", index, nil)
	}()
	testCases := []struct {
		name       string
		autoExpand string
	}{
		{"invalidOption", "foobar"},
		{"validOption", "0-1"},
		{"disableAutoExpand", "false"},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := setIndexAutoExpand(index, tc.autoExpand); err != nil {
				t.Error(err)
			}
			listIndexSettingsSummary(index)
		})
	}
}

func TestIndex_CatIndices(t *testing.T) {
	b, err := catIndices("", "", "", []string{"*"})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
}

func TestIndex_DeleteIndex(t *testing.T) {
	disableDestructiveRequires()
	idxPrefix := "test-del-idx"
	for i := 1; i <= 3; i++ {
		escConsole("put", fmt.Sprintf("%s-000%d", idxPrefix, i), nil)
	}
	b, err := catIndices("", "", "", []string{idxPrefix + "*"})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(b))
	if err := deleteIndex([]string{idxPrefix + "*"}); err != nil {
		t.Error(err)
	}
}

func TestIndex_DisableReadOnly(t *testing.T) {
	indices := []string{
		"test-all-indices-001",
		"test-all-indices-002",
		"test-all-indices-003",
		"test-readonly-one",
		"test-readonly-101",
		"test-readonly-102",
		"test-readonly-103",
	}
	testCases := []struct {
		name         string
		indexPattern string
	}{
		{"oneIndex", "test-readonly-one"},
		{"pattern", "test-readonly-1*"},
		{"all", "*"},
	}
	for _, idx := range indices {
		escConsole("put", idx, nil)
		escConsole("put", fmt.Sprintf("%s/_block/read_only", idx), nil)
	}
	defer deleteIndex([]string{"test-*"})
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			getIndexReadonly()
			if err := disableReadOnlyIdx(tc.indexPattern); err != nil {
				t.Error(err)
			}
			getIndexReadonly()
		})
	}
}
