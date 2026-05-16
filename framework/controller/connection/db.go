package connection

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const menuURL = "#"

type SubOption struct {
	ID          string   `json:"id"`
	Label       string   `json:"label"`
	Description string   `json:"description"`
	Tag         string   `json:"tag"`
	Details     []string `json:"details"`
}

type MenuItem struct {
	ID          string      `json:"id"`
	Label       string      `json:"label"`
	Icon        string      `json:"icon"`
	Description string      `json:"description"`
	Tag         string      `json:"tag"`
	Details     []string    `json:"details"`
	SubOptions  []SubOption `json:"subOptions"`
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func Connection_DB() ([]MenuItem, error) {
	resp, err := httpClient.Get(menuURL)
	if err != nil {
		return nil, fmt.Errorf("error conectando al servidor: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("respuesta inesperada del servidor: %d", resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error leyendo respuesta: %w", err)
	}

	var menu []MenuItem
	if err := json.Unmarshal(data, &menu); err != nil {
		return nil, fmt.Errorf("error parseando JSON: %w", err)
	}

	if len(menu) == 0 {
		return nil, fmt.Errorf("el menú está vacío")
	}

	return menu, nil
}
