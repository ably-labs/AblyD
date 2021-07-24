// Copyright 2013 Joe Walnes and the websocketd team.
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package libablyd

import (
	"github.com/ably/ably-go/ably"
	"context"
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
	return endpoint
}

func (ablyEndpoint *AblyEndpoint) Terminate() {
	ablyEndpoint.inputChannel.Detach(context.Background())
	ablyEndpoint.outputChannel.Detach(context.Background())
	ablyEndpoint.log.Trace("websocket", "Terminated websocket connection")
}

func (ablyEndpoint *AblyEndpoint) Output() chan []byte {
	return ablyEndpoint.output
}

func (ablyEndpoint *AblyEndpoint) Send(msg []byte) bool {
	err := ablyEndpoint.outputChannel.Publish(context.Background(), "message", msg)

	if err != nil {
		ablyEndpoint.log.Trace("ably", "Cannot send: %s", err)
		return false
	}

	return true
}

func (ablyEndpoint *AblyEndpoint) StartReading() {
	ablyEndpoint.read_frames()
}

func (ablyEndpoint *AblyEndpoint) read_frames() {
	ablyEndpoint.inputChannel.SubscribeAll(context.Background(), func(msg *ably.Message) {
		var b []byte = []byte(msg.Data.(string))
		ablyEndpoint.output <- append(b, '\n')
	})
}
