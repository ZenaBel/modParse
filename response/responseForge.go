package response

// DataForge Структура для відповіді API
type DataForge struct {
	ID                             int    `json:"id"`
	GameVersionId                  int    `json:"gameVersionId"`
	MinecraftGameVersionId         int    `json:"minecraftGameVersionId"`
	ForgeVersion                   string `json:"forgeVersion"`
	Name                           string `json:"name"`
	Type                           int    `json:"type"`
	DownloadUrl                    string `json:"downloadUrl"`
	Filename                       string `json:"filename"`
	InstallMethod                  int    `json:"installMethod"`
	Latest                         bool   `json:"latest"`
	Recommended                    bool   `json:"recommended"`
	Approved                       bool   `json:"approved"`
	DateModified                   string `json:"dateModified"`
	MavenVersionString             string `json:"mavenVersionString"`
	VersionJson                    string `json:"versionJson"`
	LibrariesInstallLocation       string `json:"librariesInstallLocation"`
	MinecraftVersion               string `json:"minecraftVersion"`
	AdditionalFilesJson            string `json:"additionalFilesJson"`
	ModLoaderGameVersionId         int    `json:"modLoaderGameVersionId"`
	ModLoaderGameVersionTypeId     int    `json:"modLoaderGameVersionTypeId"`
	ModLoaderGameVersionStatus     int    `json:"modLoaderGameVersionStatus"`
	ModLoaderGameVersionTypeStatus int    `json:"modLoaderGameVersionTypeStatus"`
	McGameVersionId                int    `json:"mcGameVersionId"`
	McGameVersionTypeId            int    `json:"mcGameVersionTypeId"`
	McGameVersionStatus            int    `json:"mcGameVersionStatus"`
	McGameVersionTypeStatus        int    `json:"mcGameVersionTypeStatus"`
	InstallProfileJson             string `json:"installProfileJson"`
}

// ApiResponseForge Основна структура для тіла відповіді
type ApiResponseForge struct {
	Data DataForge `json:"data"`
}
