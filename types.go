package runner

type Infobase interface {
	// Clear path
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
