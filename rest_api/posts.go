package rest_api

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/nireo/booru/lib"
	"github.com/nireo/booru/models"
)

func ServePostsInBoardJSON(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// fetch board name
	keys, ok := r.URL.Query()["board"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "You need to provide board name", http.StatusBadRequest)
		return
	}

	db := lib.GetDatabase()

	var board models.Board
	if err := db.Where(&models.Board{Link: keys[0]}).First(&board).Error; err != nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	var posts []models.Post
	db.Model(&board).Related(&posts)

	for index := range posts {
		comments, _ := posts[index].GetComments()
		posts[index].Comments = comments
	}

	json.NewEncoder(w).Encode(posts)
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	keys, ok := r.URL.Query()["board"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "You need to provide board name", http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Error parsing multipart data", http.StatusInternalServerError)
		return
	}

	db := lib.GetDatabase()

	var board models.Board
	if err := db.Where(&models.Board{Link: keys[0]}).First(&board).Error; err != nil {
		http.Error(w, "Board not found", http.StatusNotFound)
		return
	}

	var fileExt string
	var file multipart.File

	if r.Form.Has("file") {
		file, handler, err := r.FormFile("file")
		defer file.Close()
		if err != nil {
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		fileExt = filepath.Ext(handler.Filename)
	}

	newPost := &models.Post{
		UUID:          lib.UUID(),
		CreatedBy:     "Anonymous",
		Content:       r.FormValue("content"),
		BoardID:       board.ID,
		FileExtension: fileExt,
	}

	if file != nil && fileExt != "" {
		filepath := fmt.Sprintf("./images/%s%s", newPost.UUID, newPost.FileExtension)
		dst, err := os.Create(filepath)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	db.NewRecord(newPost)
	db.Create(newPost)
	json.NewEncoder(w).Encode(newPost)
}

func CreateComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	keys, ok := r.URL.Query()["post"]
	if !ok || len(keys[0]) < 1 {
		http.Error(w, "You need to provide post uuid", http.StatusBadRequest)
		return
	}

	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	db := lib.GetDatabase()
	var post models.Post
	if err := db.Where(&models.Post{UUID: keys[0]}).First(&post).Error; err != nil {
		http.Error(w, "Post not found", http.StatusNotFound)
		return
	}

	if r.FormValue("content") == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	withFile := true
	file, handler, err := r.FormFile("file")
	if err != nil {
		withFile = false
	}

	newComment := &models.Comment{
		UUID:      lib.UUID(),
		CreatedBy: "Anonymous",
		Content:   r.FormValue("content"),
		PostID:    post.ID,
	}

	if withFile {
		newComment.FileExtension = filepath.Ext(handler.Filename)
		// save the file
		defer file.Close()
		filePath := fmt.Sprintf("./images/%s%s", newComment.UUID, newComment.FileExtension)
		dst, err := os.Create(filePath)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer dst.Close()

		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// finally create the database entry, and redirect the user to the board page
	db.NewRecord(newComment)
	db.Create(newComment)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newComment)
}
