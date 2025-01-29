package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
)

func main() {
	// URL файла CSS
	cssURL := "https://10web-site.ai/124/wp-content/plugins/ai-builder-demo-plugin-master/assets/css/fonts.css"

	// Создаем папку fonts
	if err := os.MkdirAll("fonts", os.ModePerm); err != nil {
		fmt.Printf("Ошибка создания папки fonts: %v\n", err)
		return
	}

	// Загружаем CSS
	resp, err := http.Get(cssURL)
	if err != nil {
		fmt.Printf("Ошибка загрузки CSS: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Ошибка: статус HTTP %d\n", resp.StatusCode)
		return
	}

	// Читаем содержимое CSS
	cssData, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Ошибка чтения CSS: %v\n", err)
		return
	}

	// Регулярное выражение для поиска ссылок на файлы шрифтов
	fontRegex := regexp.MustCompile(`url\(["']?([^"')]+\.(eot|ttf|woff|woff2|otf))["']?\)`)

	// Ищем все совпадения
	matches := fontRegex.FindAllStringSubmatch(string(cssData), -1)

	// Если ничего не найдено
	if len(matches) == 0 {
		fmt.Println("Шрифты не найдены.")
		return
	}

	// Базовый URL для преобразования относительных ссылок
	baseURL, err := url.Parse(cssURL)
	if err != nil {
		fmt.Printf("Ошибка парсинга базового URL: %v\n", err)
		return
	}

	// Загружаем и сохраняем найденные шрифты
	for _, match := range matches {
		fontPath := match[1]

		// Преобразуем относительный путь в абсолютный
		fontURL, err := baseURL.Parse(fontPath)
		if err != nil {
			fmt.Printf("Ошибка формирования URL для %s: %v\n", fontPath, err)
			continue
		}

		fmt.Printf("Найден шрифт: %s\n", fontURL.String())

		// Определяем имя файла
		fileName := path.Base(fontURL.Path)

		// Загружаем файл
		resp, err := http.Get(fontURL.String())
		if err != nil {
			fmt.Printf("Ошибка загрузки файла %s: %v\n", fontURL.String(), err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("Ошибка загрузки файла %s: статус HTTP %d\n", fontURL.String(), resp.StatusCode)
			continue
		}

		// Создаем файл в папке fonts
		outFile, err := os.Create(path.Join("fonts", fileName))
		if err != nil {
			fmt.Printf("Ошибка создания файла %s: %v\n", fileName, err)
			continue
		}
		defer outFile.Close()

		// Копируем данные в файл
		_, err = io.Copy(outFile, resp.Body)
		if err != nil {
			fmt.Printf("Ошибка сохранения файла %s: %v\n", fileName, err)
			continue
		}

		fmt.Printf("Файл %s успешно загружен.\n", fileName)
	}

	fmt.Println("Парсинг завершен.")
}
