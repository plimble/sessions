package mongo

import (
	"bytes"
	"github.com/plimble/sessions"
	"gopkg.in/mgo.v2"
	"net/http"
)

type SessionData struct {
	Id   string `bson:"_id"`
	Data []byte `bson:"data"`
}

type MongoStore struct {
	session    *mgo.Session
	db         string
	collection string
}

func NewMongoStore(session *mgo.Session, db, collection string) *MongoStore {
	return &MongoStore{session, db, collection}
}

func (s *MongoStore) Get(id string, buf *bytes.Buffer) error {
	data := &SessionData{}
	if err := s.session.DB(s.db).C(s.collection).FindId(id).One(&data); err != nil {
		return err
	}

	if _, err := buf.Write(data.Data); err != nil {
		return err
	}

	return nil
}

func (s *MongoStore) Save(session *sessions.Session, buf *bytes.Buffer, w http.ResponseWriter) error {
	data := &SessionData{
		Id:   session.ID,
		Data: buf.Bytes(),
	}

	if _, err := s.session.DB(s.db).C(s.collection).UpsertId(session.ID, data); err != nil {
		return err
	}

	http.SetCookie(w, sessions.NewCookie(session.Name(), session.ID, session.Options))
	return nil
}
