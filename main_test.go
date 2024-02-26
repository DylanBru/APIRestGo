package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// --------------Vérifie si l'appel à l'endpoint /api/v1/videos renvoie statut 200--------------//
func TestVideoHandler(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(videoHandler))
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("excpected 200 but got %d", resp.StatusCode)
	}
}

func TestListVideos(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/videos", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Enregistre la réponse http pour le test
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(videoHandler)
	handler.ServeHTTP(rr, req)
	// Vérification du code de statut de la réponse
	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %d", rr.Code)
	}
	expectedBytes, er := json.Marshal(videos)
	if er != nil {
		t.Fatal(er)
	}
	expected := string(expectedBytes)
	// Compare les deux représentations JSON après les avoir nettoyées
	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(string(expected)) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestGetVideo(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v1/videos/2", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Enregistre la réponse http pour le test
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(videoHandler)
	handler.ServeHTTP(rr, req)
	// Vérification du code de statut de la réponse
	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %d", rr.Code)
	}
	expectedBytes, er := json.Marshal(videos[1])
	if er != nil {
		t.Fatal(er)
	}
	expected := string(expectedBytes)
	// Compare les deux représentations JSON après les avoir nettoyées
	if strings.TrimSpace(rr.Body.String()) != strings.TrimSpace(string(expected)) {
		t.Errorf("handler returned unexpected body: got %v want %v", rr.Body.String(), expected)
	}
}

func TestCreateVideo(t *testing.T) {
	expected := video{ID: "4", Title: "Le quatrième", Author: "Dylan Bru", PublishedDate: "2024-02-22"}
	requestBody := `{"id": "4", "title": "Le quatrième", "author": "Dylan Bru", "publishedDate": "2024-02-22"}`
	req, err := http.NewRequest("POST", "/api/v1/videos", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(videoHandler)
	handler.ServeHTTP(rr, req)
	// Vérification du code de statut de la réponse
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status Created; got %d", rr.Code)
	}
	// Récupérer la liste des vidéos à jour
	reqGet, er := http.NewRequest("GET", "/api/v1/videos", nil)
	if er != nil {
		t.Fatal(er)
	}
	reqGetRecorder := httptest.NewRecorder()
	handler.ServeHTTP(reqGetRecorder, reqGet)
	// Vérification si la vidéo ajoutée correspond aux données attendues
	var videosResponse []video
	if e := json.Unmarshal(reqGetRecorder.Body.Bytes(), &videosResponse); e != nil {
		t.Fatal(e)
	}
	found := false
	for _, v := range videosResponse {
		if v.ID == expected.ID && v.Title == expected.Title && v.Author == expected.Author && v.PublishedDate == expected.PublishedDate {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected video %+v not found in response body: %v", expected, videosResponse)
	}
}

func TestDeleteVideo(t *testing.T) {
	deleteExpected := videos[1]
	req, err := http.NewRequest("DELETE", "/api/v1/videos/2", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(videoHandler)
	handler.ServeHTTP(rr, req)
	// Vérification du code de statut de la réponse
	if rr.Code != http.StatusOK {
		t.Errorf("expected status Created; got %d", rr.Code)
	}
	// Récupérer la liste des vidéos à jour
	reqGet, er := http.NewRequest("GET", "/api/v1/videos", nil)
	if er != nil {
		t.Fatal(er)
	}
	reqGetRecorder := httptest.NewRecorder()
	handler.ServeHTTP(reqGetRecorder, reqGet)
	// Vérification si la vidéo a bien été supprimée
	var videosResponse []video
	if e := json.Unmarshal(reqGetRecorder.Body.Bytes(), &videosResponse); e != nil {
		t.Fatal(e)
	}
	found := false
	for _, v := range videosResponse {
		if v.ID == deleteExpected.ID && v.Title == deleteExpected.Title && v.Author == deleteExpected.Author && v.PublishedDate == deleteExpected.PublishedDate {
			found = true
			break
		}
	}
	if found {
		t.Errorf("video %+v found in response body: %v, expected : not found", deleteExpected, videosResponse)
	}
}
