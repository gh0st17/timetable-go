package manager

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
	"timetable/errtype"

	"golang.org/x/net/html"
)

func loadFromUrl(u *url.URL, jar http.CookieJar, proxyUrl *url.URL) (*html.Node, error) {
	var (
		bytes []byte
		req   *http.Request
		resp  *http.Response
		doc   *html.Node
		err   error
	)

	client := &http.Client{
		Jar: jar,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
		Timeout: 10 * time.Second,
	}

	if req, err = http.NewRequest("GET", u.String(), nil); err != nil {
		return nil, errtype.NetworkError(fmt.Errorf("ошибка формирования заголовка запроса: %s", err))
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_6_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15")
	if resp, err = client.Do(req); err != nil {
		return nil, errtype.NetworkError(fmt.Errorf("ошибка приема данных: %s", err))
	}
	defer resp.Body.Close()

	if bytes, err = io.ReadAll(resp.Body); err != nil {
		return nil, errtype.RuntimeError(fmt.Errorf("ошибка разбора ответа сервера: %s", err))
	}

	if doc, err = html.Parse(strings.NewReader(string(bytes))); err != nil {
		return nil, errtype.RuntimeError(fmt.Errorf("неверный формат адреса прокси: %s", err))
	} else {
		return doc, nil
	}
}

func saveCookiesToFile(jar http.CookieJar, filename string, u *url.URL) error {
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

func loadCookiesFromFile(jar http.CookieJar, filename string, u *url.URL) error {
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

type LoadPredicate func() (*html.Node, error)

func retryLoadFromUrl(attempts int8, print bool, pred LoadPredicate) (*html.Node, error) {
	var (
		doc   *html.Node
		err   error
		sleep uint64 = 5
	)

	var i int8
	for i = 0; i < attempts; i++ {
		if doc, err = pred(); err == nil {
			return doc, nil
		}

		if print {
			fmt.Printf("Повторная попытка %d через %d секунд\n", i+1, sleep)
		}
		time.Sleep(time.Second * time.Duration(sleep))
		sleep <<= 1
	}

	return nil, err
}
