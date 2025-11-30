package main

import (
	"context"
	"errors"
	"testing"
)

type MockGettingUserNameRepository struct {
	err error
}

func (m *MockGettingUserNameRepository) GetUserName(ctx context.Context, userID int) (string, error) {
	if m.err != nil {
		return "", m.err
	}
	if userID == 1 {
		return "Alice", nil
	}
	return "Unknown", nil
}

func TestSayHello(t *testing.T) {
	ctx := context.Background()

	helloService := NewHelloService(&MockGettingUserNameRepository{})

	result, err := helloService.SayHello(ctx, 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := "Hello, Alice!"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}

func TestSayHelloWhenUserNotFound(t *testing.T) {
	ctx := context.Background()

	helloService := NewHelloService(&MockGettingUserNameRepository{
		err: errors.New("user not found"),
	})

	_, err := helloService.SayHello(ctx, 2)
	if err == nil {
		t.Errorf("expected error for unknown user, got nil")
	}

	if err.Error() != "user not found" {
		t.Errorf("expected 'user not found' error, got %s", err.Error())
	}
}

func TestLogAdapterForGettingUserNameRepository(t *testing.T) {
	ctx := context.Background()

	mockRepo := &MockGettingUserNameRepository{}
	logAdapter := NewLogAdapterForGettingUserNameRepository(mockRepo)

	helloService := NewHelloService(logAdapter)

	result, err := helloService.SayHello(ctx, 1)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	expected := "Hello, Alice!"
	if result != expected {
		t.Errorf("expected %s, got %s", expected, result)
	}
}
