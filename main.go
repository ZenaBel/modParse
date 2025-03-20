package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"minecrat-api/manifest"
	"minecrat-api/request"
	"minecrat-api/response"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const url string = "https://api.curseforge.com"

// Функція для безпечного закриття ресурсів
func safeClose(closer io.Closer, resourceName string) {
	if err := closer.Close(); err != nil {
		fmt.Printf("Error closing %s: %v\n", resourceName, err)
	}
}

// Функція для отримання API Key
func getApiKey(cliApiKey string, useEnv bool) (string, error) {
	// Якщо API Key передано через командний рядок, використовуємо його
	if cliApiKey != "" {
		return cliApiKey, nil
	}

	// Якщо ключ не передано через командний рядок і використовується env, пробуємо отримати зі змінної середовища
	if useEnv {
		apiKeyValue := os.Getenv("API_KEY")
		if apiKeyValue == "" {
			return "", errors.New("змінну середовища API_KEY не встановлено")
		}
		return apiKeyValue, nil
	}

	// Якщо немає ні CLI ключа, ні змінної середовища, повертаємо помилку
	return "", errors.New("API Key не вказано ні через прапорець, ні через змінну середовища")
}

// Функція для перевірки існування файлу
func fileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

// Функція для запиту у користувача, чи перезаписувати файл
func askForOverwrite() bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Файл вже існує. Перезаписати? (y/n): ")
	r, _ := reader.ReadString('\n')
	r = strings.TrimSpace(strings.ToLower(r))

	// Перевіряємо варіанти відповіді
	if r == "y" || r == "yes" || r == "н" {
		return true
	}
	return false
}

// Функція для завантаження файлу за URL з можливістю вказати директорію
func downloadFile(downloadUrl string, filename string, directory string, overrode bool) (bool, error) {
	// Створюємо повний шлях до файлу
	fullPath := filepath.Join(directory, filename)

	// Перевіряємо, чи існує файл
	if fileExists(fullPath) {
		// Якщо прапорець -overrode встановлений, пропускаємо запит користувача і перезаписуємо файл
		if !overrode {
			// Запитуємо користувача, чи перезаписувати файл, якщо прапорець не встановлений
			if !askForOverwrite() {
				fmt.Println("Завантаження скасовано.")
				return false, nil
			}
		}
	}

	// Виконуємо HTTP-запит для отримання файлу
	resp, err := http.Get(downloadUrl)
	if err != nil {
		return false, err
	}
	defer safeClose(resp.Body, "file download response body")

	// Перевіряємо, чи існує директорія, якщо ні - створюємо її
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			return false, err
		}
	}

	// Створюємо файл за вказаним шляхом
	out, err := os.Create(fullPath)
	if err != nil {
		return false, err
	}
	defer safeClose(out, "file")

	// Копіюємо дані з відповіді в файл
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return false, err
	}

	fmt.Printf("Файл %s успішно завантажено в директорію %s!\n", filename, directory)
	return true, nil
}

// downloadForgeByVersion Функція для завантаження Forge за версією
func downloadForgeByVersion(versionForge string, apiKey string, overrode bool) {
	req, err := http.NewRequest("GET", url+"/v1/minecraft/modloader/"+versionForge, nil)
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
	}
	defer safeClose(resp.Body, "response body")

	var apiResp response.ApiResponseForge
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
	}
	downloadUrl := apiResp.Data.DownloadUrl
	fmt.Printf("Download URL: %s\n", downloadUrl)

	// Завантажуємо файл
	_, err = downloadFile(downloadUrl, apiResp.Data.Filename, "forge", overrode)
	if err != nil {
		fmt.Printf("Error downloading file: %v\n", err)
	}
}

// Кастомна функція для виведення Usage
func customUsage() {
	_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Користування: %s [опції]\n", os.Args[0])
	fmt.Println("Доступні опції:")
	flag.PrintDefaults()
}

func main() {
	apiKey := flag.String("api-key", "", "API Key для доступу до API (alias -k)")
	flag.StringVar(apiKey, "k", "", "API Key для доступу до API")

	help := flag.Bool("help", false, "Показати це повідомлення допомоги (alias -h)")
	flag.BoolVar(help, "h", false, "Показати це повідомлення допомоги")

	forge := flag.Bool("forge", false, "Завантажити Forge (alias -f)")
	flag.BoolVar(forge, "f", false, "Завантажити Forge")

	noDownload := flag.Bool("no-download", false, "Не завантажувати файли (alias -n)")

	env := flag.Bool("env", false, "Використати змінні середовища для API Key (API_KEY)")

	overrode := flag.Bool("overrode", false, "Перезаписати файли (alias -o)")
	flag.BoolVar(overrode, "o", false, "Перезаписати файли")

	// Перевизначаємо поведінку функції допомоги
	flag.Usage = customUsage

	// Парсимо аргументи командного рядка
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	// Отримуємо API Key
	apiKeyVal, err := getApiKey(*apiKey, *env)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Використовуємо функцію для завантаження маніфесту
	manifestData, err := manifest.LoadManifest("manifest.json")
	if err != nil {
		log.Fatal(err)
	}

	// Виводимо інформацію для перевірки
	fmt.Printf("Manifest Name: %s\n", manifestData.Name)

	fmt.Println("API Key:", apiKeyVal)

	if *forge {
		var versionForge string
		for _, modLoader := range manifestData.Minecraft.ModLoaders {
			if modLoader.Primary {
				versionForge = modLoader.ID
			}
		}
		downloadForgeByVersion(versionForge, apiKeyVal, *overrode)
	}

	if *noDownload {
		fmt.Println("Завантаження файлів вимкнено.")
		return
	}

	var FileIds []int

	// Завантажуємо файли
	for _, file := range manifestData.Files {
		if file.Required {
			FileIds = append(FileIds, file.FileID)
		}
	}

	apiRequest := request.Request{FileIds: FileIds}

	// Створюємо запит
	reqBody, err := json.Marshal(apiRequest)
	if err != nil {
		fmt.Printf("Error marshalling request: %v\n", err)
		return
	}

	req, err := http.NewRequest("POST", url+"/v1/mods/files", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Printf("Error creating request: %v\n", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", apiKeyVal)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending request: %v\n", err)
		return
	}
	defer safeClose(resp.Body, "response body")

	var apiResp response.ApiResponse
	err = json.NewDecoder(resp.Body).Decode(&apiResp)
	if err != nil {
		fmt.Printf("Error decoding response: %v\n", err)
		return
	}

	var success = 0
	var failed = 0
	var listFailed []string

	// Завантажуємо файли
	for _, file := range apiResp.Data {
		fmt.Printf("Завантаження файлу: %s\n", file.FileName)
		res, err := downloadFile(file.DownloadUrl, file.FileName, "mods", *overrode)
		if err != nil {
			fmt.Printf("Помилка завантаження файлу: %v\n", err)
		}
		if res {
			success++
		} else {
			failed++
			listFailed = append(listFailed, file.FileName)
		}
	}

	fmt.Printf("Завантажено файлів: %d\n", success)
	fmt.Printf("Помилок завантаження: %d\n", failed)
	if failed > 0 {
		fmt.Println("Список файлів, які не вдалося завантажити:")
		for _, file := range listFailed {
			fmt.Println(file)
		}
	}
}
