package producthandlers

import (
	"CRUD-go/config"
	"CRUD-go/entities"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func ProductIndex(w http.ResponseWriter, r *http.Request) {
	rows, err := config.DB.Query(`
		SELECT 
			p.id, p.name, p.stock, p.description, p.created_at, p.updated_at,
			c.name as category_name
		FROM products p
		JOIN categories c ON p.category_id = c.id
	`)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var products []entities.Product

	for rows.Next() {
		var product entities.Product
		err := rows.Scan(
			&product.Id,
			&product.Name,
			&product.Stock,
			&product.Description,
			&product.CreatedAt,
			&product.UpdatedAt,
			&product.Category.Name,
		)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}
func ProductFindByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	parsedID, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}
	id := uint(parsedID)

	var product entities.Product
	err = config.DB.QueryRow(`
		SELECT 
			p.id, p.name, p.stock, p.description, p.created_at, p.updated_at,
			c.id, c.name, c.created_at, c.updated_at
		FROM products p
		JOIN categories c ON p.category_id = c.id
		WHERE p.id = ?
	`, id).Scan(
		&product.Id,
		&product.Name,
		&product.Stock,
		&product.Description,
		&product.CreatedAt,
		&product.UpdatedAt,
		&product.Category.Id,
		&product.Category.Name,
		&product.Category.CreatedAt,
		&product.Category.UpdatedAt,
	)
	if err != nil {
		http.Error(w, `{"error":"Product not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}
func ProductStore(w http.ResponseWriter, r *http.Request) {
	var Product entities.Product
	if err:= json.NewDecoder(r.Body).Decode(&Product); err!=nil{
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	Product.CreatedAt = time.Now()
	Product.UpdatedAt = time.Now()

	result, err := config.DB.Exec(`
		INSERT INTO products (name, stock, description, category_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		product.Name, product.Stock, product.Description, product.Category.Id, product.CreatedAt, product.UpdatedAt
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	product.Id = uint(id)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func ProductUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars{r}
	idStr := vars["id"]
	parseID, err := strconv.Atoi(idStr)
	if err!=nil{
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
	}
	id:=uint(parseID)

	var product entities.Product
	if err:=json.NewDecoder(r.Body).Decode(&product); err!=nil{
		http.Error(w, `{"error":"Invalid JSON"}`, http.StatusBadRequest)
	}
	product.UpdatedAt = time.Now()

	_, err = config.DB.Exec(`
		UPDATE products 
		SET name=?, stock=?, description=?, category_id=?, updated_at=? 
		WHERE id=?`,
		product.Name, product.Stock, product.Description, product.Category.Id, product.UpdatedAt, id,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	product.Id = id
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func ProductDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr := vars["id"]
	parsedID, err := strconv.Atoi(idStr)
	if err:= nil{
		http.Error(w, `{"error":"Invalid ID"}`, http.StatusBadRequest)
		return
	}

	id:=uint(parsedID)

	result, err := config.DB.Exec("DELETE FROM products WHERE id = ?", id)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected,_ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, `{"error":"Product not found"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message":"Product deleted successfully"}`))
}
