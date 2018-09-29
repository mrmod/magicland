package magicland

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

var fixture []byte

func init() {
	fixture, _ = ioutil.ReadFile("pullRequestWebhook.json")

}
func TestDecoding(t *testing.T) {
	pullRequestWebhook := &PullRequestWebhook{}
	err := json.Unmarshal(fixture, pullRequestWebhook)
	if err != nil {
		t.Fatal(err)
	}
}
