package http

import (
	"GOTRAKER/internal/orders"
	"GOTRAKER/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

type Handler struct {
	serv *service.OrderService
}

func NewHandler(serv *service.OrderService) *Handler {
	return &Handler{serv: serv}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" && r.Method == http.MethodGet {
		w.Write([]byte("Добро пожаловать в GoTracker API"))
		return
	}
	path := r.URL.Path
	method := r.Method

	if path == "/orders" && method == http.MethodGet {
		h.handleGetAll(w, r)
		return
	}
	if path == "/orders" && method == http.MethodPost {
		h.handleCreate(w, r)
		return
	}
	if strings.HasPrefix(path, "/orders/") && method == http.MethodPut {
		h.handleUpdate(w, r)
		return
	}
	if strings.HasPrefix(path, "/orders/") && method == http.MethodGet {
		h.handleGetByID(w, r)
		return
	}
	http.NotFound(w, r)
}
func (h *Handler) handleGetAll(w http.ResponseWriter, r *http.Request) {
	orders := h.serv.GetAll()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
func (h *Handler) handleGetByID(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil || id <= 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	order, err := h.serv.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}
func (h *Handler) handleCreate(w http.ResponseWriter, r *http.Request) {
	var o *orders.Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	if o.ID <= 0 || o.Customer == "" || o.Address == "" {
		http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
		return
	}
	h.serv.AddOrder(r.Context(), o)
	w.WriteHeader(http.StatusCreated)
}
func (h *Handler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	id, err := strconv.Atoi(parts[2])
	if err != nil || id <= 0 {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var o orders.Order
	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	o.ID = id

	if err := h.serv.UpdateOrder(r.Context(), o); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
}
