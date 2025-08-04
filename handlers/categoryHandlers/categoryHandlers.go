package categoryhandlers

import (
	"CRUD-go/config"
	"CRUD-go/entities"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// func CategoryIndex(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	w.Write([]byte(`{"message": "Category index works!"}`))
// }

func CategoryIndex(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintln(w, "Category List")
	rows, err := config.DB.Query("SELECT * FROM categories")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	var categories []entities.Category

	for rows.Next() {
		var category entities.Category
		if err := rows.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt); err != nil {
			panic(err)
		}
		categories = append(categories, category)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(categories)
}
func CategoryFindByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	// idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error":"Missing id parameter"}`, http.StatusBadRequest)
		return
	}

	parsedID, err := strconv.Atoi(idStr)
	if err != nil || parsedID < 0 {
		http.Error(w, `{"error":"Invalid id"}`, http.StatusBadRequest)
		return
	}
	id := uint(parsedID)

	var category entities.Category
	row := config.DB.QueryRow("SELECT id, name, created_at, updated_at FROM categories WHERE id = ?", id)
	err = row.Scan(&category.Id, &category.Name, &category.CreatedAt, &category.UpdatedAt)
	if err != nil {
		http.Error(w, `{"error":"Category not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func CategoryStore(w http.ResponseWriter, r *http.Request) {
	var category entities.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	category.CreatedAt = time.Now()
	category.UpdatedAt = time.Now()

	result, err := config.DB.Exec("INSERT INTO categories (name, created_at, updated_at) VALUES (?,?,?)", category.Name, category.CreatedAt, category.UpdatedAt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	id, _ := result.LastInsertId()
	category.Id = uint(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)

}

func CategoryUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	parsedID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}
	id := uint(parsedID)

	var category entities.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	category.UpdatedAt = time.Now()

	_, err = config.DB.Exec("UPDATE categories SET name=?, updated_at=? WHERE id=?", category.Name, category.UpdatedAt, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	category.Id = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(category)
}

func CategoryDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]

	parsedID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}
	id := uint(parsedID)

	result, err := config.DB.Exec("DELETE FROM categories WHERE id = ?", id)
	if err != nil {
		http.Error(w, `{"error":"Failed to delete category"}`, http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error":"Category not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Category deleted successfully"}`))
}
