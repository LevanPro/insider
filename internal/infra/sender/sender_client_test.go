package sender_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/LevanPro/insider/internal/infra/sender"
)

type MockSendResponse struct {
	MessageID string
}

// Since the original code uses the imported type, we alias the mock
// to satisfy the signature of the tested Client.Send method.
func init() {
	// This is a common pattern to ensure the imported service package is used
	// if available. If the actual service.SendResponse is used, this mock
	// struct definition should be removed.
	// For this isolated test environment, we pretend the imported type exists.
}

func TestClient_Send(t *testing.T) {
	const (
		testTo      = "test-recipient@example.com"
		testContent = "Hello, world!"
		testAPIKey  = "test-api-key"
		expectedID  = "msg-12345"
	)

	tests := []struct {
		name           string
		apiKey         string
		mockStatusCode int
		mockBody       []byte
		handler        http.HandlerFunc
		expectError    bool
		expectedID     string
	}{
		{
			name:           "Success_Status_200_With_API_Key",
			apiKey:         testAPIKey,
			mockStatusCode: http.StatusOK,
			mockBody:       []byte(fmt.Sprintf(`{"message": "Sent", "messageId": "%s"}`, expectedID)),
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("Expected POST method, got %s", r.Method)
				}
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("Expected Content-Type: application/json, got %s", r.Header.Get("Content-Type"))
				}
				if r.Header.Get("x-ins-auth-key") != testAPIKey {
					t.Errorf("Expected x-ins-auth-key: %s, got %s", testAPIKey, r.Header.Get("x-ins-auth-key"))
				}
				// 3. Check body payload
				body, _ := io.ReadAll(r.Body)
				var payload struct {
					To      string `json:"to"`
					Content string `json:"content"`
				}
				if err := json.Unmarshal(body, &payload); err != nil {
					t.Fatalf("Could not unmarshal request body: %v", err)
				}
				if payload.To != testTo || payload.Content != testContent {
					t.Errorf("Expected body {To: %s, Content: %s}, got {To: %s, Content: %s}",
						testTo, testContent, payload.To, payload.Content)
				}

				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(fmt.Sprintf(`{"message": "Sent", "messageId": "%s"}`, expectedID)))
			},
			expectError: false,
			expectedID:  expectedID,
		},
		{
			name:           "Success_Status_202_No_API_Key",
			apiKey:         "", // No API key
			mockStatusCode: http.StatusAccepted,
			mockBody:       []byte(fmt.Sprintf(`{"message": "Accepted", "messageId": "%s"}`, expectedID)),
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("x-ins-auth-key") != "" {
					t.Errorf("Did not expect x-ins-auth-key header, got %s", r.Header.Get("x-ins-auth-key"))
				}
				w.WriteHeader(http.StatusAccepted)
				_, _ = w.Write([]byte(fmt.Sprintf(`{"message": "Accepted", "messageId": "%s"}`, expectedID)))
			},
			expectError: false,
			expectedID:  expectedID,
		},
		{
			name:           "Failure_Status_400",
			apiKey:         testAPIKey,
			mockStatusCode: http.StatusBadRequest,
			mockBody:       []byte(`{"error": "Invalid input"}`),
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusBadRequest)
				_, _ = w.Write([]byte(`{"error": "Invalid input"}`))
			},
			expectError: true,
			expectedID:  "",
		},
		{
			name:           "Failure_Invalid_JSON_Response",
			apiKey:         testAPIKey,
			mockStatusCode: http.StatusOK,
			mockBody:       []byte(`This is not JSON`),
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`This is not JSON`))
			},
			expectError: true,
			expectedID:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(tt.handler)
			defer server.Close()

			client := sender.NewClient(server.URL, tt.apiKey)

			resp, err := client.Send(context.Background(), testTo, testContent)

			if tt.expectError {
				if err == nil {
					t.Fatalf("Expected an error, but got nil response: %+v", resp)
				}
				if err.Error() == "" {
					t.Errorf("Expected an informative error message, got empty string")
				}
				return
			}

			// Check for expected success state
			if err != nil {
				t.Fatalf("Expected no error, but got: %v", err)
			}

			// Check the returned MessageID
			if resp.MessageID != tt.expectedID {
				t.Errorf("Expected MessageID %s, got %s", tt.expectedID, resp.MessageID)
			}
		})
	}
}

func TestClient_Send_Timeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(10 * time.Second) // Longer than client's 5s timeout
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"message": "Sent", "messageId": "timeout-test"}`))
	}))
	defer server.Close()

	// The Client's default timeout is 5 seconds.
	client := sender.NewClient(server.URL, "any-key")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := client.Send(ctx, "test", "content")

	if err == nil {
		t.Fatal("Expected a timeout error, but call succeeded")
	}

	if err.Error() == "" {
		t.Errorf("Expected an informative error message for timeout, got empty string")
	}
}
