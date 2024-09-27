package manifest

import (
	"encoding/json"
	"io"
	"os"
)

// ModLoader Структура для ModLoader
type ModLoader struct {
	ID      string `json:"id"`
	Primary bool   `json:"primary"`
}

// File Структура для файлів
type File struct {
	ProjectID int  `json:"projectID"`
	FileID    int  `json:"fileID"`
	Required  bool `json:"required"`
}

// Minecraft Структура для Minecraft
type Minecraft struct {
	Version    string      `json:"version"`
	ModLoaders []ModLoader `json:"modLoaders"`
}

// Manifest Структура для Manifest
type Manifest struct {
	Minecraft       Minecraft `json:"minecraft"`
	ManifestType    string    `json:"manifestType"`
	ManifestVersion int       `json:"manifestVersion"`
	Name            string    `json:"name"`
	Version         string    `json:"version"`
	Author          string    `json:"author"`
	Files           []File    `json:"files"`
	Overrides       string    `json:"overrides"`
}

// LoadManifest Функція для читання JSON з файлу
func LoadManifest(filename string) (*Manifest, error) {
	// Відкриваємо файл
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	// Читаємо вміст файлу
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// Створюємо екземпляр Manifest
	var manifest Manifest

	// Парсимо JSON
	err = json.Unmarshal(bytes, &manifest)
	if err != nil {
		return nil, err
	}

	return &manifest, nil
}
