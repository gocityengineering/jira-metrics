package jirametrics

import (
	"fmt"
	"net/http"
	"os"
)

func MinimalGetRequest(server string, username string, token string, dryrun bool) error {
	if dryrun {
		return nil
	}
	req, err := http.NewRequest("GET", "https://"+server+"/rest/api/3/serverInfo", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, `Can't create minimal GET request: %v`, err)
		return err
	}
	req.SetBasicAuth(username, token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Fprintf(os.Stderr, `Can't process minimal GET request: %v`, err)
		return err
	}

	defer resp.Body.Close()

	return nil
}
