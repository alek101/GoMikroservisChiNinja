package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/alek101/GoMikroservisChiNinja/model"
	orderrepo "github.com/alek101/GoMikroservisChiNinja/repository/order"
)

type Order struct {
	DB *sql.DB
}

func (o *Order) Create(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Create an order")

	// parse JSON request body into order struct
	var sample model.Order
	if err := json.NewDecoder(r.Body).Decode(&sample); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	repo := orderrepo.NewMariaRepo(o.DB)
	if err := repo.Insert(r.Context(), sample); err != nil {
		http.Error(w, "failed to save order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sample)
}

func (o *Order) List(w http.ResponseWriter, r *http.Request) {
	// retrieve all orders from repository
	repo := orderrepo.NewMariaRepo(o.DB)
	orders, err := repo.FindAll(r.Context())
	if err != nil {
		http.Error(w, "failed to fetch orders", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

func (o *Order) GetByID(w http.ResponseWriter, r *http.Request) {
	// parse ID from path
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	repo := orderrepo.NewMariaRepo(o.DB)
	ord, err := repo.FindByID(r.Context(), id)
	if err != nil {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ord)
}

func (o *Order) UpdateByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Update an order by ID")

	// parse ID from path
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var ord model.Order
	if err := json.NewDecoder(r.Body).Decode(&ord); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}
	ord.OrderID = id

	repo := orderrepo.NewMariaRepo(o.DB)
	if err := repo.Update(r.Context(), ord); err != nil {
		http.Error(w, "failed to update order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(ord)
}

func (o *Order) DeleteByID(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Delete an order by ID")

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	repo := orderrepo.NewMariaRepo(o.DB)
	if err := repo.Delete(r.Context(), id); err != nil {
		http.Error(w, "failed to delete order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
