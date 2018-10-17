package magicland

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/context"
)

const closedPullRequest = "closed"

// WebhookHandler Handle GitHub pull request webhook
func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[%s] %s", r.Method, r.URL.Path)
	w.Header().Add("Content-Type", "application/json")
	webhookHandler(w, r)
}

func webhookHandler(w http.ResponseWriter, r *http.Request) {
	pullRequestWebhook := &PullRequestWebhook{}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		badRequest("Failed to read request body", err, w)
		return
	}
	if err := json.Unmarshal(body, &pullRequestWebhook); err != nil {
		log.Printf("JSON\n\t%s\n\n", string(body))
		badRequest("Failed to decode webhook", err, w)
		return
	}
	if pullRequestWebhook.Action != closedPullRequest || !pullRequestWebhook.PullRequest.Merged {
		badRequest("Pull request must be closed and merged", fmt.Errorf("Pull request must be closed and merged"), w)
		return
	}
	log.Println("Integrating repository")
	// Build the deployable configuration
	serviceName := strings.Join([]string{
		strconv.Itoa(pullRequestWebhook.Repository.Owner.ID),
		pullRequestWebhook.Repository.Name,
		pullRequestWebhook.Repository.Owner.Login,
	}, "-")
	gitConfig := GitConfiguration{
		User:          pullRequestWebhook.PullRequest.User.Login,
		ServiceName:   serviceName,
		BranchName:    "master",
		RepositoryURL: pullRequestWebhook.Repository.URL,
	}
	if err := PublicClone(gitConfig); err != nil {
		failedDependency(fmt.Sprintf("Failed to clone service repository for %s", serviceName), err, w)
		return
	}
	if err := saveService(gitConfig); err != nil {
		failedDependency(fmt.Sprintf("Failed to save service %s", serviceName), err, w)
		return
	}
	log.Println("Building runtime for service", serviceName)
	rtConfig := newRuntimeConfiguration("host01.local", 8000, gitConfig)

	ctx := context.Background()
	runnableContainer, err := buildContainer(ctx, rtConfig)
	if err != nil {
		failedDependency("Failed to build runnable container", err, w)
		return
	}
	runningContainer, err := runContainer(ctx, *runnableContainer)
	if err != nil {
		failedDependency("Failed to run container", err, w)
		return
	}
	log.Printf("%s UP in container %s", serviceName, runningContainer.StagedContainer.ID)
	// End build config
	w.WriteHeader(http.StatusOK)
}
