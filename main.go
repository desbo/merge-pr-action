package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/google/go-github/v33/github"
)

const (
	eventNameVariable     = "GITHUB_EVENT_NAME"
	payloadPathVariable   = "GITHUB_EVENT_PATH"
	tokenVariable         = "INPUT_GITHUB_TOKEN"
	allowedUpdateVariable = "INPUT_ALLOWED_UPDATE"
	mergeMethodVariable   = "INPUT_MERGE_METHOD"
)

var (
	allowedEvents = []string{"pull_request", "pull_request_target"}
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

func checkAllowedEvent(event string) bool {
	for _, i := range allowedEvents {
		if i == event {
			return true
		}
	}
	return false
}

func main() {
	eventName := getRequiredEnvVar(eventNameVariable)
	if checkAllowedEvent(eventName) == false {
		log.Printf("event `%v` is not one of %v, exiting", eventName, allowedEvents)
		os.Exit(0)
		return
	}

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

	if event.PullRequest.Title == nil {
		log.Fatalf("no pull request title in event payload")
	}

	upgradeTypeString := os.Getenv(allowedUpdateVariable)

	if strings.TrimSpace(strings.ToLower(upgradeTypeString)) == "any" {
		log.Printf("any upgrade type allowed, merging")
		merge(&event.PullRequest)
		os.Exit(0)
		return
	}

	allowedUpgrade, err := parseUpgradeType(upgradeTypeString)
	if err != nil {
		log.Fatalf("error parsing allowed upgrade type: %v", err.Error())
	}

	upgrade, err := parseVersionUpgrade(*event.PullRequest.Title)
	if err != nil {
		log.Fatalf("error parsing upgrade from PR title %v: %v", event.PullRequest.Title, err.Error())
	}
	upgradeType := upgrade.UpgradeType()

	log.Printf("detected upgrade: %v", upgrade)

	if !allowed(allowedUpgrade, upgradeType) {
		log.Printf("%v upgrade not allowed, skipping", upgradeType)
		os.Exit(0)
	}

	merge(&event.PullRequest)
	os.Exit(0)
}

func merge(pr *github.PullRequest) {
	token := getRequiredEnvVar(tokenVariable)
	mergeMethod := getRequiredEnvVar(mergeMethodVariable)
	client := newAuthenticatedClient(token)

	if err := client.mergePR(pr, mergeMethod); err != nil {
		log.Fatalf("error merging PR: %v", err.Error())
	}
}
