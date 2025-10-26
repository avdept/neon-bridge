package models

import (
	"time"

	"gorm.io/gorm"
)

type Dashboard struct {
	ID            uint           `json:"id" gorm:"primaryKey"`
	Name          string         `json:"name" gorm:"not null"`
	Description   string         `json:"description"`
	GlancesConfig string         `json:"glances_config" gorm:"type:json"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"-" gorm:"index"`

	Widgets []Widget `json:"widgets" gorm:"foreignKey:DashboardID"`
}

type DashboardResponse struct {
	ID            uint             `json:"id"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	GlancesConfig string           `json:"glances_config"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
	Widgets       []WidgetResponse `json:"widgets"`
}

func (d *Dashboard) ToResponse() DashboardResponse {
	var widgetResponses []WidgetResponse
	for _, widget := range d.Widgets {
		widgetResponses = append(widgetResponses, widget.ToResponse())
	}

	return DashboardResponse{
		ID:            d.ID,
		Name:          d.Name,
		Description:   d.Description,
		GlancesConfig: d.GlancesConfig,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
		Widgets:       widgetResponses,
	}
}
