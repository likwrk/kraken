package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"kraken/api/app"
	"kraken/api/models"
)

// Mock repository for tests
type mockUserRepo struct {
	users []models.User
}

func (m *mockUserRepo) GetAll() ([]models.User, error) {
	return m.users, nil
}

func (m *mockUserRepo) Create(user models.User) error {
	m.users = append(m.users, user)
	return nil
}

func TestGetUsersHandler(t *testing.T) {
	mockRepo := &mockUserRepo{
		users: []models.User{
			{ID: 1, Name: "Alice", Age: 30},
			{ID: 2, Name: "Bob", Age: 25},
		},
	}
	a := app.NewApp(mockRepo)

	req := httptest.NewRequest("GET", "/users", nil)
	w := httptest.NewRecorder()
	a.GetUsers(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}

	var got []models.User
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatal(err)
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 users, got %d", len(got))
	}
}

func TestCreateUserHandler(t *testing.T) {
	mockRepo := &mockUserRepo{}
	a := app.NewApp(mockRepo)

	user := models.User{Name: "Charlie", Age: 28}
	body, _ := json.Marshal(user)

	req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	a.CreateUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}

	if len(mockRepo.users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(mockRepo.users))
	}

	if mockRepo.users[0].Name != "Charlie" {
		t.Fatalf("expected user name 'Charlie', got '%s'", mockRepo.users[0].Name)
	}
}
