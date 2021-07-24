// Copyright 2013 Joe Walnes and the websocketd team.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"

	"ablyD/libablyd"
)

type Config struct {
	MaxForks          int      // Number of allowable concurrent forks
	LogLevel          libablyd.LogLevel
	AblyAPIKey		  string
	*libablyd.Config
}

func parseCommandLine() *Config {
	var mainConfig Config
	var config libablyd.Config

	flag.Usage = func() {}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// config options
	versionFlag := flag.Bool("version", false, "Print version and exit")
	logLevelFlag := flag.String("loglevel", "access", "Log level, one of: debug, trace, access, info, error, fatal")
	maxForksFlag := flag.Int("maxforks", 20, "Max forks, zero means unlimited")
	ablyApiKey := flag.String("apikey", "INSERT_API_KEY", "Ably API key")

	err := flag.CommandLine.Parse(os.Args[1:])
	if err != nil {
		if err == flag.ErrHelp {
			PrintHelp()
			os.Exit(0)
		} else {
			ShortHelp()
			os.Exit(2)
		}
	}

	mainConfig.AblyAPIKey = *ablyApiKey
	mainConfig.MaxForks = *maxForksFlag
	mainConfig.LogLevel = libablyd.LevelFromString(*logLevelFlag)

	if mainConfig.LogLevel == libablyd.LogUnknown {
		fmt.Printf("Incorrect loglevel flag '%s'. Use --help to see allowed values.\n", *logLevelFlag)
		ShortHelp()
		os.Exit(1)
	}

	if len(os.Args) == 1 {
		fmt.Printf("Command line arguments are missing.\n")
		ShortHelp()
		os.Exit(1)
	}

	if *versionFlag {
		fmt.Printf("%s %s\n", HelpProcessName(), Version())
		os.Exit(0)
	}

	if mainConfig.AblyAPIKey == "INSERT_API_KEY" {
		fmt.Printf("Please provide your Ably API key with -apikey=API_KEY\n")
		ShortHelp()
		os.Exit(1)
	}

	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Please specify COMMAND or provide --dir, --staticdir or --cgidir argument.\n")
		ShortHelp()
		os.Exit(1)
	}

	if len(args) > 0 {
		if path, err := exec.LookPath(args[0]); err == nil {
			config.CommandName = path // This can be command in PATH that we are able to execute
			fmt.Println(config.CommandName)
			config.CommandArgs = flag.Args()[1:]
			// config.UsingScriptDir = false
		} else {
			fmt.Fprintf(os.Stderr, "Unable to locate specified COMMAND '%s' in OS path.\n", args[0])
			ShortHelp()
			os.Exit(1)
		}
	}

	mainConfig.Config = &config

	return &mainConfig
}
