package e2e

import (
	"context"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"io"
	"net/http"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	uri := buildServerContainer(t)

	t.Run("Empty result", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("%s/Test/changes", uri))
		if err != nil {
			t.Fatalf("%v", err)
		}
		content, _ := io.ReadAll(resp.Body)
		contentStr := strings.TrimSpace(string(content))
		if contentStr != "[]" {
			t.Fatalf("Expected [], got %s", contentStr)
		}
	})

}

func buildServerContainer(t *testing.T) string {
	ctx := context.Background()
	request := testcontainers.ContainerRequest{
		FromDockerfile: testcontainers.FromDockerfile{
			Context:    "../server",
			Dockerfile: "Dockerfile",
			Repo:       "logsync",
			Tag:        "latest",
		},
		ExposedPorts: []string{"3000/tcp"},
		WaitingFor:   wait.ForListeningPort("3000"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: request,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Cleanup(func() {
		_ = container.Terminate(ctx)
	})

	ip, err := container.Host(ctx)
	if err != nil {
		t.Fatalf("%v", err)
	}

	mappedPort, err := container.MappedPort(ctx, "3000")
	if err != nil {
		t.Fatalf("%v", err)
	}

	uri := fmt.Sprintf("http://%s:%s", ip, mappedPort.Port())
	return uri
}
