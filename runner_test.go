package runner

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/v8platform/marshaler"
	"io/ioutil"
	"os"
	"testing"
)

type runTestSuite struct {
	suite.Suite
}

func (b *runTestSuite) SetupSuite() {

}

func (s *runTestSuite) r() *require.Assertions {
	return s.Require()
}

type v8runnerTestSuite struct {
	runTestSuite
}

func Test_runnerTestSuite(t *testing.T) {
	suite.Run(t, new(v8runnerTestSuite))
}

func (t *v8runnerTestSuite) TestCmdRunnerCreateInfobase() {

	runner := NewPlatformRunner(testInfoBase{}, CreateFileInfoBaseOptions{
		File:     "./file_ib",
		DBFormat: "8.3.8",
	})
	args := runner.Args()
	t.r().Contains(args, CreateInfobase)
	t.r().Contains(args, "File='./file_ib';DBFormat=8.3.8")

}

func (t *v8runnerTestSuite) TestCmdRunnerV8path() {

	if testing.Short() {
		t.T().Skip("skipping test in short mode.")
	}

	path, _ := ioutil.TempDir("", "1c_DB_")

	runner := NewPlatformRunner(testInfoBase{}, CreateFileInfoBaseOptions{
		File: path,
	})

	args := runner.Args()
	t.r().Contains(args, CreateInfobase)
	t.r().Contains(args, fmt.Sprintf("File='%s'", path))

	err := runner.Run(context.Background())

	t.r().NoError(err, "no error command")

}

func Exists(name string) (bool, error) {
	_, err := os.Stat(name)
	if os.IsNotExist(err) {
		return false, err
	}
	return true, nil
}

func NewTempIB() testInfoBase {

	path, _ := ioutil.TempDir("", "1c_DB_")

	ib := testInfoBase{
		File: path,
	}

	return ib
}

type testInfoBase struct {

	// имя каталога, в котором размещается файл информационной базы;
	File string `v8:"File, equal_sep, quotes" json:"file"`

	// язык (страна), который будет использован при открытии или создании информационной базы.
	// Допустимые значения такие же как у параметра <Форматная строка> метода Формат().
	// Параметр Locale задавать не обязательно.
	// Если не задан, то будут использованы региональные установки текущей информационной базы;
	Locale string `v8:"Locale, optional, equal_sep" json:"locale"`
}

func (d testInfoBase) ConnectionString() string {

	return "File=" + d.File

}

type CreateInfoBaseOptions struct {
	DisableStartupDialogs bool   `v8:"/DisableStartupDialogs" json:"disable_startup_dialogs"`
	UseTemplate           string `v8:"/UseTemplate" json:"use_template"`
	AddToList             bool   `v8:"/AddToList" json:"add_to_list"`
}

type CreateFileInfoBaseOptions struct {
	CreateInfoBaseOptions `v8:",inherit" json:"common"`

	// имя каталога, в котором размещается файл информационной базы;
	File string `v8:"File, equal_sep, quotes" json:"file"`

	// язык (страна), который будет использован при открытии или создании информационной базы.
	// Допустимые значения такие же как у параметра <Форматная строка> метода Формат().
	// Параметр Locale задавать не обязательно.
	// Если не задан, то будут использованы региональные установки текущей информационной базы;
	Locale string `v8:"Locale, optional, equal_sep" json:"locale"`

	// формат базы данных
	// Допустимые значения: 8.2.14, 8.3.8.
	// Значение по умолчанию — 8.2.14
	DBFormat string `v8:"DBFormat, optional, equal_sep" json:"db_format"`

	// размер страницы базы данных в байтах
	// Допустимые значения:
	//   4096(или 4k),
	//   8192(или 8k),
	//   16384(или 16k),
	//   32768(или 32k),
	//   65536(или 64k),
	// Значение по умолчанию —  4k
	DBPageSize int64 `v8:"DBPageSize, optional, equal_sep" json:"db_page_size"`
}

func (d CreateInfoBaseOptions) Command() string {
	return CreateInfobase
}

func (d CreateInfoBaseOptions) Check() error {

	return nil
}

func (d CreateInfoBaseOptions) Values() []string {

	v, _ := marshaler.Marshal(d)
	return v

}

func (d CreateFileInfoBaseOptions) Values() []string {

	v, _ := marshaler.Marshal(d)
	return v

}
