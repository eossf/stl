package main

type Track struct {
	Id   string `json:"id"`
	Name  string `json:"name"`
	Author string `json:"author"`
	Steps  int    `json:"steps"`
}

// TODO: mongodb
var trackstore = make(map[string]*Track)
