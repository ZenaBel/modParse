package response

import "time"

// FileHash Структура для хешів файлу
type FileHash struct {
	Value string `json:"value"`
	Algo  int    `json:"algo"`
}

// SortableGameVersion Структура для версій гри
type SortableGameVersion struct {
	GameVersionName        string    `json:"gameVersionName"`
	GameVersionPadded      string    `json:"gameVersionPadded"`
	GameVersion            string    `json:"gameVersion"`
	GameVersionReleaseDate time.Time `json:"gameVersionReleaseDate"`
	GameVersionTypeId      int       `json:"gameVersionTypeId"`
}

// Dependency Структура для залежностей модуля
type Dependency struct {
	ModId        int `json:"modId"`
	RelationType int `json:"relationType"`
}

// Module Структура для модулів файлу
type Module struct {
	Name        string `json:"name"`
	Fingerprint int    `json:"fingerprint"`
}

// FileData Основна структура для даних про файл
type FileData struct {
	Id                   int                   `json:"id"`
	GameId               int                   `json:"gameId"`
	ModId                int                   `json:"modId"`
	IsAvailable          bool                  `json:"isAvailable"`
	DisplayName          string                `json:"displayName"`
	FileName             string                `json:"fileName"`
	ReleaseType          int                   `json:"releaseType"`
	FileStatus           int                   `json:"fileStatus"`
	Hashes               []FileHash            `json:"hashes"`
	FileDate             time.Time             `json:"fileDate"`
	FileLength           int64                 `json:"fileLength"`
	DownloadCount        int                   `json:"downloadCount"`
	FileSizeOnDisk       int64                 `json:"fileSizeOnDisk"`
	DownloadUrl          string                `json:"downloadUrl"`
	GameVersions         []string              `json:"gameVersions"`
	SortableGameVersions []SortableGameVersion `json:"sortableGameVersions"`
	Dependencies         []Dependency          `json:"dependencies"`
	ExposeAsAlternative  bool                  `json:"exposeAsAlternative"`
	ParentProjectFileId  int                   `json:"parentProjectFileId"`
	AlternateFileId      int                   `json:"alternateFileId"`
	IsServerPack         bool                  `json:"isServerPack"`
	ServerPackFileId     int                   `json:"serverPackFileId"`
	IsEarlyAccessContent bool                  `json:"isEarlyAccessContent"`
	EarlyAccessEndDate   time.Time             `json:"earlyAccessEndDate"`
	FileFingerprint      int64                 `json:"fileFingerprint"`
	Modules              []Module              `json:"modules"`
}

// ApiResponse Основна структура для відповіді API
type ApiResponse struct {
	Data []FileData `json:"data"`
}
