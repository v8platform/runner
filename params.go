package runner

import (
	"regexp"
	"strings"
)

var re = regexp.MustCompile(`(?m)^((\-|\/)\w+)`) // Название переданного аргумента

type Params struct {
	values []string
}

func (p *Params) Values() []string {

	return p.values

}

func findValueName(value string) string {

	name := re.FindString(value)
	return name

}

func (p *Params) Append(arr ...string) {

	if len(p.values) == 0 {
		p.values = append(p.values, arr...)
		return
	}

	for _, value := range arr {
		p.addValue(value)
	}

}

func (p *Params) addValue(value string) {

	name := findValueName(value)

	for i, v := range p.values {

		if len(name) > 0 && strings.HasPrefix(v, name) {
			p.values[i] = value
			return
		}
	}

	p.values = append(p.values, value)

}
