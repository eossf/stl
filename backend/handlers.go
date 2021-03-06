package main

import (
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome!\n")
}

// DELETE /track/:id
func TrackDelete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	i, _ := strconv.Atoi(id)
	result := delTrack(i)
	writeOKResponse(w, result)
}

// PUT /track
func TrackUpdate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	track := &Track{}
	if err := populateModelFromHandler(w, r, params, track); err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessible Entity")
		return
	}
	putTrack(*track)
	writeOKResponse(w, track)
}

// POST /tracks
func TrackCreate(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	track := &Track{}
	if err := populateModelFromHandler(w, r, params, track); err != nil {
		writeErrorResponse(w, http.StatusUnprocessableEntity, "Unprocessible Entity")
		return
	}
	postTrack(*track)
	writeOKResponse(w, track)
}

// GET /tracks
func TrackIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	tracks := []Track{}
	results := getTracks()
	for _, track := range results {
		tracks = append(tracks, track)
	}
	writeOKResponse(w, tracks)
}

// GET /tracks/:id
func TrackShow(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	id := params.ByName("id")
	i, _ := strconv.Atoi(id)
	result := getTrack(i)
	writeOKResponse(w, result)
}

// Writes the response as a standard JSON response with StatusOK
func writeOKResponse(w http.ResponseWriter, m interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.WriteHeader(http.StatusOK)
	v := json.NewEncoder(w).Encode(m)
	if err := v; err != nil {
		writeErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
	}
}

// Writes the error response as a Standard API JSON response with a response code
func writeErrorResponse(w http.ResponseWriter, errorCode int, errorMsg string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.WriteHeader(errorCode)
	json.
		NewEncoder(w).
		Encode(&JsonErrorResponse{Error: &ApiError{Status: errorCode, Title: errorMsg}})
}

//Populates a model from the params in the Handler
func populateModelFromHandler(w http.ResponseWriter, r *http.Request, params httprouter.Params, model interface{}) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return err
	}
	if err := r.Body.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(body, model); err != nil {
		return err
	}
	return nil
}
