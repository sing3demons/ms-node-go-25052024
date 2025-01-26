package store

import (
	"testing"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestStartSession(t *testing.T) {
	mockClient := new(MockMongoClient)

	// Create a SessionOptions pointer (not a slice)
	sessionOptions := &options.SessionOptions{}

	// Mock StartSession to expect the specific SessionOptions pointer
	mockClient.On("StartSession", sessionOptions).Return(&MockSession{}, nil)

	// Call StartSession with the correct argument
	session, err := mockClient.StartSession(sessionOptions)

	// Validate the error is nil
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	// Validate the session object
	if session == nil {
		t.Fatal("expected non-nil session")
	}

	// Assert that the mock expectations were met
	mockClient.AssertExpectations(t)
}
