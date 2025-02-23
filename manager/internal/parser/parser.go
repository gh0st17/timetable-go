// Пакет parser предоставляет набор функции для
// разбора исходного кода страницы с расписанием
// или со списком групп
//
// Основные функции:
//   - FindNode: Поиск узлов с нужными
//     параметрами [NodeParam]
package parser

import (
	"strings"

	"github.com/gh0st17/timetable-go/manager/internal/basic_types"

	"golang.org/x/net/html"
)

// Критерии поиска узла
type NodeParam struct {
	Tag       string
	Attr_name string
	Attr_val  string
}

var (
	educator_param = NodeParam{
		Tag:       "a",
		Attr_name: "",
		Attr_val:  "",
	}
)

// Поиск узлов с нужными параметрами [NodeParam]
func FindNode(doc *html.Node, param NodeParam) []html.Node {
	var found []html.Node

	if doc.Type == html.ElementNode && doc.Data == param.Tag {
		if len(doc.Attr) > 0 && param.Attr_name != "" {
			for _, attr := range doc.Attr {
				if attr.Key == param.Attr_name && attr.Val == param.Attr_val {
					found = append(found, *doc)
				}
			}
		} else if param.Attr_name == "" {
			found = append(found, *doc)
		}
	}

	// Обходим все дочерние узлы
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		found = append(found, FindNode(c, param)...)
	}

	return found
}

// Извлекает название предмета из html
func ExtractSubject(html_subj html.Node, subject *basic_types.Subject) {
	event_name_type := ExtractText(&html_subj)
	var (
		event_name string = ""
		event_type string = ""
	)

	var splited = strings.Split(event_name_type, " ")
	for _, s := range splited[:len(splited)-1] {
		event_name += s + " "
	}
	event_type = splited[len(splited)-1]

	subject.Event_name = strings.TrimSpace(event_name)
	subject.Event_type = strings.TrimSpace(event_type)
}

// Извлекает время и место проведения занятия из html
func ExtractPlace(html_place *html.Node, subject *basic_types.Subject) {
	var (
		tmp_str    string
		educs_html = []html.Node{}
	)

	educs_html = FindNode(html_place, educator_param)
	for _, html_edu := range educs_html {
		subject.Educators = append(subject.Educators, ExtractText(&html_edu))
	}

	html_place = html_place.FirstChild.NextSibling
	subject.Event_time = ExtractText(html_place)

	for i := 0; i < len(subject.Educators)*2; i++ {
		html_place = html_place.NextSibling
	}

	for c := html_place.NextSibling; c != nil; c = c.NextSibling {
		tmp_str = ExtractText(c)
		if tmp_str != "" {
			subject.Places = append(subject.Places, tmp_str)
		}
	}
}

// Извлекает текст внутри тэга
func ExtractText(n *html.Node) (result string) {
	if n.Type == html.TextNode {
		trimWhitespaces(&n.Data)
		if n.Data != "" {
			return n.Data
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text := ExtractText(c)
		if text != "" {
			if result != "" && !strings.HasSuffix(result, " ") {
				result += " "
			}
			result += text
		}
	}

	trimWhitespaces(&result)
	return result
}

// Удаляет все табуляции, переводы строк и лишние пробелы
func trimWhitespaces(str *string) {
	*str = strings.ReplaceAll(*str, "\t", "")
	*str = strings.ReplaceAll(*str, "\n", "")
	*str = strings.ReplaceAll(*str, "  ", " ")
	*str = strings.ReplaceAll(*str, "\u00a0", " ")
	*str = strings.TrimSpace(*str)
}
