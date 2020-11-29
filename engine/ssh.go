// Copyright 2018 Drone.IO Inc
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package engine

import (
	"bytes"
	"io"
	"io/ioutil"

	"github.com/drone/autoscaler"
	"golang.org/x/crypto/ssh"
)

type sshClient struct {
	*ssh.Client
}

// clientFunc defines a builder funciton used to build and return
// the SSH client to a Server. This is primarily used for
// mock unit testing.
// type sshClientFunc func(interface{}) (*sshClient, io.Closer, error)
type sshClientFunc func(*autoscaler.Server) (*sshClient, io.Closer, error)

// SSH client for unit testing
func testSSHClient(server *autoscaler.Server) (*sshClient, io.Closer, error) {
	var privateKeyFilePath = "/Users/niinai/.ssh/ttousai_rsa"
	key, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, nil, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, nil, err
	}

	config := &ssh.ClientConfig{
		User: "niinai",
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	_client, err := ssh.Dial("tcp", "173.255.202.60:22", config)
	if err != nil {
		return nil, nil, err
	}

	_sshClient := &sshClient{_client}
	return _sshClient, _sshClient, nil
}

// newSSHClient returns a new SSH client configured for the
// Server host and certificate chain.
func newSSHClient(server *autoscaler.Server) (*sshClient, io.Closer, error) {
	// TODO: get privateKeyFilePath from autoscaler config
	// in the meantime mount private key file to /keys/id_rsa
	var privateKeyFilePath = "/keys/id_rsa"
	key, err := ioutil.ReadFile(privateKeyFilePath)
	if err != nil {
		return nil, nil, err
	}

	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, nil, err
	}

	// TODO: DO NOT USE INSECUREIGNOREHOSTKEY
	// TODO: get User from autoscaler config.
	// in the meantime hardcode to ubuntu.
	config := &ssh.ClientConfig{
		User: "ubuntu",
		Auth: []ssh.AuthMethod{
			// Use the PublicKeys method for remote authentication.
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// TODO: get sshPort from autoscaler config
	// Connect to the remote server and perform the SSH handshake.
	_client, err := ssh.Dial("tcp", server.Address+":22", config)
	if err != nil {
		return nil, nil, err
	}

	_sshClient := &sshClient{_client}
	return _sshClient, _sshClient, nil
}

// Wrapper session
func (client *sshClient) Run(cmd string) (string, error) {
	// create a session and call Run then close the session
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	var stdoutBuf bytes.Buffer
	var stderrBuf bytes.Buffer
	session.Stdout = &stdoutBuf
	session.Stderr = &stderrBuf
	err = session.Run(cmd)
	if err != nil {
		return stderrBuf.String(), err
	}

	return stdoutBuf.String(), err
}

// Ping
func (client *sshClient) Ping() (string, error) {
	// out, err := client.Run("ps -C sshd -ocomm=")
	out, err := client.Run("ps -C drone-runner-exec -ocomm=")
	if err != nil {
		return out, err
	}
	return out, nil
}
