package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestListVideos(t *testing.T) {
	req, err := http.NewRequest("GET", "/api/v2/videos", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Enregistre la réponse http pour le test
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(ListVideos)
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
	req, err := http.NewRequest("GET", "/api/v2/videos/2", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Enregistre la réponse http pour le test
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetVideo)
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
	expected := Video{ID: "4", Title: "Le quatrième", Author: "Dylan Bru", PublishedDate: "2024-02-22", IsActive: true, DeletedAt: ""}
	requestBody := `{"title": "Le quatrième", "author": "Dylan Bru", "publishedDate": "2024-02-22"}`
	req, err := http.NewRequest("POST", "/api/v2/videos", strings.NewReader(requestBody))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateVideo)
	handler.ServeHTTP(rr, req)
	// Vérification du code de statut de la réponse
	if rr.Code != http.StatusCreated {
		t.Errorf("expected status Created; got %d", rr.Code)
	}
	// Vérification si la vidéo créée est dans la liste des vidéos
	found := false
	for _, v := range videos {
		if v.ID == expected.ID && v.Title == expected.Title && v.Author == expected.Author && v.PublishedDate == expected.PublishedDate && v.IsActive == expected.IsActive && v.DeletedAt == expected.DeletedAt {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected video %+v not found in response body: %v", expected, videos)
	}
}

func TestDeleteVideo(t *testing.T) {
	req, err := http.NewRequest("DELETE", "/api/v2/videos/2", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(DeleteVideo)
	handler.ServeHTTP(rr, req)
	// Vérification du code de statut de la réponse
	if rr.Code != http.StatusOK {
		t.Errorf("expected status Created; got %d", rr.Code)
	}
	// Vérification si les données de la vidéo ont bien été supprimées de la liste des vidéos
	found := false
	for _, v := range videos {
		if v.ID == "2" && v.Title == "" && v.Author == "" && v.PublishedDate == "" && v.IsActive == false && v.DeletedAt != "" {
			found = true
			break
		}
	}
	if !found {
		t.Error("vDeleted video not found")
	}
}
