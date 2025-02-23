package manager

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gh0st17/timetable-go/errtype"

	"golang.org/x/net/html"
)

// Возвращает корневой *html.Node страницы, загруженной по ссылке u
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
		return nil, errtype.ErrNetwork(fmt.Errorf("ошибка формирования заголовка запроса: %s", err))
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 14_6_1) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.5 Safari/605.1.15")
	if resp, err = client.Do(req); err != nil {
		return nil, errtype.ErrNetwork(fmt.Errorf("ошибка приема данных: %s", err))
	}
	defer resp.Body.Close()

	if bytes, err = io.ReadAll(resp.Body); err != nil {
		return nil, errtype.ErrRuntime(fmt.Errorf("ошибка разбора ответа сервера: %s", err))
	}

	if doc, err = html.Parse(strings.NewReader(string(bytes))); err != nil {
		return nil, errtype.ErrRuntime(fmt.Errorf("неверный формат адреса прокси: %s", err))
	} else {
		return doc, nil
	}
}

// Предикат для функции загрузки страницы с повтором
// при неудаче
type LoadPredicate func() (*html.Node, error)

// Ззагружает страницу с повтором при неудаче
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
