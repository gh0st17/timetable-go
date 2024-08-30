package parser

import (
	"strings"
	"timetable/manager/basic_types"

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

// Поиск узлов
func FindNode(doc *html.Node, found *[]html.Node, param *NodeParam) {
	if doc.Type == html.ElementNode && doc.Data == param.Tag {
		if len(doc.Attr) > 0 && param.Attr_name != "" {
			for _, attr := range doc.Attr {
				if attr.Key == param.Attr_name && attr.Val == param.Attr_val {
					*found = append(*found, *doc)
				}
			}
		} else if param.Attr_name == "" {
			*found = append(*found, *doc)
		}
	}

	// Обходим все дочерние узлы
	for c := doc.FirstChild; c != nil; c = c.NextSibling {
		FindNode(c, found, param)
	}
}

func ExtractSubject(html_subj *[]html.Node, subject *basic_types.Subject) {
	var event_name_type string
	for _, html_s := range *html_subj {
		event_name_type = ExtractText(&html_s)

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
}

func ExtractPlace(html_place *html.Node, subject *basic_types.Subject) {
	var (
		tmp_str    string
		educs_html []html.Node
	)

	FindNode(html_place, &educs_html, &educator_param)
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

// Выкусывем текст
func ExtractText(n *html.Node) string {
	if n.Type == html.TextNode {
		trimWhitespaces(&n.Data)
		if n.Data != "" {
			return n.Data
		}
	}
	var result string
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

func trimWhitespaces(str *string) {
	// Удаляем все табуляции, переводы строк и лишние пробелы
	*str = strings.ReplaceAll(*str, "\t", "")
	*str = strings.ReplaceAll(*str, "\n", "")
	*str = strings.ReplaceAll(*str, "  ", " ")
	*str = strings.TrimSpace(*str)
}
