package runner

import "fmt"

type Params interface {
	Set(key, sep, value string)
	Append(v2 Params)
	Values() []Param
	StringValues() []string
	String() string
}

type Param interface {
	Key() string
	Sep() string
	Value() string
}

type Infobase interface {
	Path() string
	// Возвращает
	// - /IBConnectionString <СтрокаПодключения>
	// - /F<ПУтьКБазе>
	// - /S<ПутьКСервернойБазе>
	ConnectionString() string
	Values() []string
}

type Command interface {
	Command() string
	Check() error
	Values() []string
}

type Values struct {
	keys   []string
	values map[string]Param
}

func (v Values) Len() int {
	return len(v.keys)
}

type param struct {
	key string
	sep string
	val string
}

func (v param) Key() string {
	return v.key
}
func (v param) Sep() string {
	return v.sep
}
func (v param) Value() string {
	return v.val
}

type ValueSep string

func NewValues() *Values {
	return &Values{
		values: make(map[string]string),
	}
}

const (
	SpaceSep ValueSep = " "
	EqualSep ValueSep = "="
	NoSep    ValueSep = ""
)

func (v *Values) Values() []string {

	var str []string

	for _, value := range v.values {
		str = append(str, fmt.Sprintf("%s%s", value.Sep(), value.Value()))
	}

	return str
}

func (v *Values) Set(key string, sep string, value string) {

	v.Map(param{
		key, sep, value,
	})

}

func (v *Values) Map(val Param) {

	key := val.Key()

	_, ok := v.values[key]

	if !ok {
		v.keys = append(v.keys, key)
	}

	v.values[key] = val

}

func (v *Values) GetMap() {

}

func (v *Values) Append(v2 Params) {

	for _, s2 := range v2.Len() {
		v.Map(s2.key, s2.val)
	}

}
