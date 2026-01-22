package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Global variable to store categories (for demo purpose)
var categories = []Category{
	{ID: 1, Name: "Buku", Description: "Pengelompokan berbagai jenis buku."},
	{ID: 2, Name: "Elektronik", Description: "Pengelompokan perangkat elektronik."},
	{ID: 3, Name: "Pakaian", Description: "Pengelompokan berbagai jenis pakaian."},
}

func main() {
	http.HandleFunc("/api/categories", categoriesHandler)
	http.HandleFunc("/api/categories/", categoryHandler)

	fmt.Println("Server running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

/*
|--------------------------------------------------------------------------
| Handlers
|--------------------------------------------------------------------------
*/

func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {

	case http.MethodGet:
		respondJSON(w, http.StatusOK, categories)

	case http.MethodPost:
		defer r.Body.Close()

		var input Category
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		input.ID = generateID()
		categories = append(categories, input)

		respondJSON(w, http.StatusCreated, input)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func categoryHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	index := findCategoryIndex(id)
	if index == -1 {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	switch r.Method {

	case http.MethodGet:
		respondJSON(w, http.StatusOK, categories[index])

	case http.MethodPut:
		defer r.Body.Close()

		var input Category
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		input.ID = id
		categories[index] = input

		respondJSON(w, http.StatusOK, input)

	case http.MethodDelete:
		categories = append(categories[:index], categories[index+1:]...)
		respondJSON(w, http.StatusOK, map[string]string{
			"message": "Category deleted",
		})

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

/*
|--------------------------------------------------------------------------
| Helpers
|--------------------------------------------------------------------------
*/

func respondJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func parseID(path string) (int, error) {
	idStr := strings.TrimPrefix(path, "/api/categories/")
	return strconv.Atoi(idStr)
}

func findCategoryIndex(id int) int {
	for i, c := range categories {
		if c.ID == id {
			return i
		}
	}
	return -1
}

func generateID() int {
	if len(categories) == 0 {
		return 1
	}
	return categories[len(categories)-1].ID + 1
}
