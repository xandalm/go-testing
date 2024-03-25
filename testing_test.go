package testing_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	tpkg "github.com/xandalm/go-testing"
)

func TestServerLauncher(t *testing.T) {
	ctx := context.Background()
	launcher := tpkg.NewServerLauncher(
		ctx,
		"cmd/",
		"main.go",
		&tpkg.HTTPServerChecker{
			"http://localhost:5000",
			&http.Client{},
		},
	)

	if err := launcher.StartAndWait(5 * time.Second); err != nil {
		t.Fatalf("cannot start the server, %v", err)
	}

	time.Sleep(2 * time.Second)

	if err := launcher.EndAndClean(); err != nil {
		t.Errorf("cannot graceful shutdown the server, %v", err)
	}
}
