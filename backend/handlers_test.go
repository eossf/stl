package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/julienschmidt/httprouter"
)

func TestTrackShow(t *testing.T) {
	t.Log("When the Tracks' isdn does not exist")
	// A request with a non-existant isdn
	req1, err := http.NewRequest("GET", "/tracks/1234", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr1 := newRequestRecorder(req1, "GET", "/tracks/:isdn", TrackShow)
	if rr1.Code != 404 {
		t.Error("Expected response code to be 404")
	}
	// expected response
	er1 := "{\"error\":{\"status\":404,\"title\":\"Record Not Found\"}}\n"
	if rr1.Body.String() != er1 {
		t.Error("Response body does not match")
	}

	t.Log("When the Track exists")
	// Create an entry of the Track to the trackstore map
	testTrack := &Track{
		Id:   "111",
		Name:  "test title",
		Author: "test author",
		Steps:  42,
	}
	trackstore["111"] = testTrack
	// A request with an existing isdn
	req2, err := http.NewRequest("GET", "/tracks/111", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr2 := newRequestRecorder(req2, "GET", "/tracks/:isdn", TrackShow)
	if rr2.Code != 200 {
		t.Error("Expected response code to be 200")
	}
	// expected response
	er2 := "{\"meta\":null,\"data\":{\"isdn\":\"111\",\"title\":\"test title\",\"author\":\"test author\",\"pages\":42}}\n"
	if rr2.Body.String() != er2 {
		t.Error("Response body does not match")
	}
}

func TestTrackIndex(t *testing.T) {
	// Create an entry of the Track to the trackstore map
	testTrack := &Track{
		Id:   "111",
		Name:  "test title",
		Author: "test author",
		Steps:  42,
	}
	trackstore["111"] = testTrack
	// A request with an existing isdn
	req1, err := http.NewRequest("GET", "/tracks", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr1 := newRequestRecorder(req1, "GET", "/tracks", TrackIndex)
	if rr1.Code != 200 {
		t.Error("Expected response code to be 200")
	}
	// expected response
	er1 := "{\"meta\":null,\"data\":[{\"id\":\"111\",\"name\":\"test title\",\"author\":\"test author\",\"steps\":42}]}\n"
	if rr1.Body.String() != er1 {
		t.Error("Response body does not match")
	}
}

// Mocks a handler and returns a httptest.ResponseRecorder
func newRequestRecorder(req *http.Request, method string, strPath string, fnHandler func(w http.ResponseWriter, r *http.Request, param httprouter.Params)) *httptest.ResponseRecorder {
	router := httprouter.New()
	router.Handle(method, strPath, fnHandler)
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router.ServeHTTP(rr, req)
	return rr
}
