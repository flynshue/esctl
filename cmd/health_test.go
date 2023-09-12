package cmd

import "testing"

func TestShowHealth(t *testing.T) {
	if err := showHealth(); err != nil {
		t.Error(err)
	}
}
