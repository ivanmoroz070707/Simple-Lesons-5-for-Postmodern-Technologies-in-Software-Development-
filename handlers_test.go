package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"house-api/models"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ==========================================
// 1. MOCK РЕПОЗИТОРІЙ
// ==========================================
type MockHouseRepo struct {
	mock.Mock
}

func (m *MockHouseRepo) Create(house *models.House) error {
	args := m.Called(house)
	return args.Error(0)
}

func (m *MockHouseRepo) GetAll() ([]models.House, error) {
	args := m.Called()
	return args.Get(0).([]models.House), args.Error(1)
}

func (m *MockHouseRepo) GetByID(id int) (*models.House, error) {
	args := m.Called(id)
	if args.Get(0) != nil {
		return args.Get(0).(*models.House), args.Error(1)
	}
	return nil, args.Error(1)
}

// ВИПРАВЛЕНО: Тепер відповідає інтерфейсу (приймає тільки house)
func (m *MockHouseRepo) UpdateFull(house *models.House) error {
	args := m.Called(house)
	return args.Error(0)
}

func (m *MockHouseRepo) UpdatePartial(id int, updates map[string]interface{}) error {
	args := m.Called(id, updates)
	return args.Error(0)
}

func (m *MockHouseRepo) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// ==========================================
// 2. ТЕСТИ ДЛЯ ЕНДПОІНТІВ
// ==========================================

func TestCreateHouse(t *testing.T) {
	mockRepo := new(MockHouseRepo)
	handler := &HouseHandler{repo: mockRepo}

	newHouse := &models.House{Address: "Main St", Price: 1000}
	mockRepo.On("Create", mock.AnythingOfType("*models.House")).Return(nil)

	body, _ := json.Marshal(newHouse)
	req, _ := http.NewRequest("POST", "/houses", bytes.NewBuffer(body))
	rr := httptest.NewRecorder()

	handler.CreateHouse(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
}

func TestPatchHouse(t *testing.T) {
	mockRepo := new(MockHouseRepo)
	handler := &HouseHandler{repo: mockRepo}

	houseID := 1
	updateData := map[string]interface{}{"price": 150000.0}
	updatedHouse := &models.House{ID: houseID, Address: "Test St", Price: 150000.0}

	mockRepo.On("UpdatePartial", houseID, updateData).Return(nil)
	mockRepo.On("GetByID", houseID).Return(updatedHouse, nil)

	body, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PATCH", "/houses/1", bytes.NewBuffer(body))
	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.PatchHouse(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockRepo.AssertExpectations(t)
}

func TestDeleteHouse(t *testing.T) {
	mockRepo := new(MockHouseRepo)
	handler := &HouseHandler{repo: mockRepo}

	mockRepo.On("Delete", 1).Return(nil)

	req, _ := http.NewRequest("DELETE", "/houses/1", nil)
	
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	rr := httptest.NewRecorder()
	handler.DeleteHouse(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	mockRepo.AssertExpectations(t)
}
