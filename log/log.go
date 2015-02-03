package log

import (
	"fmt"
	"github.com/wsxiaoys/terminal/color"
	"io"
	"log"
	"os"
)

const (
	// Color settings for the different logging functions
	// See: https://godoc.org/github.com/wsxiaoys/terminal/color
	defaultColor = "@w" // white
	infoColor    = "@c" // cyan
	warnColor    = "@y" // yellow
	successColor = "@g" // green
	errorColor   = "@r" // red
)

var (
	Default = NewLogger(os.Stdout, defaultColor)
	Info    = NewLogger(os.Stdout, infoColor)
	Warn    = NewLogger(os.Stdout, warnColor)
	Success = NewLogger(os.Stdout, successColor)
	Error   = NewLogger(os.Stdout, errorColor)
)

type Logger struct {
	out   io.Writer
	color string
}

func NewLogger(out io.Writer, color string) *Logger {
	return &Logger{
		out:   out,
		color: color,
	}
}

func (l *Logger) Print(v ...interface{}) {
	log.SetOutput(l.out)
	log.Print(color.Sprint(l.color + fmt.Sprint(v...)))
}

func (l *Logger) Println(v ...interface{}) {
	log.SetOutput(l.out)
	log.Println(color.Sprint(l.color + fmt.Sprint(v...)))
}

func (l *Logger) Printf(format string, v ...interface{}) {
	log.SetOutput(l.out)
	log.Printf(color.Sprint(l.color + fmt.Sprintf(format, v...)))
}

func (l *Logger) Panic(v ...interface{}) {
	log.SetOutput(l.out)
	log.Panic(color.Sprint(l.color + fmt.Sprint(v...)))
}

func (l *Logger) Panicln(v ...interface{}) {
	log.SetOutput(l.out)
	log.Panicln(color.Sprint(l.color + fmt.Sprint(v...)))
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	log.SetOutput(l.out)
	log.Panicf(color.Sprint(l.color + fmt.Sprintf(format, v...)))
}

func (l *Logger) Fatal(v ...interface{}) {
	log.SetOutput(l.out)
	log.Fatal(color.Sprint(l.color + fmt.Sprint(v...)))
}

func (l *Logger) Fatalln(v ...interface{}) {
	log.SetOutput(l.out)
	log.Fatalln(color.Sprint(l.color + fmt.Sprint(v...)))
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	log.SetOutput(l.out)
	log.Fatalf(color.Sprint(l.color + fmt.Sprintf(format, v...)))
}
