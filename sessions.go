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
	s.sessions[name] = newSession(name, s.options)

	c, err := s.req.Cookie(name)
	if err != nil {
		s.sessions[name].Values = make(map[string]interface{})
		return s.sessions[name], err
	}

	buf := s.bufPool.Get()
	defer s.bufPool.Put(buf)

	if err := s.store.Get(c.Value, buf); err != nil {
		s.sessions[name].Values = make(map[string]interface{})
		return s.sessions[name], err
	}

	s.sessions[name].DecodeMsg(msgp.NewReader(buf))
	s.sessions[name].IsNew = false

	return s.sessions[name], nil
}

func (s *Sessions) Save(w http.ResponseWriter) error {
	var err error
	var b []byte
	buf := s.bufPool.Get()
	defer s.bufPool.Put(buf)
	for _, session := range s.sessions {
		if !session.isWriten {
			continue
		}

		session.ID = generateUUID()

		b, err = session.MarshalMsg(nil)
		if err != nil {
			return err
		}

		buf.Write(b)

		if err = s.store.Save(session, buf, w); err != nil {
			return err
		}

		buf.Reset()
	}

	return nil
}
