package magicland

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const closedPullRequest = "closed"

func handler(w http.ResponseWriter, r *http.Request) {
	pullRequestWebhook := &PullRequestWebhook{}
	body, decodeError := ioutil.ReadAll(r.Body)
	if decodeError != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(body, &pullRequestWebhook); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if pullRequestWebhook.Action == closedPullRequest && pullRequestWebhook.PullRequest.Merged {
		w.WriteHeader(http.StatusOK)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8000", nil)
}
