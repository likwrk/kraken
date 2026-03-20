package app

import (
	"encoding/json"
	"log"
	"net/http"

	"kraken/api/models"
)

func (a *App) SensorsHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("device")
	if deviceID == "" {
		http.Error(w, "no device id", http.StatusBadRequest)
		return
	}
	temps, err := a.SensorsRepo.GetSensorsValues(deviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(temps)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (a *App) SensorsAllReadings(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("device")
	if deviceID == "" {
		http.Error(w, "no device id", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Content-Type-Options", "nosniff") // Security
	w.WriteHeader(http.StatusOK)
	err := a.SensorsRepo.GetAllReadings(deviceID, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) SensorsAllSets(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("device")
	if deviceID == "" {
		http.Error(w, "no device id", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-Content-Type-Options", "nosniff") // Security
	w.WriteHeader(http.StatusOK)
	err := a.SensorsRepo.GetAllSets(deviceID, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) SensorDataAddHandler(w http.ResponseWriter, r *http.Request) {
	deviceID := r.URL.Query().Get("id")
	if deviceID == "" {
		http.Error(w, "no device id", http.StatusBadRequest)
		return
	}
	var request []models.SensorsSetRequest
	log.Println("new request from ")
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println("requset timestamp", request[0].Timestamp)
	result := make([]models.SensorsSet, len(request))
	for i, v := range request {
		result[i] = v.ToSensorsSet()
	}

	log.Println("Total data len:", len(result))

	if len(result) < 1 {
		w.WriteHeader(http.StatusCreated)
		return
	}

	if err := a.SensorsRepo.WriteSensorsValues(deviceID, result); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := a.SensorsRepo.WriteSensorsSetValues(deviceID, result); err != nil {
		log.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}
