package session

import (
	"parentscontactform/internal/models"
	"sync"
)

var (
	sessionStore = make(map[string]*models.SessionData)
	sessionMutex sync.RWMutex
)

func Get(sessionID string) (*models.SessionData, bool) {
	sessionMutex.RLock()
	defer sessionMutex.RUnlock()
	sessionData, exists := sessionStore[sessionID]
	return sessionData, exists
}

func Set(sessionID string, sessionData *models.SessionData) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	sessionStore[sessionID] = sessionData
}

func Delete(sessionID string) {
	sessionMutex.Lock()
	defer sessionMutex.Unlock()
	delete(sessionStore, sessionID)
}
