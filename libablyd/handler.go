package libablyd

import (
	"errors"
	"strconv"
	"encoding/json"

	"github.com/ably/ably-go/ably"
	"context"
	"sync"
)

var ScriptNotFoundError = errors.New("script not found")


type AblyDHandler struct {
	config *Config
	ablyRealtime   *ably.Realtime // Ably Realtime Instance
	ablyCommandChannel *ably.RealtimeChannel

	ablyDState AblyDInstanceState // Active instances running

	log *LogScope

	command string
}

// TODO: Replace s
func NewAblyDHandler(ablyRealtime *ably.Realtime, config *Config, newLog *LogScope) (ablyDHandler *AblyDHandler, err error) {
	ablyDHandler = &AblyDHandler{config: config, ablyRealtime: ablyRealtime}
	ablyDHandler.command = config.CommandName

	ablyDHandler.ablyCommandChannel = startUpCommandChannel(ablyRealtime)
	ablyDHandler.ablyDState = AblyDInstanceState{config.MaxForks, make(map[int]string)}

	ablyDHandler.enterPresence()

	ablyDHandler.log = newLog
	ablyDHandler.log.Associate("command", ablyDHandler.command)

	return ablyDHandler, nil
}

func (ablyDHandler *AblyDHandler) enterPresence() {
	ablyDHandler.ablyCommandChannel.Presence.Enter(context.Background(), ablyDHandler.ablyDState)
}

func startUpCommandChannel(ablyRealtime *ably.Realtime) (ablyCommandChannel *ably.RealtimeChannel) {
	commandChannel := ablyRealtime.Channels.Get("command")

	return commandChannel
}

func (ablyDHandler *AblyDHandler) ListenForCommands(wg *sync.WaitGroup) {
	// Subscribe to messages sent on the channel
	ablyDHandler.ablyCommandChannel.SubscribeAll(context.Background(), func(msg *ably.Message) {
		stringData := msg.Data.(string)
	    data := AblyDInstanceStartMessage{}
	    json.Unmarshal([]byte(stringData), &data)

		if data.Action == "stop" {
			wg.Done()
		} else if data.Action != "" {
			if (ablyDHandler.ablyDState.MaxInstances <= len(ablyDHandler.ablyDState.Instances)) {
				ablyDHandler.ablyCommandChannel.Publish(context.Background(), "Error", "Failed to create new instance: Max instances reached")
			} else {
				// TODO: Should use msg ID not data, 
				// but this does not currently work https://github.com/ably/ably-go/issues/58
				go ablyDHandler.Accept(data.MessageID, data.Args)
			}
		}
	})
	ablyDHandler.log.Access("ablyD", "READY")
}

func (ablyDHandler *AblyDHandler) Accept(messageID string, args []string) {
	allArgs := append(ablyDHandler.config.CommandArgs, args...)
	launched, err := launchCmd(ablyDHandler.command, allArgs, ablyDHandler.config.Env)
	if err != nil {
		ablyDHandler.log.Error("process", "Could not launch process (%s)", err)
		return
	}

	log := ablyDHandler.log

	pid := strconv.Itoa(launched.cmd.Process.Pid)

	log.Associate("pid", pid)

	log.Access("session", "CONNECT")
	defer log.Access("session", "DISCONNECT")

	process := NewProcessEndpoint(launched, ablyDHandler.log)

    channelOutput := ablyDHandler.ablyRealtime.Channels.Get(pid + ":serveroutput")
	channelInput := ablyDHandler.ablyRealtime.Channels.Get("[?rewind=10]"+ pid + ":serverinput")

	// Enter presence of serverinput
	channelInput.Presence.Enter(context.Background(), "")

	ablyEndpoint := NewAblyEndpoint(channelInput, channelOutput, log)


	newInstanceMessage := &NewInstanceMessage{MessageID: messageID, Pid: pid}

	ablyDHandler.ablyCommandChannel.Publish(context.Background(), "new-instance", newInstanceMessage)

	// Add to our list of active instances
	ablyDHandler.ablyDState.Instances[launched.cmd.Process.Pid] = "Running"
	ablyDHandler.ablyCommandChannel.Presence.Update(context.Background(), ablyDHandler.ablyDState)

	PipeEndpoints(process, ablyEndpoint)

	// TODO: Remove from list here and update presence
	delete(ablyDHandler.ablyDState.Instances, launched.cmd.Process.Pid)
	ablyDHandler.ablyCommandChannel.Presence.Update(context.Background(), ablyDHandler.ablyDState)
	channelInput.Presence.Leave(context.Background(), "")
}

