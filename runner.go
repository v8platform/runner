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

type v8Runner struct {
	Options   *Options
	Where     Infobase
	What      Command
	ctx       context.Context
	commandV8 string
}

func newRunner(ctx context.Context, where Infobase, what Command, opts ...interface{}) v8Runner {

	options := defaultOptions()

	inlineOptions := getOptions(opts...)
	if inlineOptions != nil {
		options = inlineOptions
	}

	o := clearOpts(opts)

	options.Options(o...)

	r := v8Runner{
		Where:   where,
		What:    what,
		Options: options,
		ctx:     ctx,
	}

	return r
}

func Run(where Infobase, what Command, opts ...interface{}) error {

	ctx := context.Background()

	p, err := Background(ctx, where, what, opts...)

	if err != nil {
		return err
	}

	return <-p.Wait()
}

func Background(ctx context.Context, where Infobase, what Command, opts ...interface{}) (Process, error) {

	r := newRunner(ctx, where, what, opts...)

	err := checkCommand(r.What)

	if err != nil {
		return nil, err
	}

	r.commandV8, err = getV8Path(*r.Options)

	if err != nil {
		return nil, err
	}

	p := r.run()

	return p, nil

}

func (r *v8Runner) run() Process {

	args := getCmdArgs(r.Where, r.What, *r.Options)

	runner := prepareRunner(r.ctx, r.commandV8, args, *r.Options)

	p := background(runner, r.ctx)

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

		whereParams, whereValues := getConnectionsStringParams(where.Values())
		whatParams, whatValues := getConnectionsStringParams(what.Values())

		connectionString := joinConnectionStringParams(whatParams, whereParams)

		params.Append(connectionString)
		params.Append(whatValues...)
		params.Append(whereValues...)

	} else {

		params.Append(where.ConnectionString())
		params.Append(what.Values()...)

	}

	params.Append(options.commonValues...)
	params.Append(options.Values()...)

	return params.Values()
}

func prepareRunner(ctx context.Context, command string, args []string, options Options) Runner {

	if options.Context == nil {
		options.Context = ctx
	}

	r := cmd.NewCmdRunner(command, args,
		cmd.WithContext(options.Context),
		cmd.WithOutFilePath(options.Out),
		cmd.WithDumpResultFilePath(options.DumpResult),
	)

	return r
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
