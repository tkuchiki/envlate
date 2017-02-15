package main

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

type tpl struct {
	unsetErr bool
}

func split(s, sep string) []string {
	return strings.Split(s, sep)
}

func NewTemplate(unset bool) *tpl {
	return &tpl{
		unsetErr: unset,
	}
}

func (t *tpl) render(tplb []byte, envs map[string]string) (string, error) {
	option := "zero"
	if t.unsetErr {
		option = "error"
	}

	tpl := template.New("").Option("missingkey=" + option).Funcs(template.FuncMap{
		"split": split,
	})

	tpl, err := tpl.Parse(string(tplb))
	if err != nil {
		return "", err
	}

	var txt bytes.Buffer
	err = tpl.Execute(&txt, envs)
	if err != nil {
		re := regexp.MustCompile(`template: :([0-9]+):([0-9]+): executing .+: map has no entry for key "(.+)"`)
		group := re.FindStringSubmatch(err.Error())
		if len(group) != 4 {
			return "", err
		}
		err = fmt.Errorf("line %s char %s : no entry for key `%s`", group[1], group[2], group[3])
	}
	return txt.String(), err
}
