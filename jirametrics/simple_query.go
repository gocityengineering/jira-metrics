package jirametrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func SimpleQuery(config *Config, username string, token string, dryrun bool, query string) (Result, error) {
	server := config.Server
	var result = Result{}

	// fetch response from Jira
	bytes, err := json.Marshal(Data{Jql: query})
	body := strings.NewReader(string(bytes))

	req, err := http.NewRequest("POST", "https://"+server+"/rest/api/2/search", body)
	if err != nil {
		return result, fmt.Errorf("failed to create request: %v", err)
	}
	req.SetBasicAuth(username, token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return result, fmt.Errorf("failed to parse response body: %v", err)
	}

	return result, nil
}
