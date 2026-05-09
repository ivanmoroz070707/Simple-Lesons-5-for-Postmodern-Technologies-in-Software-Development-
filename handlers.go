package main

import (
	"encoding/json"
	"net/http"
	"strconv" 
	"house-api/models" 
)

type HouseHandler struct {
	repo models.HouseRepository
}

// 1. POST /houses — Створення
func (h *HouseHandler) CreateHouse(w http.ResponseWriter, r *http.Request) {
	var house models.House
	if err := json.NewDecoder(r.Body).Decode(&house); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := house.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.Create(&house); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(house)
}

// 2. GET /houses — Отримання всіх
func (h *HouseHandler) GetAllHouses(w http.ResponseWriter, r *http.Request) {
	houses, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(houses)
}

// 3. GET /houses/{id} — Отримання одного
func (h *HouseHandler) GetHouseByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	house, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "House not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(house)
}

// 4. PUT /houses/{id} — Повне оновлення
func (h *HouseHandler) UpdateHouse(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var house models.House
	if err := json.NewDecoder(r.Body).Decode(&house); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	house.ID = id // Гарантуємо, що ID з URL співпадає з моделлю

	if err := house.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdateFull(&house); err != nil {
		http.Error(w, "Update failed", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(house)
}

// 5. PATCH /houses/{id} — Часткове оновлення (UpdatePartial)
func (h *HouseHandler) UpdateHousePartial(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	// Декодуємо JSON у map, щоб передати в UpdatePartial напарника
	var updateData map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdatePartial(id, updateData); err != nil {
		http.Error(w, "Partial update failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// 6. DELETE /houses/{id} — Видалення
func (h *HouseHandler) DeleteHouse(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, "Delete failed", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}