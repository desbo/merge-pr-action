package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/go-github/v33/github"
)

const (
	payloadPathVariable   = "GITHUB_EVENT_PATH"
	tokenVariable         = "INPUT_GITHUB_TOKEN"
	allowedUpdateVariable = "INPUT_ALLOWED_UPDATE"
	maxRetries            = 10
)

type pullRequestEvent struct {
	PullRequest github.PullRequest `json:"pull_request"`
}

func getRequiredEnvVar(name string) string {
	value := os.Getenv(name)
	if value == "" {
		log.Fatalf("required env variable %v not set", name)
	}

	return value
}

func main() {
	payloadPath := getRequiredEnvVar(payloadPathVariable)

	payload, err := ioutil.ReadFile(payloadPath)
	if err != nil {
		log.Fatalf("error opening %v: %v", payloadPath, err.Error())
	}

	var event pullRequestEvent
	err = json.Unmarshal(payload, &event)
	if err != nil {
		log.Fatalf("error parsing event JSON: %v", err.Error())
	}

	upgrade, err := parseVersionUpgrade(*event.PullRequest.Title)
	if err != nil {
		log.Fatalf("error parsing upgrade from PR title %v: %v", event.PullRequest.Title, err.Error())
	}
	upgradeType := upgrade.UpgradeType()

	allowedUpgrade, err := parseUpgradeType(os.Getenv(allowedUpdateVariable))
	if err != nil {
		log.Fatalf("error parsing allowed upgrade type: %v", err.Error())
	}

	if !allowed(allowedUpgrade, upgradeType) {
		log.Printf("upgrade of type %v not allowed, skipping", upgradeType)
	}

	token := getRequiredEnvVar(tokenVariable)
	client := newAuthenticatedClient(token)

	if err := client.mergePR(&event.PullRequest, maxRetries); err != nil {
		log.Fatalf("error merging PR: %v", err.Error())
	}

	os.Exit(0)
}
