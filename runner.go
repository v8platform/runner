package runner

import (
	"github.com/v8platform/errors"
	"github.com/v8platform/find"
	"github.com/v8platform/runner/cmd"
	"strings"

	"context"
)

const CreateInfobase = "CREATEINFOBASE"

var defaultVersion = "8.3"

type PlatformRunner interface {

	//CreateInfobase()
	Run(ctx context.Context) error
	Background(ctx context.Context) (Process, error)
	Check() error
	Args() []string
	Opts() Options
}

func NewPlatformRunner(where Infobase, what Command, opts ...interface{}) PlatformRunner {

	runner := newRunner(where, what, opts...)

	return &runner
}

type platformRunner struct {
	Options *Options
	Where   Infobase
	What    Command
	command string
	args    []string
}

func (r *platformRunner) Run(ctx context.Context) error {

	p, err := r.Background(ctx)

	if err != nil {
		return err
	}

	return <-p.Wait()
}

func (r *platformRunner) Background(ctx context.Context) (Process, error) {

	if err := r.Check(); err != nil {
		return nil, err
	}

	p := r.background(ctx)

	return p, nil
}

func (r *platformRunner) Check() error {

	_, err := getV8Path(*r.Options)

	if err != nil {
		return err
	}

	return checkCommand(r.What)
}

func (r *platformRunner) Args() []string {

	commandV8, _ := getV8Path(*r.Options)
	return append([]string{
		commandV8}, r.args...)

}

func (r *platformRunner) Opts() Options {
	return *r.Options
}

func newRunner(where Infobase, what Command, opts ...interface{}) platformRunner {

	options := defaultOptions()

	inlineOptions := getOptions(opts...)
	if inlineOptions != nil {
		options = inlineOptions
	}

	o := clearOpts(opts)

	options.Options(o...)

	args := getCmdArgs(where, what, *options)

	r := platformRunner{
		Where:   where,
		What:    what,
		Options: options,
		args:    args,
	}

	return r
}

func Run(where Infobase, what Command, opts ...interface{}) error {

	return NewPlatformRunner(where, what, opts...).Run(context.Background())

}

func Background(ctx context.Context, where Infobase, what Command, opts ...interface{}) (Process, error) {

	return NewPlatformRunner(where, what, opts...).Background(ctx)

}

func (r *platformRunner) background(ctx context.Context) Process {

	if r.Options.Context == nil {
		r.Options.Context = ctx
	}

	cmdRunner := cmd.NewCmdRunner(r.command, r.args,
		cmd.WithContext(r.Options.Context),
		cmd.WithOutFilePath(r.Options.Out),
		cmd.WithDumpResultFilePath(r.Options.DumpResult),
	)

	p := background(cmdRunner, ctx)

	return p

}

func checkCommand(what Command) (err error) {
	err = what.Check()
	return
}

func isCreateInfobase(what Command) bool {

	return strings.EqualFold(strings.ToUpper(what.Command()), CreateInfobase)

}

func getConnectionsStringParams(values []string) (params []string, additional []string) {

	for _, value := range values {

		if strings.HasPrefix(value, "/") || strings.HasPrefix(value, "-") {
			additional = append(additional, value)
		} else {
			params = append(params, value)
		}
	}

	return
}

func joinConnectionStringParams(whereParams, whatParams []string) string {

	// TODO Сделать поиск одинаковых параметров
	params := append(whereParams, whatParams...)
	return strings.Join(params, ";")
}

func getCmdArgs(where Infobase, what Command, options Options) []string {

	params := &Params{
		values: []string{what.Command()},
	}

	if isCreateInfobase(what) {

		connectionStringParams, values := getConnectionsStringParams(what.Values())

		connectionString := strings.Join(connectionStringParams, ";")
		params.Append(connectionString)
		params.Append(values...)

	} else {

		params.Append(where.ConnectionString())
		params.Append(what.Values()...)

	}

	params.Append(options.commonValues...)
	params.Append(options.Values()...)

	return params.Values()
}

func getV8Path(options Options) (string, error) {
	if len(options.v8path) > 0 {
		return options.v8path, nil
	}

	v8 := defaultVersion
	if len(options.Version) > 0 {
		v8 = options.Version
	}

	v8path, err := find.PlatformByVersion(v8, find.WithBitness(find.V8_x64x32))

	if err != nil {

		err = errors.NotExist.Newf("Version %s not found", options.Version)
		_ = errors.AddErrorContext(err, "version", options.Version)

		return "", err
	}

	return v8path, nil

}

func defaultOptions() *Options {

	options := Options{}

	options.NewOutFile()
	options.NewDumpResultFile()
	//options.customValues = *types.NewValues()
	//options.commonValues = *types.NewValues()

	return &options
}

func getOptions(opts ...interface{}) *Options {

	for _, opt := range opts {

		switch o := opt.(type) {

		case Options:
			return &o
		case *Options:
			return o
		}

	}

	return nil
}

func clearOpts(opts []interface{}) []Option {

	var o []Option

	for _, opt := range opts {

		if fn, ok := opt.(Option); ok {
			o = append(o, fn)
		}
	}
	return o
}
