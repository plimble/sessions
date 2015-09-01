package sessions

import (
	"bytes"
	"net/http"
	"sync"
)

type Store interface {
	Get(id string, buf *bytes.Buffer) error
	Save(session *Session, buf *bytes.Buffer, w http.ResponseWriter) error
	Delete(session *Session, w http.ResponseWriter) error
}

type MemoryStore struct {
	sync.RWMutex
	data map[string][]byte
}

type memoryError struct {
	text string
}

func newMemoryError(text string) *memoryError {
	return &memoryError{text}
}

func (m *memoryError) Error() string {
	return m.text
}

func NewMemoryStore() *MemoryStore {
	s := &MemoryStore{
		data: make(map[string][]byte),
	}

	return s
}

func (s *MemoryStore) Get(id string, buf *bytes.Buffer) error {
	b, ok := s.data[id]
	if !ok {
		return newMemoryError("not found")
	}

	_, err := buf.Write(b)

	return err
}

func (s *MemoryStore) Delete(id string) error {
	s.Lock()
	delete(s.data, id)
	s.Unlock()

	return nil
}

func (s *MemoryStore) Save(session *Session, buf *bytes.Buffer, w http.ResponseWriter) error {
	s.Lock()
	s.data[session.ID] = buf.Bytes()
	s.Unlock()

	http.SetCookie(w, NewCookie(session.Name(), session.ID, session.Options))
	return nil
}
