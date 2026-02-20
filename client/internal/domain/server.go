package domain

import "time"

// represent selectable server
type Server struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Port        int       `json:"port"`
	Protocol    string    `json:"protocol"`
	Username    string    `json:"username"`
	Desc        string    `json:"description"`
	Environment string    `json:"environment"`
	Region      string    `json:"region"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Implements item.Item (bubbles)
func (s Server) FilterValue() string {
	return s.Name
}

func (s Server) Title() string {
	return s.Name
}

func (s Server) Description() string {
	return s.Address
}
