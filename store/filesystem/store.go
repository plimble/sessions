package filesystem

import (
	"bytes"
	"github.com/plimble/sessions"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

var fileMutex sync.RWMutex

type FileSystemStore struct {
	path string
}

func NewFileSystemStore(path string) *FileSystemStore {
	if path == "" {
		path = os.TempDir()
	}

	return &FileSystemStore{
		path: path,
	}
}

func (s *FileSystemStore) Get(id string, buf *bytes.Buffer) error {
	fileMutex.RLock()
	defer fileMutex.RUnlock()

	f, err := os.OpenFile(filepath.Join(s.path, "session_"+id), os.O_RDONLY, 0755)
	if err != nil {
		return err
	}

	_, err = buf.ReadFrom(f)

	return err
}

func (s *FileSystemStore) Save(session *sessions.Session, buf *bytes.Buffer, w http.ResponseWriter) error {
	fileMutex.RLock()
	defer fileMutex.RUnlock()

	f, err := os.OpenFile(filepath.Join(s.path, "session_"+session.ID), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0755)
	if err != nil {
		return err
	}

	_, err = io.Copy(f, buf)
	if err != nil {
		return err
	}

	http.SetCookie(w, sessions.NewCookie(session.Name(), session.ID, session.Options))

	return nil
}
