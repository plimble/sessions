package sessions

import (
	"github.com/oxtoacart/bpool"
	"github.com/tinylib/msgp/msgp"
	"sync"

	"net/http"
)

type Sessions struct {
	sync.RWMutex
	req      *http.Request
	store    Store
	options  *Options
	sessions map[string]*Session
	bufPool  *bpool.BufferPool
}

func (s *Sessions) Get(name string) (*Session, error) {
	if session, ok := s.sessions[name]; ok {
		return session, nil
	}

	s.Lock()
	defer s.Unlock()
	s.sessions[name] = newSession(name, s.options)

	var err error
	var c *http.Cookie

	if c, err = s.req.Cookie(name); err != nil {
		s.sessions[name].ID = generateUUID()
		s.sessions[name].Values = make(map[string]interface{})
		return s.sessions[name], nil
	}

	buf := s.bufPool.Get()
	defer s.bufPool.Put(buf)

	if err = s.store.Get(c.Value, buf); err == nil {

		if err = msgp.Decode(buf, s.sessions[name]); err == nil {
			s.sessions[name].ID = c.Value
			s.sessions[name].IsNew = false
			return s.sessions[name], nil
		}

	}

	s.sessions[name].ID = generateUUID()
	s.sessions[name].Values = make(map[string]interface{})

	return s.sessions[name], err
}

func (s *Sessions) Delete(name string) error {
	s.Lock()
	defer s.Unlock()

	session, ok := s.sessions[name]
	if !ok {
		return nil
	}

	if err := s.store.Delete(session.ID); err != nil {
		return err
	}

	delete(s.sessions, name)

	return nil
}

func (s *Sessions) Save(w http.ResponseWriter) error {
	buf := s.bufPool.Get()
	defer s.bufPool.Put(buf)
	for _, session := range s.sessions {
		if !session.isWriten {
			continue
		}

		if err := msgp.Encode(buf, session); err != nil {
			return err
		}

		if err := s.store.Save(session, buf, w); err != nil {
			return err
		}

		buf.Reset()
	}

	return nil
}
