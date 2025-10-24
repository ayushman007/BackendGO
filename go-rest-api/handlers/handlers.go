package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)



type Customer struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Contacted bool   `json:"contacted"`
}

// store holds items in memory
var (
	mu        sync.RWMutex
	
	customers = []Customer{
		{ID: 1, Name: "Alice Johnson", Role: "Manager", Email: "alice@example.com", Phone: "555-0101", Contacted: true},
		{ID: 2, Name: "Bob Smith", Role: "Engineer", Email: "bob@example.com", Phone: "555-0202", Contacted: false},
		{ID: 3, Name: "Eve Davis", Role: "Support", Email: "eve@example.com", Phone: "555-0303", Contacted: false},
	}
	nextCustomer = 4
	
)



// CustomersHandler handles GET (list) and POST (create) on /customers
func CustomersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "application/json")
		mu.RLock()
		cpy := make([]Customer, len(customers))
		copy(cpy, customers)
		mu.RUnlock()
		json.NewEncoder(w).Encode(cpy)
	case http.MethodPost:
		var in struct {
			Name      string `json:"name"`
			Role      string `json:"role"`
			Email     string `json:"email"`
			Phone     string `json:"phone"`
			Contacted bool   `json:"contacted"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(in.Name) == "" {
			http.Error(w, "name required", http.StatusBadRequest)
			return
		}
		mu.Lock()
		c := Customer{
			ID:        nextCustomer,
			Name:      in.Name,
			Role:      in.Role,
			Email:     in.Email,
			Phone:     in.Phone,
			Contacted: in.Contacted,
		}
		nextCustomer++
		customers = append(customers, c)
		mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(c)
	default:
		log.Printf("Unsupported method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// CustomerByIDHandler handles GET, PUT, DELETE on /customers/{id}
func CustomerByIDHandler(w http.ResponseWriter, r *http.Request) {
	// support both mux and plain net/http style paths
	idStr := strings.TrimPrefix(r.URL.Path, "/customers/")
	if idStr == "" || strings.Contains(idStr, "/") {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		mu.RLock()
		var found Customer
		ok := false
		for _, c := range customers {
			if c.ID == id {
				found = c
				ok = true
				break
			}
		}
		mu.RUnlock()
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(found)
	case http.MethodPut:
		var in struct {
			Name      string `json:"name"`
			Role      string `json:"role"`
			Email     string `json:"email"`
			Phone     string `json:"phone"`
			Contacted bool   `json:"contacted"`
		}
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "invalid json", http.StatusBadRequest)
			return
		}
		mu.Lock()
		defer mu.Unlock()
		for i, c := range customers {
			if c.ID == id {
				customers[i].Name = in.Name
				customers[i].Role = in.Role
				customers[i].Email = in.Email
				customers[i].Phone = in.Phone
				customers[i].Contacted = in.Contacted
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(customers[i])
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	case http.MethodDelete:
		mu.Lock()
		defer mu.Unlock()
		for i, c := range customers {
			if c.ID == id {
				// remove by index
				customers = append(customers[:i], customers[i+1:]...)
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
