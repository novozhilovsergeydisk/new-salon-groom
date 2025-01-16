package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
)

func main() {
	// Проверка аргументов
	if len(os.Args) < 2 {
		fmt.Println("Использование: go run link_grabber.go <URL>")
		return
	}

	// Получение URL из аргументов
	url := os.Args[1]

	// Загружаем содержимое страницы
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Ошибка загрузки страницы: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Ошибка HTTP-запроса: %s\n", resp.Status)
		return
	}

	// Парсим HTML
	fmt.Println("Поиск тегов <img> и их атрибутов...")
	if err := parseAndPrintImgAttributes(resp.Body); err != nil {
		fmt.Printf("Ошибка парсинга HTML: %v\n", err)
	}
}

// parseAndPrintImgAttributes парсит HTML и выводит значения атрибутов src и srcset из тегов <img>
func parseAndPrintImgAttributes(r io.Reader) error {
	tokenizer := html.NewTokenizer(r)

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			if tokenizer.Err() == io.EOF {
				return nil
			}
			return tokenizer.Err()
		case html.StartTagToken:
			token := tokenizer.Token()
			if token.Data == "img" {
				for _, attr := range token.Attr {
					if attr.Key == "src" || attr.Key == "srcset" {
						fmt.Printf("%s: %s\n", attr.Key, attr.Val)
					}
				}
			}
		}
	}
}
