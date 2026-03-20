package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"kraken/api/models"
	"log"
	"net/http"
)

type sensorsRepo struct {
	db            *sql.DB
	allowedTables map[string]bool
}

func NewSensorsRepository(db *sql.DB) SensorsRepository {
	return &sensorsRepo{
		db: db,
		allowedTables: map[string]bool{
			"sensor_sets":     true,
			"sensor_readings": true,
		},
	}
}

func (r *sensorsRepo) getLatestValues(deviceId string, tableName string) ([]models.SensorTemp, error) {
	if !r.allowedTables[tableName] {
		return nil, fmt.Errorf("invalid table: %s", tableName)
	}
	query := fmt.Sprintf(`
		SELECT device_id, sensor_id, sensor_value, timestamp
		FROM (
		    SELECT *,
		           ROW_NUMBER() OVER (
		               PARTITION BY sensor_id
		               ORDER BY timestamp DESC
		           ) AS rn
		    FROM %s
		    WHERE device_id = ?
		) t
		WHERE rn = 1;
`, tableName)
	rows, err := r.db.Query(query, deviceId)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	var temps []models.SensorTemp

	for rows.Next() {
		var t models.SensorTemp

		err := rows.Scan(
			&t.DeviceID,
			&t.SensorID,
			&t.Value,
			&t.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		temps = append(temps, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return temps, nil
}

func (r *sensorsRepo) innerGetSensorsTableRows(deviceID string, tableName string) (*sql.Rows, error) {
	if !r.allowedTables[tableName] {
		return nil, fmt.Errorf("invalid table: %s", tableName)
	}
	query := fmt.Sprintf(`
    SELECT device_id, sensor_id, sensor_value, timestamp
    FROM %s
    WHERE device_id = ?
    ORDER BY timestamp ASC
`, tableName)

	rows, err := r.db.Query(query, deviceID)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (r *sensorsRepo) innerGetSensorsValuesFromTable(deviceID string, tableName string) ([]models.SensorTemp, error) {
	rows, err := r.innerGetSensorsTableRows(deviceID, tableName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var temps []models.SensorTemp

	for rows.Next() {
		var t models.SensorTemp

		err := rows.Scan(
			&t.DeviceID,
			&t.SensorID,
			&t.Value,
			&t.Timestamp,
		)
		if err != nil {
			return nil, err
		}
		temps = append(temps, t)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return temps, nil
}

func (r *sensorsRepo) innerWriteSensorsValuesToTable(tableName string, deviceID string, sets []models.SensorsSet) error {
	if !r.allowedTables[tableName] {
		return fmt.Errorf("invalid table: %s", tableName)
	}
	count := len(sets)
	if count < 1 {
		return nil
	}
	latestValues, err := r.getLatestValues(deviceID, tableName)
	latestValuesPerSensor := make(map[string]*models.SensorTemp, len(latestValues))
	for i := range latestValues {
		log.Println(latestValues[i].SensorID, latestValues[i].Value)
		latestValuesPerSensor[latestValues[i].SensorID] = &latestValues[i]
	}
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	// Always ensure rollback if something fails
	defer tx.Rollback()

	for _, set := range sets {
		for _, value := range set.Values {
			exisitngValue, ok := latestValuesPerSensor[value.SensorID]
			if ok {
				if exisitngValue.Value == value.Value {
					continue
				}
				exisitngValue.Value = value.Value
			}

			query := fmt.Sprintf(`INSERT INTO %s (device_id, sensor_id, sensor_value, timestamp) VALUES (?, ?, ?, ?)`, tableName)
			_, err = tx.Exec(
				query,
				deviceID, value.SensorID, value.Value, value.Timestamp.Unix(),
			)

			if err != nil {
				return err
			}
		}
	}

	// Commit only if everything succeeded
	return tx.Commit()
}

func (r *sensorsRepo) innerGetAllFromTable(tableName string, deviceID string, w http.ResponseWriter) error {
	if !r.allowedTables[tableName] {
		return fmt.Errorf("invalid table: %s", tableName)
	}
	flusher := w.(http.Flusher)
	rows, err := r.innerGetSensorsTableRows(deviceID, tableName)
	if err != nil {
		return err
	}
	defer rows.Close()
	if _, err := w.Write([]byte("[")); err != nil {
		return err
	}
	if flusher != nil {
		flusher.Flush() // Send '[' immediately
	}
	first := true
	for rows.Next() {
		var sensor models.SensorTemp
		if err := rows.Scan(&sensor.DeviceID, &sensor.SensorID, &sensor.Value, &sensor.Timestamp); err != nil {
			return err
		}

		if !first {
			if _, err := w.Write([]byte(",")); err != nil {
				return err
			}
		}
		first = false

		if err := json.NewEncoder(w).Encode(sensor); err != nil {
			return err
		}

		// Optional: Flush periodically for large datasets
		if flusher != nil && !first {
			flusher.Flush()
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if _, err := w.Write([]byte("]")); err != nil {
		return err
	}
	return nil
}

func (r *sensorsRepo) GetAllReadings(deviceID string, w http.ResponseWriter) error {
	return r.innerGetAllFromTable("sensor_readings", deviceID, w)
}

func (r *sensorsRepo) GetAllSets(deviceID string, w http.ResponseWriter) error {
	return r.innerGetAllFromTable("sensor_sets", deviceID, w)
}

func (r *sensorsRepo) WriteSensorsSetValues(deviceID string, sets []models.SensorsSet) error {
	return r.innerWriteSensorsValuesToTable("sensor_sets", deviceID, sets)
}

func (r *sensorsRepo) GetSensorsValues(deviceId string) ([]models.SensorTemp, error) {
	return r.innerGetSensorsValuesFromTable(deviceId, "sensor_readings")
}

func (r *sensorsRepo) GetSensorsSetValues(deviceId string) ([]models.SensorTemp, error) {
	return r.innerGetSensorsValuesFromTable(deviceId, "sensor_sets")
}

func (r *sensorsRepo) WriteSensorsValues(deviceID string, sets []models.SensorsSet) error {
	return r.innerWriteSensorsValuesToTable("sensor_readings", deviceID, sets)
}
