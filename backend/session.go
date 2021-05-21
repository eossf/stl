package main

import (
	"github.com/google/uuid"
	"gopkg.in/mgo.v2"
)

type Session struct {
	MgoHost       string
	DbName        string
	User          string
	Password      string
	SourceName    string
	CollName      string
	MgoSession    *mgo.Session
	MgoDb         *mgo.Database
	MgoCollection *mgo.Collection
	Err           error
	Uuid          string
	Locked        bool
}

func New(mgohost string, dbname string, user string, password string, sourcename string, collname string) *Session {
	s := Session{
		MgoHost:    mgohost,
		DbName:     dbname,
		User:       user,
		Password:   password,
		SourceName: sourcename,
		CollName:   collname,
		Uuid:       uuid.NewString(),
		Locked:     false,
	}
	s.MgoSession, s.Err = mgo.Dial(s.MgoHost)
	if s.Err != nil {
		panic(s.Err)
	}
	s.MgoSession.SetMode(mgo.Monotonic, true)
	s.MgoSession.Login(&mgo.Credential{Username: s.User, Password: s.Password, Source: s.SourceName})
	return &s
}

func (s Session) getCollection() *mgo.Collection {
	s.MgoDb = s.MgoSession.DB(s.DbName)
	s.MgoCollection = s.MgoDb.C(s.CollName)
	return s.MgoCollection
}
