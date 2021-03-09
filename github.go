package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/google/go-github/v33/github"
	"golang.org/x/oauth2"
)

type authenticatedGitHubClient struct {
	ctx    context.Context
	client *github.Client
}

func newAuthenticatedClient(token string) *authenticatedGitHubClient {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	return &authenticatedGitHubClient{ctx, client}
}

func (c *authenticatedGitHubClient) refetchPR(pr *github.PullRequest) error {
	pr, _, err := c.client.PullRequests.Get(
		c.ctx,
		pr.Base.Repo.Owner.GetName(),
		pr.Base.Repo.GetName(),
		pr.GetNumber(),
	)

	return err
}

func (c *authenticatedGitHubClient) getMergeableState(pr *github.PullRequest, attempts int, maxAttempts int) (string, error) {
	if attempts+1 == maxAttempts {
		return "", fmt.Errorf("gave up waiting for mergeable state to be determined")
	}

	state := pr.GetMergeableState()
	log.Printf("current PR state: %v", state)

	if strings.EqualFold(state, "unknown") {
		time.Sleep(time.Duration(math.Pow(2, float64(attempts+1)) * float64(time.Second)))
		return c.getMergeableState(pr, attempts+1, maxAttempts)
	}

	return state, nil
}

func (c *authenticatedGitHubClient) mergePR(pr *github.PullRequest, maxRetries int) error {
	state := pr.GetMergeableState()

	if strings.EqualFold(state, "conflicting") {
		return errors.New("PR has conflicts, will not merge")
	}

	options := &github.PullRequestOptions{
		MergeMethod: "rebase",
	}

	result, _, err := c.client.PullRequests.Merge(
		c.ctx,
		pr.Base.Repo.Owner.GetLogin(),
		pr.Base.Repo.GetName(),
		pr.GetNumber(),
		"",
		options,
	)

	if err != nil {
		return err
	}

	if !result.GetMerged() {
		return fmt.Errorf("PR was not merged: %v", result.GetMessage())
	}

	log.Printf(result.GetMessage())
	return nil
}
