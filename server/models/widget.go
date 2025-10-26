package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

type JSON map[string]interface{}

func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSON)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into JSON", value)
	}

	result := make(JSON)
	err := json.Unmarshal(bytes, &result)
	*j = result
	return err
}

func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

type FilteredJSON map[string]interface{}

var sensitiveFields = []string{
	"password",
	"token",
	"api_key",
	"apikey",
	"api_token",
	"secret",
	"auth",
	"credential",
	"key",
	"private_key",
	"access_token",
	"refresh_token",
}

func filterSensitiveFields(data map[string]interface{}) map[string]interface{} {
	filtered := make(map[string]interface{})

	for key, value := range data {
		isSensitive := false
		lowerKey := strings.ToLower(key)
		for _, sensitiveField := range sensitiveFields {
			if strings.Contains(lowerKey, sensitiveField) {
				isSensitive = true
				break
			}
		}

		if !isSensitive {
			if nestedMap, ok := value.(map[string]interface{}); ok {
				filtered[key] = filterSensitiveFields(nestedMap)
			} else {
				filtered[key] = value
			}
		}
	}

	return filtered
}

type Widget struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	DashboardID uint           `json:"dashboard_id" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"not null"`
	Type        string         `json:"type" gorm:"not null"`
	Position    int            `json:"position" gorm:"default:0"`
	Config      JSON           `json:"config" gorm:"type:text"`
	LastState   JSON           `json:"last_state" gorm:"type:text"`
	IsEnabled   bool           `json:"is_enabled" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	Dashboard Dashboard `json:"dashboard" gorm:"foreignKey:DashboardID"`
}

type WidgetResponse struct {
	ID          uint         `json:"id"`
	DashboardID uint         `json:"dashboard_id"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Position    int          `json:"position"`
	Config      FilteredJSON `json:"config"`
	LastState   JSON         `json:"last_state"`
	IsEnabled   bool         `json:"is_enabled"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

func (w *Widget) ToResponse() WidgetResponse {
	return WidgetResponse{
		ID:          w.ID,
		DashboardID: w.DashboardID,
		Name:        w.Name,
		Type:        w.Type,
		Position:    w.Position,
		Config:      FilteredJSON(filterSensitiveFields(map[string]interface{}(w.Config))),
		LastState:   w.LastState,
		IsEnabled:   w.IsEnabled,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
	}
}
