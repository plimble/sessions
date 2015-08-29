package sessions

import (
	"github.com/oxtoacart/bpool"
	"net/http"
	"sync"
)

type Options struct {
	Path   string
	Domain string
	// MaxAge=0 means no 'Max-Age' attribute specified.
	// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'.
	// MaxAge>0 means Max-Age attribute present and given in seconds.
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

func (o *Options) mergeDefault() {
	if o.Path == "" {
		o.Path = "/"
	}
}

type SessionManager struct {
	pool    sync.Pool
	bufPool *bpool.BufferPool
}

func New(poolSize int, store Store, options *Options) *SessionManager {
	if poolSize == 0 {
		poolSize = 10000
	}

	if options == nil {
		options = &Options{}
	}

	options.mergeDefault()

	s := &SessionManager{
		bufPool: bpool.NewBufferPool(poolSize),
	}

	s.pool.New = func() interface{} {
		opts := *options
		return &Sessions{
			sessions: make(map[string]*Session),
			options:  &opts,
			store:    store,
			bufPool:  s.bufPool,
		}
	}

	return s
}

func (s *SessionManager) GetSessions(r *http.Request) *Sessions {
	sessions := s.pool.Get().(*Sessions)
	sessions.req = r
	return sessions
}

func (s *SessionManager) Close(sessions *Sessions) {
	for key := range sessions.sessions {
		delete(sessions.sessions, key)
	}
	s.pool.Put(sessions)
}
