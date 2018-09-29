package magicland

import (
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

// NewSSHClient creates a client configuration suitable to
// fetching a remote repository which requires authentication
func NewSSHClient(gitConfig GitConfiguration) (*ssh.ClientConfig, error) {
	var hostKey ssh.PublicKey
	privateKeyBuf, err := ioutil.ReadFile(privateKeyFor(gitConfig))
	if err != nil {
		return nil, err
	}
	privateKey, err := ssh.ParsePrivateKey(privateKeyBuf)
	if err != nil {
		return nil, err
	}
	config := &ssh.ClientConfig{
		User: gitConfig.User,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(privateKey),
		},
		HostKeyCallback: ssh.FixedHostKey(hostKey),
	}
	return config, nil
}

//
func privateKeyFor(gitConfig GitConfiguration) string {
	return "TODO Real value"
}
