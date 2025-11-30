package main

import (
	"context"
	"log"
)

type SayHelloPort interface {
	SayHello(ctx context.Context, userID int) (string, error)
}

type GettingUserNameRepository interface {
	GetUserName(ctx context.Context, userID int) (string, error)
}

type HelloService struct {
	gettingUserNameRepo GettingUserNameRepository
}

func NewHelloService(gettingUserNameRepo GettingUserNameRepository) *HelloService {
	return &HelloService{
		gettingUserNameRepo: gettingUserNameRepo,
	}
}

func (s *HelloService) SayHello(ctx context.Context, userID int) (string, error) {
	userName, err := s.gettingUserNameRepo.GetUserName(ctx, userID)
	if err != nil {
		return "", err
	}
	return s.makeHelloMessage(userName), nil
}

func (s *HelloService) makeHelloMessage(userName string) string {
	return "Hello, " + userName + "!"
}

// Log Adapter for driven port

type LogAdapterForGettingUserNameRepository struct {
	repo GettingUserNameRepository
}

func NewLogAdapterForGettingUserNameRepository(repo GettingUserNameRepository) *LogAdapterForGettingUserNameRepository {
	return &LogAdapterForGettingUserNameRepository{
		repo: repo,
	}
}

func (l *LogAdapterForGettingUserNameRepository) GetUserName(ctx context.Context, userID int) (string, error) {
	log.Println("Getting user name for userID:", userID)
	userName, err := l.repo.GetUserName(ctx, userID)
	if err != nil {
		// Here we would log the error
		log.Printf("Error getting user name for userID %d: %v", userID, err)
		return "", err
	}
	// Here we would log the successful retrieval
	log.Printf("Successfully got user name for userID %d: %s", userID, userName)
	return userName, nil
}
