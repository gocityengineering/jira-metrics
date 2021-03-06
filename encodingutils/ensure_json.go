package encodingutils

import (
	"encoding/json"
	"fmt"
	"unicode/utf8"

	"github.com/ghodss/yaml"
)

// ensure YAML as well as JSON can be read
// applies only to file-based processing; the server only accepts JSON
func EnsureJson(a *[]byte, isJSON bool) error {
	if len(*a) == 0 {
		return fmt.Errorf("input must not be empty")
	}

	if utf8.Valid(*a) == false {
		return fmt.Errorf("input must be valid UTF-8")
	}

	// attempt to parse JSON first
	var any interface{}
	err := json.Unmarshal(*a, &any)

	// input is valid JSON
	if err == nil {
		return nil
	}

	// exit condition: flagged as JSON but error found
	if isJSON {
		return fmt.Errorf("invalid JSON: %s", err.Error())
	}

	// not JSON
	json, err := yaml.YAMLToJSON(*a)
	if err != nil {
		return fmt.Errorf("invalid YAML: %s", err.Error())
	}

	// successful conversion
	*a = json

	return nil
}
