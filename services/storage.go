package services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

type Session struct {
	ID              string
	Filename        string
	DocumentContent string
	CreatedAt       time.Time
	ExpiresAt       time.Time
}

type Storage struct {
	sessions map[string]*Session
	mutex    sync.RWMutex
}

var globalStorage *Storage

func InitStorage() {
	globalStorage = &Storage{
		sessions: make(map[string]*Session),
	}
	
	// Start cleanup goroutine
	go globalStorage.cleanupExpiredSessions()
}

func GetStorage() *Storage {
	return globalStorage
}

func (s *Storage) CreateSession(filename, content string) (*Session, error) {
	sessionID, err := generateSessionID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate session ID: %w", err)
	}

	session := &Session{
		ID:              sessionID,
		Filename:        filename,
		DocumentContent: content,
		CreatedAt:       time.Now(),
		ExpiresAt:       time.Now().Add(30 * time.Minute),
	}

	s.mutex.Lock()
	s.sessions[sessionID] = session
	s.mutex.Unlock()

	return session, nil
}

func (s *Storage) GetSession(sessionID string) (*Session, error) {
	s.mutex.RLock()
	session, exists := s.sessions[sessionID]
	s.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("session not found")
	}

	if time.Now().After(session.ExpiresAt) {
		s.DeleteSession(sessionID)
		return nil, fmt.Errorf("session expired")
	}

	return session, nil
}

func (s *Storage) DeleteSession(sessionID string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if _, exists := s.sessions[sessionID]; !exists {
		return fmt.Errorf("session not found")
	}
	
	delete(s.sessions, sessionID)
	return nil
}

func (s *Storage) cleanupExpiredSessions() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mutex.Lock()
		now := time.Now()
		
		for sessionID, session := range s.sessions {
			if now.After(session.ExpiresAt) {
				delete(s.sessions, sessionID)
			}
		}
		s.mutex.Unlock()
	}
}

func (s *Storage) GetSessionCount() int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return len(s.sessions)
}

func generateSessionID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// Convenience functions for global storage
func CreateSession(filename, content string) (*Session, error) {
	return globalStorage.CreateSession(filename, content)
}

func GetSession(sessionID string) (*Session, error) {
	return globalStorage.GetSession(sessionID)
}

func DeleteSession(sessionID string) error {
	return globalStorage.DeleteSession(sessionID)
}