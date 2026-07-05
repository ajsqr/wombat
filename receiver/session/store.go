package session

import (
	"net"
	"sync/atomic"

	"github.com/ajsqr/wombat/cmap"
	"github.com/ajsqr/wombat/frame"
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

func (ss *SessionStore) New(conn net.Conn) *Session {
	var s Session
	s.sessionID = ss.counter.Add(1)
	s.inbox = make(chan *frame.Frame)
	s.conn = conn
	// once the session has been crafted, add it to the in-memory store
	ss.store.Set(s.sessionID, &s)
	return &s
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
