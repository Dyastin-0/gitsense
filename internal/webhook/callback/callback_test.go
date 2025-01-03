package callback

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	code := m.Run()

	os.Exit(code)
}

func TestHandler(t *testing.T) {
	url := "http://localhost:4000/api/v1/callback/Dyastin-0/cdmlms/hooks/test"

	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	fmt.Println(string(body))

	assert.Equal(t, 200, resp.StatusCode, "Expected status code 200")
}
