package sessions

//go:generate msgp

// Default flashes key.
const flashesKey = "_flash"

type Session struct {
	ID       string                 `msg:"-"`
	Values   map[string]interface{} `msg:"v"`
	Options  *Options               `msg:"-"`
	IsNew    bool                   `msg:"-"`
	name     string                 `msg:"-"`
	isWriten bool                   `msg:"-"`
}

func newSession(name string, options *Options) *Session {
	return &Session{
		Options: options,
		IsNew:   true,
		name:    name,
	}
}

func (s *Session) Name() string {
	return s.name
}

func (s *Session) Set(key string, v interface{}) {
	s.Values[key] = v
	s.isWriten = true
}

func (s *Session) SetAll(vals map[string]interface{}) {
	s.Values = vals
	s.isWriten = true
}

//GetString return string value
func (s *Session) GetString(key string, def string) string {
	v, ok := s.Values[key]
	if !ok {
		return def
	}

	return v.(string)
}

//GetStrings return array string value
func (s *Session) GetStrings(key string, def []string) []string {
	v, ok := s.Values[key]
	if !ok {
		return def
	}

	valinf := v.([]interface{})
	vals := make([]string, len(valinf))
	for i := 0; i < len(valinf); i++ {
		vals[i] = valinf[i].(string)
	}

	return vals
}

//GetInt return int64 value
func (s *Session) GetInt(key string, def int64) int64 {
	v, ok := s.Values[key]
	if !ok {
		return def
	}

	return v.(int64)
}

//GetInts return array int64 value
func (s *Session) GetInts(key string, def []int64) []int64 {
	v, ok := s.Values[key]
	if !ok {
		return def
	}

	valinf := v.([]interface{})
	vals := make([]int64, len(valinf))
	for i := 0; i < len(valinf); i++ {
		vals[i] = valinf[i].(int64)
	}

	return vals
}

//GetFloat return float64 value
func (s *Session) GetFloat(key string, def float64) float64 {
	v, ok := s.Values[key]
	if !ok {
		return def
	}

	return v.(float64)
}

//GetFloats return array float64 value
func (s *Session) GetFloats(key string, def []float64) []float64 {
	v, ok := s.Values[key]
	if !ok {
		return def
	}

	valinf := v.([]interface{})
	vals := make([]float64, len(valinf))
	for i := 0; i < len(valinf); i++ {
		vals[i] = valinf[i].(float64)
	}

	return vals
}

//GetBool return bool value
func (s *Session) GetBool(key string, def bool) bool {
	v, ok := s.Values[key]
	if !ok {
		return def
	}

	return v.(bool)
}

// Flashes returns a slice of flash messages from the session.
//
// A single variadic argument is accepted, and it is optional: it defines
// the flash key. If not defined "_flash" is used by default.
func (s *Session) Flashes() []string {
	if v, ok := s.Values[flashesKey]; ok {
		delete(s.Values, flashesKey)
		s.isWriten = true
		valinf := v.([]interface{})
		vals := make([]string, len(valinf))
		for i := 0; i < len(valinf); i++ {
			vals[i] = valinf[i].(string)
		}

		return vals
	}
	return nil
}

// AddFlash adds a flash message to the session.
//
// A single variadic argument is accepted, and it is optional: it defines
// the flash key. If not defined "_flash" is used by default.
func (s *Session) AddFlash(value string) {
	if v, ok := s.Values[flashesKey]; ok {
		flashes := v.([]string)
		s.Values[flashesKey] = append(flashes, value)
	} else {
		s.Values[flashesKey] = []string{value}
	}

	s.isWriten = true
}

func (s *Session) AddFlashs(values []string) {
	if v, ok := s.Values[flashesKey]; ok {
		flashes := v.([]string)
		s.Values[flashesKey] = append(flashes, values...)
	} else {
		s.Values[flashesKey] = values
	}

	s.isWriten = true
}
