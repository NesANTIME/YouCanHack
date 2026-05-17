package database

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const menuURL = "https://gist.githubusercontent.com/NesANTIME/48b1f120d87039819d86d6a7e09ad8af/raw/814ea28955ebcd738a31cbaf65deb4dacfdbe801/prueba"

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
