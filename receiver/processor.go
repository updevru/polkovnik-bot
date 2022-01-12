package receiver

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"polkovnik/app"
	"polkovnik/domain"
	"text/template"
)

type TemplateDto struct {
	Method string
	Body   interface{}
	Params url.Values
	Header http.Header
}

func CreateTemplateDto(receiver *domain.Receiver, r *http.Request) TemplateDto {
	var data TemplateDto
	var err error
	var post *interface{}
	var body []byte

	data = TemplateDto{
		Method: r.Method,
		Header: r.Header,
		Params: r.URL.Query(),
	}

	if r.Method != http.MethodGet {
		body, err = ioutil.ReadAll(r.Body)
		if err == nil && len(body) > 0 {
			var bodyType string
			if receiver.Format == domain.DataReceiverFormatAuto {
				bodyType = r.Header.Get("Content-Type")
			} else {
				bodyType = fmt.Sprintf("application/%s", receiver.Format)
			}

			switch bodyType {
			case "application/json":
				err = json.Unmarshal(body, &post)
				if err == nil {
					data.Body = post
				}
			case "application/xml":
				err = xml.Unmarshal(body, &post)
				if err == nil {
					data.Body = post
				}
			case "application/text":
				data.Body = body
			}
		}

	}

	return data
}

func renderData(templateString string, dto TemplateDto) (string, error) {
	buf := &bytes.Buffer{}
	tpl, err := template.New("test").Funcs(
		template.FuncMap{
			"getValue": func(key string, values map[string][]string) string {
				if values == nil {
					return ""
				}
				vs := values[key]
				if len(vs) == 0 {
					return ""
				}
				return vs[0]
			},
		},
	).Parse(templateString)

	if err != nil {
		return "", err
	}

	err = tpl.Execute(buf, dto)

	return buf.String(), err
}

type Processor struct {
	Tpl *app.TemplateEngine
}

func (p Processor) Run(team *domain.Team, receiver *domain.Receiver, data TemplateDto) error {
	switch receiver.Type {
	case domain.ReceiverTypeMessage:
		return p.sendTeamMessage(team, receiver, data)
	}

	return errors.New(fmt.Sprintf("receiver type %s unknown", receiver.Type))
}
