package repository

import (
	"kraken/api/models"
	"net/http"
)

type UserRepository interface {
	GetAll() ([]models.User, error)
	Create(user models.User) error
}

type SensorsRepository interface {
	GetAllReadings(deviceID string, w http.ResponseWriter) error
	GetAllSets(deviceID string, w http.ResponseWriter) error
	GetSensorsValues(setID string) ([]models.SensorTemp, error)
	GetSensorsSetValues(setID string) ([]models.SensorTemp, error)
	WriteSensorsValues(deviceID string, sets []models.SensorsSet) error
	WriteSensorsSetValues(deviceID string, sets []models.SensorsSet) error
}
