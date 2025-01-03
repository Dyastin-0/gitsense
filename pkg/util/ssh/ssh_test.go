package ssh

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	err := godotenv.Load("../../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	code := m.Run()

	os.Exit(code)
}

func TestRunCommand(t *testing.T) {
	privateKey := os.Getenv("PRIVATE_KEY")

	hostKey := os.Getenv("HOST_KEY")
	instanceIP := os.Getenv("INSTANCE_IP")
	user := os.Getenv("USER")
	command := `echo Hello, SSH!`

	stdout, stderr, err := RunCommand(privateKey, instanceIP, hostKey, user, command)

	fmt.Println(stdout)

	assert.NoError(t, err, "Expected no error when running command")
	assert.Equal(t, "Hello, SSH!\n", stdout, "Expected output from command")
	assert.Empty(t, stderr, "Expected no stderr output")
}
