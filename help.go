// Copyright 2013 Joe Walnes and the websocketd team.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	help = `
{{binary}} ({{version}})

{{binary}} is a command line tool that will allow any executable program
that accepts input on stdin and produces output on stdout to be made 
available with Ably.

Usage:

  Start an AblyD instance for COMMAND:
    {{binary}} [options] COMMAND [command args]

Options:

  --help                         Print help and exit.

  --version                      Print version and exit.

  --loglevel=LEVEL               Log level to use (default access).
                                 From most to least verbose:
                                 debug, trace, access, info, error, fatal

BSD license: Run '{{binary}} --license' for details.
`
	short = `
Usage:

  Start an AblyD instance for COMMAND:
    {{binary}} [options] COMMAND [command args]

  Or, show extended help message using:
    {{binary}} --help
`
)

func get_help_message(content string) string {
	msg := strings.Trim(content, " \n")
	msg = strings.Replace(msg, "{{binary}}", HelpProcessName(), -1)
	return strings.Replace(msg, "{{version}}", Version(), -1)
}

func HelpProcessName() string {
	binary := os.Args[0]
	if strings.Contains(binary, "/go-build") { // this was run using "go run", let's use something appropriate
		binary = "ablyd"
	} else {
		binary = filepath.Base(binary)
	}
	return binary
}

func PrintHelp() {
	fmt.Fprintf(os.Stderr, "%s\n", get_help_message(help))
}

func ShortHelp() {
	// Shown after some error
	fmt.Fprintf(os.Stderr, "\n%s\n", get_help_message(short))
}
