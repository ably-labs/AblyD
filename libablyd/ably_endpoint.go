// Copyright 2013 Joe Walnes and the websocketd team.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libablyd

import (
	"fmt"
	"github.com/ably/ably-go/ably"
)

type AblyEndpoint struct {
	outputChannel *ably.RealtimeChannel
	inputChannel *ably.RealtimeChannel
	output chan []byte
	log    *LogScope
	mtype  int
}

func NewAblyEndpoint(inputAblyChannel *ably.RealtimeChannel, outputAblyChannel *ably.RealtimeChannel, log *LogScope) *AblyEndpoint {
	endpoint := &AblyEndpoint{
		inputChannel:   inputAblyChannel,
		outputChannel:  outputAblyChannel,
		output: make(chan []byte),
		log:    log,
		mtype:  0,// websocket.TextMessage,
	}
	// if bin {
	// 	endpoint.mtype = websocket.BinaryMessage
	// }
	return endpoint
}

func (ablyEndpoint *AblyEndpoint) Terminate() {
	ablyEndpoint.log.Trace("websocket", "Terminated websocket connection")
}

func (ablyEndpoint *AblyEndpoint) Output() chan []byte {
	return ablyEndpoint.output
}

func (ablyEndpoint *AblyEndpoint) Send(msg []byte) bool {
	_, err := ablyEndpoint.outputChannel.Publish("console msg", msg)

	if err != nil {
		ablyEndpoint.log.Trace("ably", "Cannot send: %s", err)
		return false
	}

	return true
}

func (ablyEndpoint *AblyEndpoint) StartReading() {
	go ablyEndpoint.read_frames()
}

func (ablyEndpoint *AblyEndpoint) read_frames() {
	sub, err := ablyEndpoint.inputChannel.Subscribe()
	if err != nil {
		ablyEndpoint.log.Trace("ably", "Cannot subscribe in read_frames: %s", err)	
	}

	// For each message we receive from the subscription, append to output
	for msg := range sub.MessageChannel() {
		fmt.Println("LOLOLOLOLOL")
	    var b []byte = []byte(msg.Data.(string))
		ablyEndpoint.output <- append(b, '\n')
	}
}
