package main

/*
=== Утилита wget ===

Реализовать утилиту wget с возможностью скачивать сайты целиком

Программа должна проходить все тесты. Код должен проходить проверки go vet и golint.
*/

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html"
)

// downloadFile скачивает файл по указанному URL и сохраняет его в заданное местоположение.
func downloadFile(filePath string, fileURL string) error {
	// Создаем директорию для файла, если она не существует
	if err := os.MkdirAll(path.Dir(filePath), os.ModePerm); err != nil {
		return err
	}

	// Открываем файл для записи
	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Выполняем HTTP-запрос
	resp, err := http.Get(fileURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Записываем тело ответа в файл
	_, err = io.Copy(outputFile, resp.Body)
	return err
}

// updateLinks парсит HTML и заменяет ссылки на локальные пути.
func updateLinks(baseURL *url.URL, filePath string) error {
	// Открываем HTML файл для чтения
	inputFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	doc, err := html.Parse(inputFile)
	if err != nil {
		return err
	}

	// Открываем HTML файл для записи
	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	// Рекурсивная функция для обработки узлов HTML
	var convert func(*html.Node)
	convert = func(n *html.Node) {
		if n.Type == html.ElementNode {
			// Ищем атрибуты href и src
			for i, attr := range n.Attr {
				if attr.Key == "href" || attr.Key == "src" {
					link, err := baseURL.Parse(attr.Val)
					if err != nil {
						continue
					}
					if link.Host == baseURL.Host {
						relPath := link.Path
						n.Attr[i].Val = path.Join(baseURL.Path, relPath)
					}
				}
			}
		}
		// Рекурсивный вызов для дочерних узлов
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			convert(c)
		}
	}

	// Запускаем конвертацию ссылок
	convert(doc)

	// Записываем измененный HTML в файл
	html.Render(outputFile, doc)

	return nil
}

// extractResourceLinks парсит HTML и извлекает все ссылки на ресурсы.
func extractResourceLinks(baseURL *url.URL, body io.Reader) ([]string, error) {
	var links []string
	tokenizer := html.NewTokenizer(body)

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			// Если достигнут конец файла, возвращаем найденные ссылки
			if tokenizer.Err() == io.EOF {
				return links, nil
			}
			return nil, tokenizer.Err()
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			// Обрабатываем теги a, link, img, script
			switch token.Data {
			case "a", "link", "img", "script":
				for _, attr := range token.Attr {
					if attr.Key == "href" || attr.Key == "src" {
						link, err := baseURL.Parse(attr.Val)
						if err != nil {
							continue
						}
						links = append(links, link.String())
					}
				}
			}
		}
	}
}

// downloadWebPage скачивает страницу и все связанные с ней ресурсы.
func downloadWebPage(baseURL string, root string) error {
	// Выполняем HTTP-запрос для главной страницы
	resp, err := http.Get(baseURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Сохраняем главную страницу
	u, err := url.Parse(baseURL)
	if err != nil {
		return err
	}
	filePath := path.Join(root, u.Host, u.Path)
	if strings.HasSuffix(filePath, "/") || filePath == "" {
		filePath = path.Join(filePath, "index.html")
	}

	// Сохраняем главную страницу
	if err := downloadFile(filePath, baseURL); err != nil {
		return err
	}

	// Конвертируем ссылки в главной странице
	if err := updateLinks(u, filePath); err != nil {
		return err
	}

	// Извлекаем ссылки
	links, err := extractResourceLinks(u, resp.Body)
	if err != nil {
		return err
	}

	// Скачиваем все найденные ссылки
	for _, link := range links {
		fmt.Println("Downloading", link)

		// Проверяем, что длина link больше длины baseURL
		if len(link) < len(baseURL) {
			fmt.Println("Skipping link (too short):", link)
			continue
		}

		// Получаем относительный путь
		relPath, err := url.PathUnescape(link[len(baseURL):])
		if err != nil {
			relPath = link[len(baseURL):]
		}
		if relPath == "" {
			relPath = "index.html"
		}

		// Формируем путь к файлу
		filePath := path.Join(root, u.Host, relPath)

		if err := downloadFile(filePath, link); err != nil {
			fmt.Println("Failed to download", link, ":", err)
		}
	}

	return nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: wget URL [destination]")
		return
	}

	// Первый аргумент - команда wget, его игнорируем.
	baseURL := os.Args[2]
	root := "site"
	if len(os.Args) > 3 {
		root = os.Args[3]
	}

	// Проверяем, существует ли root и является ли он директорией
	fileInfo, err := os.Stat(root)
	if os.IsNotExist(err) {
		// Создаем корневую директорию, если она не существует
		if err := os.MkdirAll(root, os.ModePerm); err != nil {
			fmt.Println("Error creating root directory:", err)
			return
		}
	} else if err != nil {
		fmt.Println("Error accessing root directory:", err)
		return
	} else if !fileInfo.IsDir() {
		fmt.Println("Error:", root, "is not a directory")
		return
	}

	// Загружаем страницу и ресурсы
	if err := downloadWebPage(baseURL, root); err != nil {
		fmt.Println("Error:", err)
	}
}
