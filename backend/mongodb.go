package main

import (
	"gopkg.in/mgo.v2/bson"
	"log"
	"os"
)

var sessions = make([]*Session, 0, MAX)

const MAX = 99

func postTrack(track Track) {
	s := getSession()
	s.Locked = true
	coll := s.getCollection()
	err := coll.Insert(&Track{track.Id, track.Name, track.Author, track.Steps})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("postTrack session UUID:", s.Uuid)
	s.Locked = false
}

func getTracks() []Track {
	results := []Track{}
	s := getSession()
	s.Locked = true
	coll := s.getCollection()
	err := coll.Find(bson.M{}).All(&results)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("getTracks session UUID:", s.Uuid)
	s.Locked = false
	return results
}

func getTrack(id int) Track {
	result := Track{}
	s := getSession()
	s.Locked = true
	coll := s.getCollection()
	err := coll.Find(bson.M{"id": id}).One(&result)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("getTrack session UUID:", s.Uuid)
	s.Locked = false
	return result
}

func getSession() *Session {
	// get next available session
	countSessions := 0
	countLockedSessions := 0
	addSession := true
	for _, s := range sessions {
		if s != nil {
			if s.Locked {
				countLockedSessions++
			} else {
				addSession = false
				break
			}
			countSessions++
			addSession = true
		}
	}
	if addSession && countSessions > MAX {
		log.Fatal("No more connection available ", countSessions, "/", MAX, " reached")
	}
	if addSession && countSessions <= MAX {
		s := New(os.Getenv("MONGODB_HOST"), "stl", "stluser", os.Getenv("MONGODB_ROOT_PASSWORD"), "stl", "tracks")
		sessions = append(sessions, s)
		log.Println("New Session created: ", s.Uuid)
	} else {
		log.Println("Use existing connection")
		sessions[countSessions].MgoSession.Refresh()
	}
	return sessions[countSessions]
}

func closeSession() {
	for _, s := range sessions {
		if s != nil {
			defer s.MgoSession.Close()
		}
	}
}
