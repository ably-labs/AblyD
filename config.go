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
	"github.com/rs/xid"
)

type Config struct {
	MaxForks          int      // Number of allowable concurrent forks
	LogLevel          libablyd.LogLevel
	ServerID		  string
	ChannelNamespace  string
	ChannelPrefix  	  string	 // Ably channel prefix to use
	*libablyd.ProcessConfig
}

func parseCommandLine() *Config {
	var mainConfig Config
	var config libablyd.ProcessConfig

	flag.Usage = func() {}
	flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	// config options
	versionFlag := flag.Bool("version", false, "Print version and exit")
	logLevelFlag := flag.String("loglevel", "access", "Log level, one of: debug, trace, access, info, error, fatal")
	maxForksFlag := flag.Int("maxforks", 20, "Max forks, zero means unlimited")
	channelNamespace := flag.String("namespace", "ablyd", "Ably Channel Namespace")
	serverID := flag.String("serverid", xid.New().String(), "Unique ID for the server")

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

	mainConfig.ServerID = *serverID
	mainConfig.ChannelNamespace= *channelNamespace
	mainConfig.ChannelPrefix = *channelNamespace + ":" + *serverID + ":"
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

	args := flag.Args()

	if len(args) < 1 {
		fmt.Fprintf(os.Stderr, "Please specify COMMAND or provide --dir, --staticdir or --cgidir argument.\n")
		ShortHelp()
		os.Exit(1)
	}

	if len(args) > 0 {
		if path, err := exec.LookPath(args[0]); err == nil {
			config.CommandName = path // This can be command in PATH that we are able to execute
			config.CommandArgs = flag.Args()[1:]
			// config.UsingScriptDir = false
		} else {
			fmt.Fprintf(os.Stderr, "Unable to locate specified COMMAND '%s' in OS path.\n", args[0])
			ShortHelp()
			os.Exit(1)
		}
	}

	mainConfig.ProcessConfig = &config

	return &mainConfig
}
