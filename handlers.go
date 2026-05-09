package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv" 
	"house-api/models" 
	"github.com/go-chi/chi/v5"
)

type HouseHandler struct {
	repo models.HouseRepository
}

func (h *HouseHandler) CreateHouse(w http.ResponseWriter, r *http.Request) {
    var house models.House
    
    if err := json.NewDecoder(r.Body).Decode(&house); err != nil {
        log.Printf("Помилка декодування JSON: %v", err)
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    if err := house.Validate(); err != nil {
        log.Printf("Помилка валідації: %v", err)
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    err := h.repo.Create(&house)
    if err != nil {
        log.Printf("КРИТИЧНА ПОМИЛКА БАЗИ ДАНИХ: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(house)
}

func (h *HouseHandler) GetAllHouses(w http.ResponseWriter, r *http.Request) {
	houses, err := h.repo.GetAll()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(houses)
}

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

func (h *HouseHandler) GetHouses(w http.ResponseWriter, r *http.Request) {
    houses, err := h.repo.GetAll()
    if err != nil {
        log.Printf("Помилка при отриманні списку будинків: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(houses)
}

func (h *HouseHandler) UpdateHouse(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некоректний ID", http.StatusBadRequest)
		return
	}

	var house models.House
	if err := json.NewDecoder(r.Body).Decode(&house); err != nil {
		http.Error(w, "Некоректний JSON", http.StatusBadRequest)
		return
	}

	house.ID = id 

	if err := h.repo.UpdateFull(&house); err != nil {
		log.Printf("Помилка оновлення: %v", err)
		http.Error(w, "Помилка сервера", http.StatusInternalServerError)
		return
	}

        w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(house)
}

func (h *HouseHandler) UpdateHousePartial(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

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

func (h *HouseHandler) DeleteHouse(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некоректний ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		log.Printf("Помилка видалення: %v", err)
		http.Error(w, "Помилка сервера", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusNoContent)
}

func (h *HouseHandler) PatchHouse(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Некоректний ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Некоректний JSON", http.StatusBadRequest)
		return
	}

	if price, ok := updates["price"].(float64); ok && price <= 0 {
		http.Error(w, "Ціна має бути більшою за нуль", http.StatusBadRequest)
		return
	}

	if err := h.repo.UpdatePartial(id, updates); err != nil {
		log.Printf("Помилка PATCH: %v", err)
		http.Error(w, "Помилка сервера", http.StatusInternalServerError)
		return
	}

	house, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Будинок оновлено, але не знайдено", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(house)
}
