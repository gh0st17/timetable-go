// Пакет filesystem предоставляет набор функции для
// работы с файловой системой
//
// Основные функции:
//   - SaveCookiesToFile: Сохраняет куки в текстовый файл
//   - LoadCookiesFromFile: Загружает куки из текстового файла
package filesystem

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gh0st17/timetable-go/errtype"
)

// Сохраняет куки в текстовый файл
func SaveCookiesToFile(jar http.CookieJar, filename string, u *url.URL) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Сохраняем куки в файл
	for _, cookie := range jar.Cookies(u) {
		_, err := file.WriteString(fmt.Sprintf("%s\t%s\t%s\t%s\t%s\n", cookie.Name, cookie.Value, cookie.Path, cookie.Domain, cookie.Expires.Format(time.RFC1123)))
		if err != nil {
			return err
		}
	}
	return nil
}

// Загружает куки из текстового файла
func LoadCookiesFromFile(jar http.CookieJar, filename string, u *url.URL) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	var cookies []*http.Cookie
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "\t")
		if len(parts) != 5 {
			continue // Если формат строки неверный, пропускаем её
		}

		// Парсим данные из строки
		name := parts[0]
		value := parts[1]
		path := parts[2]
		domain := parts[3]
		expires, err := time.Parse(time.RFC1123, parts[4])
		if err != nil {
			return err
		}

		// Создаем куку и добавляем её в список
		cookie := &http.Cookie{
			Name:    name,
			Value:   value,
			Path:    path,
			Domain:  domain,
			Expires: expires,
		}
		cookies = append(cookies, cookie)
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	// Добавляем куки в cookie jar
	jar.SetCookies(u, cookies)
	return nil
}

// Функция для записи строки в файл
func WriteString(filePath string, data string) error {
	var (
		file *os.File
		err  error
	)

	if file, err = os.Create(filePath); err != nil {
		return errtype.ErrRuntime(fmt.Errorf("ошибка создания файла %s: %s", filePath, err))
	}
	defer file.Close()

	if _, err = file.WriteString(data); err != nil {
		return errtype.ErrRuntime(fmt.Errorf("ошибка записи в файл %s: %s", filePath, err))
	}

	return nil
}

// Функция для получения абсолютного пути запускаемой программы
func GetWd() (string, error) {
	if executable, err := os.Getwd(); err != nil {
		return "", err
	} else {
		return filepath.Abs(executable)
	}
}
