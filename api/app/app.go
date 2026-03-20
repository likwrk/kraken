package app

import (
	"kraken/api/repository"

	"github.com/gorilla/mux"
)

type App struct {
	UserRepo    repository.UserRepository
	SensorsRepo repository.SensorsRepository
}

func NewApp(userRepo repository.UserRepository, sensorsRepo repository.SensorsRepository) *App {
	return &App{UserRepo: userRepo, SensorsRepo: sensorsRepo}
}

func (a *App) Router() *mux.Router {
	r := mux.NewRouter()
	// r.HandleFunc("/users", a.GetUsers).Methods("GET")
	// r.HandleFunc("/users", a.CreateUser).Methods("POST")
	r.HandleFunc("/sensors", a.SensorsAllReadings).Methods("GET")
	r.HandleFunc("/sensors_set", a.SensorsAllSets).Methods("GET")
	r.HandleFunc("/sensors2", a.SensorsHandler).Methods("GET")
	r.HandleFunc("/sensors", a.SensorDataAddHandler).Methods("POST")
	return r
}
