/*
Color Print(cp) is ANSI coloured text to the standard output by Golang. Vesion 0.0.1

Website https://github.com/kenshin/cprint, depends on https://github.com/daviddengcn/go-colortext

Copyright (c) 2014 Kenshin Wang <kenshin@ksria.com>
*/
package cprint

import (

	// lib
	"3rdparty/src/go-colortext"

	// go
	"fmt"
	"os"
	"reflect"
	"strings"
)

// Print flag, include:
//  - DEBUG
//  - WARING
//  - ERROR
//  - NOTICE
const (
	DEBUG  = 1
	WARING = 2
	ERROR  = 3
	NOTICE = 4
)

// Parse identifying
const SPLIT = "%v"

// Color uint, include:
//  - None
//  - Black
//  - Red
//  - Green
//  - Yellow
//  - Blue
//  - Magenta
//  - Cyan
//  - White
const (
	None = iota
	Black
	Red
	Green
	Yellow
	Blue
	Magenta
	Cyan
	White
)

// Color Print
//  - FgColor : Foreground color
//  - FgBright: Foreground color is it bright?
//  - BgColor : Background color
//  - BgBright: Background color is it bright?
//  - Value   : Color Message
type CP struct {
	FgColor  int
	FgBright bool
	BgColor  int
	BgBright bool
	Value    string
}

// ANSI coloured text to the standard output
//
//  flag   : Include 'DEBUG', 'Waring', 'Error', 'Notice'
//  message: Print content
//  args   : Variable parameter, include string, CP type, when args last value is "\n", auto new line.
//
// For example 1:
//  P(WARING, "Remote latest version %v = latest version %v.\n", param1, param2)
//
// For example 2:
//  cp := CP{1, true, 2, true, localVersion}
//  P(DEBUG, "Current version %v, publish data: ", cp, "2014-05-31")
//
// For example 3:
//  P(DEBUG, "Current version %v", localVersion, "\n")
func P(nflag int, message interface{}, args ...interface{}) {

	// try catch
	defer func() {
		if err := recover(); err != nil {
			Error(ERROR, "util/print.go an error has occurred. Error: ", err)
			os.Exit(0)
		}
	}()

	// set state
	stateColor(nflag)

	// set color message
	msgArr := strings.Split(message.(string), SPLIT)
	for k, v := range msgArr {
		fmt.Print(v)
		if k < len(args) {
			t := reflect.TypeOf(args[k])
			switch t.Name() {
			case "string":
				normalColor(args[k])
			case "CP":
				customColor(args[k])
			default:
				normalColor(args[k])
			}
		}
	}

}

// ANSI coloured erro text to the standard output
//
//  flag   : Include 'DEBUG', 'Waring', 'Error', 'Notice'
//  message: Print content
//  err    : Error content
//
// For example:
//
//  Error(ERROR, "util/print.go an error has occurred. Error: ", err)
func Error(nFlag int, message string, err interface{}) {

	// set flag
	stateColor(nFlag)

	// color message
	ct.ChangeColor(ct.Red, false, ct.Green, false)
	fmt.Printf(message)

	// print err
	fmt.Println(err)

	// reset color
	ct.ResetColor()
}

func stateColor(nState int) {
	switch nState {
	case DEBUG:
		ct.ChangeColor(ct.White, true, ct.Cyan, false)
		fmt.Printf("Debug: ")
	case WARING:
		ct.ChangeColor(ct.White, false, ct.Magenta, false)
		fmt.Printf("Waring: ")
	case ERROR:
		ct.ChangeColor(ct.White, true, ct.Red, false)
		fmt.Printf("Error: ")
	case NOTICE:
		ct.ChangeColor(ct.White, true, ct.Green, false)
		fmt.Printf("Notice: ")
	default:
		//ct.ChangeColor(ct.Blue, false, ct.White, false)
		//fmt.Printf("Notice: ")
	}
	ct.ResetColor()
}

func customColor(cp interface{}) {

	value := reflect.ValueOf(cp)

	fgColor := value.FieldByName("FgColor").Int()
	fgBright := value.FieldByName("FgBright").Bool()
	bgColor := value.FieldByName("BgColor").Int()
	bgBright := value.FieldByName("BgBright").Bool()
	msg := value.FieldByName("Value").String()

	if fgColor > 8 || fgColor < 0 || bgColor > 8 || bgColor < 0 {
		normalColor(msg)
		fmt.Println()
		Error(WARING, "values range error, values range include 0 ~ 8, Error: ", "index out of range")
		return
	}

	ct.ChangeColor(ct.Color(fgColor), fgBright, ct.Color(bgColor), bgBright)
	fmt.Printf(msg)
	ct.ResetColor()
}

func normalColor(msg interface{}) {
	ct.ChangeColor(ct.Green, true, ct.None, false)
	fmt.Printf(msg.(string))
	ct.ResetColor()
}
