// Copyright 2013 Joe Walnes and the websocketd team.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"ablyD/libablyd"
	"sync"

	"github.com/ably/ably-go/ably"
	"github.com/joho/godotenv"
)

func main() {
	config := parseCommandLine()

	log := libablyd.RootLogScope(config.LogLevel, libablyd.Logfunc)

	godotenv.Load()

	client, err := ably.NewRealtime(
		ably.WithKey(config.AblyAPIKey),
		ably.WithEchoMessages(false),
		ably.WithClientID("ablyD"))

	if err != nil {
		log.Error("ablyD", "%s", err)
		return
	}

	config.Config.MaxForks = config.MaxForks
	config.Config.LogLevel = config.LogLevel
	handler, _ := libablyd.NewAblyDHandler(client, config.Config, log)

	var wg sync.WaitGroup
	wg.Add(1)
	handler.ListenForCommands(&wg)
	wg.Wait()
}
