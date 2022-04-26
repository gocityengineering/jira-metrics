package main

import (
	"testing"
)

func TestRealMain(t *testing.T) {
	var tests = []struct {
		description string
		configPath  string
		schemaPath  string
		user        string
		token       string
		expected    int
	}{
		{"dryrun_full_config", "config/config.yaml", "schema/schema.yaml", "user", "token", 0},
		{"dryrun_empty_config", "config/config_empty.yaml", "schema/schema.yaml", "user", "token", 0},
		{"invalid_config", "config/config_invalid.yaml", "schema/schema.yaml", "user", "token", 1},
		{"broken_config", "config/config_broken.yaml", "schema/schema.yaml", "user", "token", 1},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			retVal := realMain(test.configPath, test.schemaPath, test.user, test.token, true)
			if retVal != test.expected {
				t.Errorf("Unexpected return value '%d' for %s", retVal, test.description)
			}
		})
	}
}
