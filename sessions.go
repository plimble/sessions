package sessions

import (
	"github.com/plimble/utils/pool"
	"github.com/tinylib/msgp/msgp"

	"net/http"
)

type Sessions struct {
	req      *http.Request
	store    Store
	options  *Options
	sessions map[string]*Session
	bufPool  *pool.BufferPool
}

func (s *Sessions) Get(name string) (*Session, error) {
	if session, ok := s.sessions[name]; ok {
		return session, nil
	}

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
