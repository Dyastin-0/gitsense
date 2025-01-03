package ssh

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"net"

	sshUtil "golang.org/x/crypto/ssh"
)

func RunCommand(privateKey string, instanceIP string, hostKey string, user string, commands string) (string, string, error) {
	if privateKey == "" || instanceIP == "" || user == "" || commands == "" {
		return "", "", fmt.Errorf("missing required parameters")
	}

	signer, err := sshUtil.ParsePrivateKey([]byte(privateKey))
	if err != nil {
		return "", "", fmt.Errorf("failed to parse private key: %w", err)
	}

	config := &sshUtil.ClientConfig{
		User: user,
		Auth: []sshUtil.AuthMethod{
			sshUtil.PublicKeys(signer),
		},
		HostKeyCallback: verifyHostKey(hostKey),
	}

	client, err := sshUtil.Dial("tcp", instanceIP+":22", config)
	if err != nil {
		return "", "", fmt.Errorf("failed to dial: %w", err)
	}

	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "", "", fmt.Errorf("failed to create session: %w", err)
	}

	defer session.Close()

	var stdout, stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	err = session.Run(fmt.Sprintf("/bin/bash -c \"%s\"", commands))
	if err != nil {
		return "", "", fmt.Errorf("failed to run command: %w", err)
	}

	return stdout.String(), stderr.String(), nil
}

func verifyHostKey(hostKey string) func(hostname string, remote net.Addr, key sshUtil.PublicKey) error {
	decodedHostKey, err := base64.StdEncoding.DecodeString(hostKey)
	if err != nil {
		return func(hostname string, remote net.Addr, key sshUtil.PublicKey) error {
			return fmt.Errorf("failed to decode host key: %w", err)
		}
	}

	return func(hostname string, remote net.Addr, key sshUtil.PublicKey) error {
		if bytes.Equal(decodedHostKey, key.Marshal()) {
			return nil
		}

		return fmt.Errorf("host key mismatch: expected %s, got %s", hostKey, base64.StdEncoding.EncodeToString(key.Marshal()))
	}
}
