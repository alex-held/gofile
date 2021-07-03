package log

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/onsi/ginkgo/reporters/stenographer/support/go-colorable"
	"github.com/sirupsen/logrus"
	golog "github.com/withmandala/go-log"
)

var (
	YellowFn = color.New(color.FgYellow).SprintFunc()
	RedFn    = color.New(color.FgRed).SprintFunc()
	GreenFn  = color.New(color.FgRed).SprintFunc()

	Out, Err = colorable.NewColorableStdout(), colorable.NewColorableStderr()
	Log      *golog.Logger
	Logger   *logrus.Logger
)

func init() {
	Logger = logrus.StandardLogger()
	Logger.SetOutput(Out)

	Log = golog.New(os.Stdout).
		WithColor().
		WithoutTimestamp()

	Logger.SetFormatter(&logrus.TextFormatter{
		ForceColors:      true,
		DisableTimestamp: true,
		DisableSorting:   true,
	})
}

func Yellow(a ...interface{}) string     { return color.New(color.FgYellow).SprintFunc()(a...) }
func Red() func(a ...interface{}) string { return color.New(color.FgRed).SprintFunc() }
func Green(a ...interface{}) string      { return color.New(color.FgGreen).SprintFunc()(a...) }
func White(a ...interface{}) string      { return color.New(color.FgWhite).SprintFunc()(a...) }
func WhiteHi(a ...interface{}) string    { return color.New(color.FgHiWhite).SprintFunc()(a...) }

func LogErr(format string, args ...interface{}) {
	// format = fmt.Sprintf("%s\t%s", Red("[ERRO]"), format)
	msg := fmt.Sprintf(format, args)
	Logger.SetOutput(Err)
	Logger.Errorf(msg)
	Logger.SetOutput(Err)
}

func LogInfo(format string, args ...interface{}) {
	// format = fmt.Sprintf("%s\t%s", WhiteHi("[INFO]"), fmt.Sprintf(format, args...))
	format = fmt.Sprint(color.CyanString("[INFO]\t"), fmt.Sprintf(format, args...))
	msg := fmt.Sprintf(format, args)
	// Logger.Printf("%s\t%s", WhiteHi("[INFO]"), format)
	Logger.Infof(msg)
}

func LogDebugf(format string, args ...interface{}) {
	// format = fmt.Sprintf("%s\t%s", Green("[DEBU]"), format)
	msg := fmt.Sprintf(format, args)
	Logger.Debugf(msg)
}
