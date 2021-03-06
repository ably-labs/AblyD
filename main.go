// Copyright 2013 Joe Walnes and the websocketd team.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"ablyD/libablyd"
	"sync"
	"os"

	"github.com/ably/ably-go/ably"
	"github.com/joho/godotenv"
)

func main() {
	config := parseCommandLine()

	log := libablyd.RootLogScope(config.LogLevel, libablyd.Logfunc)

	godotenv.Load()

	client, err := ably.NewRealtime(
		ably.WithKey(os.Getenv("ABLYD_API_KEY")),
		ably.WithEchoMessages(false),
		ably.WithClientID(config.ServerID))

	if err != nil {
		log.Error("ablyD", "%s", err)
		return
	}

	config.ProcessConfig.MaxForks = config.MaxForks
	config.ProcessConfig.LogLevel = config.LogLevel
	config.ProcessConfig.ServerID = config.ServerID
	config.ProcessConfig.ChannelNamespace = config.ChannelNamespace
	config.ProcessConfig.ChannelPrefix = config.ChannelPrefix

	handler, _ := libablyd.NewAblyDHandler(client, config.ProcessConfig, log)

	var wg sync.WaitGroup
	wg.Add(1)
	handler.ListenForCommands(&wg)
	wg.Wait()
}
