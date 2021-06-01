// Copyright 2013 Joe Walnes and the websocketd team.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"ablyD/libablyd"
	"fmt"
	"os"
	"runtime"

	"github.com/ably/ably-go/ably"
	"github.com/joho/godotenv"
)

func logfunc(l *libablyd.LogScope, level libablyd.LogLevel, levelName string, category string, msg string, args ...interface{}) {
	if level < l.MinLevel {
		return
	}
	// fullMsg := fmt.Sprintf(msg, args...)

	assocDump := ""
	for index, pair := range l.Associated {
		if index > 0 {
			assocDump += " "
		}
		assocDump += fmt.Sprintf("%s:'%s'", pair.Key, pair.Value)
	}

	l.Mutex.Lock()
	// fmt.Printf("%s | %-6s | %-10s | %s | %s\n", libablyd.Timestamp(), levelName, category, assocDump, fullMsg)
	l.Mutex.Unlock()
}

func main() {
	config := parseCommandLine()

	log := libablyd.RootLogScope(config.LogLevel, logfunc)

	if runtime.GOOS != "windows" { // windows relies on env variables to find its libs... e.g. socket stuff
		os.Clearenv() // it's ok to wipe it clean, we already read env variables from passenv into config
	}

	godotenv.Load()

	config.Config.AblyKey = os.Getenv("ABLY_KEY")

	handler, _ := libablyd.NewAblyDHandler(os.Args[1], config.Config, log)

	client, _ := ably.NewRealtimeClient(ably.NewClientOptions(os.Getenv("ABLY_KEY")))

	commandChannel := client.Channels.Get("command")
	sub, err := commandChannel.Subscribe()
	if err != nil {
		panic(err)
	}

	// For each message we receive from the subscription, print it out
	for msg := range sub.MessageChannel() {
		channelName, ok := msg.Data.(string)
		if ok {
			go handler.Accept(log, channelName)
		}
	}
}
