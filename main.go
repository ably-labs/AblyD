// Copyright 2013 Joe Walnes and the websocketd team.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"ablyD/libablyd"
	"fmt"
	"sync"

	"github.com/ably/ably-go/ably"
	"github.com/joho/godotenv"
)

func logfunc(l *libablyd.LogScope, level libablyd.LogLevel, levelName string, category string, msg string, args ...interface{}) {
	if level < l.MinLevel {
		return
	}
	fullMsg := fmt.Sprintf(msg, args...)

	assocDump := ""
	for index, pair := range l.Associated {
		if index > 0 {
			assocDump += " "
		}
		assocDump += fmt.Sprintf("%s:'%s'", pair.Key, pair.Value)
	}

	l.Mutex.Lock()
	fmt.Printf("%s | %-6s | %-10s | %s | %s\n", libablyd.Timestamp(), levelName, category, assocDump, fullMsg)
	l.Mutex.Unlock()
}

func main() {
	config := parseCommandLine()

	log := libablyd.RootLogScope(config.LogLevel, logfunc)

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
	handler, _ := libablyd.NewAblyDHandler(client, config.Config, log)

	var wg sync.WaitGroup
	wg.Add(1)
	handler.ListenForCommands(&wg)
	wg.Wait()
}
