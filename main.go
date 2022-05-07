package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

const (
	// Timeout used for the context in the request to github api
	REQUEST_TIMEOUT = 5 * time.Second
	// Time used to sleep between calls of the github api
	SLEEP_TIME_BETWEEN_CALLS = 10 * time.Second
	// Number of times we will try to query the github api
	REPEAT = 20
)

var (
	sha           = flag.String("sha", "", "Commit has to check the status")
	statusContext = flag.String("context", "", "Context to check the status")
	repository    = flag.String("repository", "", "Repository where the status and commit are")
	token         = flag.String("token", "", "Github token to use against api")
)

type Response struct {
	Statuses []Status `json:"statuses"`
}

type Status struct {
	State   string `json:"state"`
	Context string `json:"context"`
}

func getStatus(ctx context.Context, sha, repository, token string) (*Response, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/commits/%s/status", repository, sha)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	cli := &http.Client{}

	resp, err := cli.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get commit status")
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	parsed := &Response{}
	err = json.Unmarshal(bodyBytes, parsed)

	if err != nil {
		return nil, err
	}

	return parsed, nil
}

func findStateForContext(resp *Response, statusContext string) string {
	if resp == nil {
		return ""
	}

	for _, status := range resp.Statuses {
		if statusContext == status.Context {
			return status.State
		}
	}

	return ""
}

func main() {
	flag.Parse()

	for i := 0; i < REPEAT; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), REQUEST_TIMEOUT)
		defer cancel()

		resp, err := getStatus(ctx, *sha, *repository, *token)
		if err != nil {
			fmt.Printf("Failed to get the status, wait and retry. Err: %+v\n", err)
		}

		state := findStateForContext(resp, *statusContext)

		switch state {
		case "", "pending":
			fmt.Printf("The deployment is in progress. State: %s\n", state)
			time.Sleep(SLEEP_TIME_BETWEEN_CALLS)
		case "success":
			fmt.Println("The deployment succeed")
			os.Exit(0)
		default:
			fmt.Printf("The deployment has failed. State: %s\n", state)
			os.Exit(2)
		}
	}

	os.Exit(3)
}
