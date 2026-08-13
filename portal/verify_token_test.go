package portal_test

import (
	"testing"

	"github.com/nvvers/orbit-portal-client-golang/portal"
	"github.com/nvvers/orbit-portal-client-golang/portaltestcontainer"
)

func TestClient_VerifyToken(t *testing.T) {
	pc, err := portaltestcontainer.New()
	if err != nil {
		t.Fatalf("Failed to create test container: %v", err)
	}

	err = pc.Start(t.Context())
	if err != nil {
		t.Fatalf("Failed to start test container: %v", err)
	}
	defer pc.Stop(t.Context())

	client, err := pc.GetClient(t.Context())
	if err != nil {
		t.Fatalf("Failed to get client: %v", err)
	}

	t.Run("verify token", func(t *testing.T) {
		err := client.VerifyToken()
		if err != nil {
			t.Fatalf("Failed to notify event: %v", err)
		}
	})

	t.Run("verify invalid token", func(t *testing.T) {
		clientWithInvalidToken, err := pc.GetClient(t.Context())
		if err != nil {
			t.Fatalf("Failed to set invalid token: %v", err)
		}

		err = portal.WithToken("invalid-token")(clientWithInvalidToken)
		if err != nil {
			t.Fatalf("Failed to set invalid token: %v", err)
		}

		err = clientWithInvalidToken.VerifyToken()
		if err == nil {
			t.Fatalf("Expected error when verifying invalid token, but got none")
		}
	})
}
