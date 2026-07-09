package session

import (
	"sync/atomic"

	"github.com/ajsqr/wombat/cmap"
	"github.com/ajsqr/wombat/receiver"
)

type SessionStore struct {
	counter atomic.Uint32
	store   *cmap.Cmap[uint32, *Session]
}

func NewSessionStore() *SessionStore {
	data := cmap.New[uint32, *Session]()
	return &SessionStore{
		store:   data,
		counter: atomic.Uint32{},
	}
}

func (ss *SessionStore) Create(s *Session) {
	ss.store.Set(s.sessionID, s)
}

func (ss *SessionStore) GetByID(id uint32) (*Session, error) {
	session, ok := ss.store.Get(id)
	if !ok {
		return nil, receiver.ErrSessionNotFound
	}

	return session, nil
}

func (ss *SessionStore) Destroy(id uint32) {
	s, ok := ss.store.Get(id)
	if !ok {
		return
	}

	close(s.inbox)
	ss.store.Delete(id)
}
