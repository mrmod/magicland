package magicland

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

var unmergedPullRequest []byte
var mergedPullRequest []byte

func init() {
	unmergedPullRequest, _ = ioutil.ReadFile("pullRequestWebhook.json")
	mergedPullRequest, _ = ioutil.ReadFile("mergedPullRequestWebhook.json")
}

func TestHandler(t *testing.T) {
	body := bytes.NewReader([]byte("Invalid body"))
	request := httptest.NewRequest("POST", "/", body)
	writer := &httptest.ResponseRecorder{}

	webhookHandler(writer, request)
	if writer.Code != http.StatusBadRequest {
		t.Fatalf("Expected %d, got %d", http.StatusBadRequest, writer.Code)
	}

	t.Log("Given an unmerged PR")
	writer = &httptest.ResponseRecorder{}
	body = bytes.NewReader(unmergedPullRequest)
	request = httptest.NewRequest("Post", "/", body)
	webhookHandler(writer, request)
	if writer.Code != http.StatusAccepted {
		t.Fatalf("Expected %d, got %d", http.StatusAccepted, writer.Code)
	}

	t.Log("Given a merged PR")
	writer = &httptest.ResponseRecorder{}
	body = bytes.NewReader(mergedPullRequest)
	request = httptest.NewRequest("Post", "/", body)
	webhookHandler(writer, request)
	if writer.Code != http.StatusOK {
		t.Fatalf("Expected %d, got %d", http.StatusOK, writer.Code)
	}
}
