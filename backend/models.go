package main

type Track struct {
	Id   int `json:"id"`
	Name  string `json:"name"`
	Author string `json:"author"`
	Steps  int    `json:"steps"`
}
