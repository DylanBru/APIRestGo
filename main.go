package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//--------------AMELIORATIONS--------------//

// Gérer la génération automatique de l'id lors de la création d'une vidéo : package uuid de google

//--------------MEMO--------------//
// & : adresse mémoire
// * : valeur pointée par l'adresse mémoire

//--------------REQUEST--------------//

// GET : curl -X GET http://localhost:8000/api/v2/videos
// POST : curl -X POST -H 'content-type: application/json' --data '{"title": "Le quatrième", "author": "Dylan Bru", "publishedDate": "2024-02-22"}' http://localhost:8000/api/v2/videos
// DELETE : curl -X DELETE http://localhost:8000/api/v2/videos/3

//--------------BDD--------------//

type Video struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	Author        string `json:"author"`
	PublishedDate string `json:"publishedDate"`
	IsActive      bool   `json:"isActive"`
	DeletedAt     string `json:"deletedAt"`
}

var videos = []Video{
	{ID: "1", Title: "Voyage Culinaire", Author: "Sophie Dubois", PublishedDate: "2023-06-15", IsActive: true, DeletedAt: ""},
	{ID: "2", Title: "Le Chemin de l'Aventure", Author: "Thomas Leduc", PublishedDate: "2021-03-05", IsActive: true, DeletedAt: ""},
	{ID: "3", Title: "Exploration : Mars", Author: "Mariel Lefèvre", PublishedDate: "2020-06-01", IsActive: true, DeletedAt: ""},
}

//--------------CONTROLLER--------------//

func ListVideos(w http.ResponseWriter, r *http.Request) {
	// Encode les vidéos en JSON et les écrit dans la réponse
	json.NewEncoder(w).Encode(videos)
}

func GetVideo(w http.ResponseWriter, r *http.Request) {
	// Extraction de l'ID de l'URL
	urlParts := strings.Split(r.URL.Path, "/")
	videoID := urlParts[len(urlParts)-1]
	if videoID == "" {
		http.Error(w, "Missing video ID", http.StatusBadRequest)
		return
	}
	// Recherche de la vidéo avec l'ID donné
	index := -1
	for i, v := range videos {
		if v.ID == videoID {
			index = i
			break
		}
	}
	// Si la vidéo n'est pas trouvée, 404
	if index == -1 {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}
	// Ecrit la vidéo recherchée dans la réponse
	json.NewEncoder(w).Encode(videos[index])
}

func CreateVideo(w http.ResponseWriter, r *http.Request) {
	var newVideo Video
	// Erreur si le corps de la requête n'est pas conforme à la structure de video
	err := json.NewDecoder(r.Body).Decode(&newVideo)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	// Auto-incrémentation de l'id unique
	newVideo.ID = strconv.Itoa(len(videos) + 1)
	newVideo.IsActive = true
	videos = append(videos, newVideo)
	w.WriteHeader(http.StatusCreated)
}

func DeleteVideo(w http.ResponseWriter, r *http.Request) {
	// Extraction de l'ID de l'URL
	urlParts := strings.Split(r.URL.Path, "/")
	videoID := urlParts[len(urlParts)-1]
	if videoID == "" {
		http.Error(w, "Missing video ID", http.StatusBadRequest)
		return
	}
	// Recherche de la vidéo avec l'ID donné
	index := -1
	for i, v := range videos {
		if v.ID == videoID {
			index = i
			break
		}
	}
	// Si la vidéo n'est pas trouvée, 404
	if index == -1 {
		http.Error(w, "Video not found", http.StatusNotFound)
		return
	}
	// Supprimer les éléments de la vidéo ciblée
	now := time.Now().Format("2006-01-02")
	videos[index].Title = ""
	videos[index].Author = ""
	videos[index].PublishedDate = ""
	videos[index].IsActive = false
	videos[index].DeletedAt = now

	// videos = append(videos[:index], videos[index+1:]...)
	w.WriteHeader(http.StatusOK)
}

func main() {
	//--------------ROUTER+DISPATCHER--------------//
	mux := http.NewServeMux()
	// Enregistrement de chaque handler pour chaque méthode
	mux.HandleFunc("GET /api/v2/videos", ListVideos)
	mux.HandleFunc("GET /api/v2/videos/{id}", GetVideo)
	mux.HandleFunc("POST /api/v2/videos", CreateVideo)
	mux.HandleFunc("DELETE /api/v2/videos/{id}", DeleteVideo)
	// Lancement du serveur HTTP sur le port 8000
	fmt.Println("Serveur démarré sur le port 8000")
	http.ListenAndServe("localhost:8000", mux)
}
