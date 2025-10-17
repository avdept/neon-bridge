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

// FilteredJSON is a custom type for returning filtered config to frontend
type FilteredJSON map[string]interface{}

// List of sensitive field names to filter out
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

// filterSensitiveFields removes sensitive fields from a JSON object
func filterSensitiveFields(data map[string]interface{}) map[string]interface{} {
	filtered := make(map[string]interface{})

	for key, value := range data {
		// Check if the key is sensitive (case-insensitive)
		isSensitive := false
		lowerKey := strings.ToLower(key)
		for _, sensitiveField := range sensitiveFields {
			if strings.Contains(lowerKey, sensitiveField) {
				isSensitive = true
				break
			}
		}

		if !isSensitive {
			// If value is a nested object, recursively filter it
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
	Type        string         `json:"type" gorm:"not null"` // e.g., "adguard-home", "pihole", etc.
	Position    int            `json:"position" gorm:"default:0"`
	Config      JSON           `json:"config" gorm:"type:text"`     // Widget configuration
	LastState   JSON           `json:"last_state" gorm:"type:text"` // Last known state/data
	IsEnabled   bool           `json:"is_enabled" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// Relationships
	Dashboard Dashboard `json:"dashboard" gorm:"foreignKey:DashboardID"`
}

// WidgetResponse is the response structure for frontend (with filtered config)
type WidgetResponse struct {
	ID          uint         `json:"id"`
	DashboardID uint         `json:"dashboard_id"`
	Name        string       `json:"name"`
	Type        string       `json:"type"`
	Position    int          `json:"position"`
	Config      FilteredJSON `json:"config"`     // Filtered configuration
	LastState   JSON         `json:"last_state"` // Last known state/data
	IsEnabled   bool         `json:"is_enabled"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// ToResponse converts Widget to WidgetResponse with filtered config
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
