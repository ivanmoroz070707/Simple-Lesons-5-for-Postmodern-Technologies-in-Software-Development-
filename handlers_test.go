package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"house-api/models"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ==========================================
// 1. ПОВНИЙ MOCK РЕПОЗИТОРІЙ
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

func (m *MockHouseRepo) UpdateFull(house *models.House) error {
	args := m.Called(house)
	return args.Error(0)
}

func (m *MockHouseRepo) UpdatePartial(id int, updateData map[string]interface{}) error {
	args := m.Called(id, updateData)
	return args.Error(0)
}

func (m *MockHouseRepo) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

// ==========================================
// 2. ТЕСТИ
// ==========================================

func TestCreateHouse_Success(t *testing.T) {
	mockRepo := new(MockHouseRepo)
	handler := &HouseHandler{repo: mockRepo}

	newHouse := models.House{Address: "Kyiv, Khreshchatyk", Rooms: 2, Price: 150000}
	
	mockRepo.On("Create", mock.AnythingOfType("*models.House")).Return(nil)

	body, _ := json.Marshal(newHouse)
	req := httptest.NewRequest(http.MethodPost, "/houses", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateHouse(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockRepo.AssertExpectations(t)
}

func TestCreateHouse_ValidationError(t *testing.T) {
	mockRepo := new(MockHouseRepo)
	handler := &HouseHandler{repo: mockRepo}

	invalidHouse := models.House{Address: "", Rooms: 0, Price: 100}
	
	body, _ := json.Marshal(invalidHouse)
	req := httptest.NewRequest(http.MethodPost, "/houses", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.CreateHouse(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestGetHouseByID_Success(t *testing.T) {
	mockRepo := new(MockHouseRepo)
	handler := &HouseHandler{repo: mockRepo}

	// Виправлено Floors на Rooms
	expectedHouse := &models.House{ID: 1, Address: "Odessa, Deribasivska", Rooms: 3, Price: 200000}
	
	mockRepo.On("GetByID", 1).Return(expectedHouse, nil)

	req := httptest.NewRequest(http.MethodGet, "/houses/1", nil)
	req.SetPathValue("id", "1") 
	w := httptest.NewRecorder()

	handler.GetHouseByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	
	var responseHouse models.House
	json.NewDecoder(w.Body).Decode(&responseHouse)
	assert.Equal(t, expectedHouse.Address, responseHouse.Address)
	
	mockRepo.AssertExpectations(t)
}

