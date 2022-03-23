package libablyd

import (
	"strconv"
	"encoding/json"

	"github.com/ably/ably-go/ably"
	"context"
	"sync"
)

type AblyDHandler struct {
	config *ProcessConfig
	ablyRealtime   *ably.Realtime // Ably Realtime Instance
	ablyCommandChannel *ably.RealtimeChannel

	ablyDState AblyDProcessState // Active instances running

	log *LogScope

	command string
}

func NewAblyDHandler(ablyRealtime *ably.Realtime, config *ProcessConfig, newLog *LogScope) (ablyDHandler *AblyDHandler, err error) {
	ablyDHandler = &AblyDHandler{config: config, ablyRealtime: ablyRealtime}
	ablyDHandler.command = config.CommandName

	ablyDHandler.startUpCommandChannel(ablyRealtime)
	ablyDHandler.ablyDState = AblyDProcessState{config.ServerID, config.ChannelNamespace, config.MaxForks, make(map[int]string)}

	ablyDHandler.enterPresence()

	ablyDHandler.log = newLog

	ablyDHandler.log.Associate("command", ablyDHandler.command)

	return ablyDHandler, nil
}

func (ablyDHandler *AblyDHandler) enterPresence() {
	ablyDHandler.ablyCommandChannel.Presence.Enter(context.Background(), ablyDHandler.ablyDState)
}

func (ablyDHandler *AblyDHandler) startUpCommandChannel(ablyRealtime *ably.Realtime) {
	commandChannel := ablyDHandler.config.ChannelNamespace + ":command"
	ablyDHandler.ablyCommandChannel = ablyRealtime.Channels.Get(commandChannel)
}

func (ablyDHandler *AblyDHandler) ListenForCommands(wg *sync.WaitGroup) {
	// Subscribe to messages sent on the channel
	ablyDHandler.ablyCommandChannel.Subscribe(context.Background(), "start", func(msg *ably.Message) {
		stringData := msg.Data.(string)
	    data := AblyDProcessStartMessage{}
	    json.Unmarshal([]byte(stringData), &data)

	    if data.MessageID != "" {
	    	if data.ServerID != "" && data.ServerID != ablyDHandler.config.ServerID {
	    		return
	    	}
			if (ablyDHandler.ablyDState.MaxProcesses <= len(ablyDHandler.ablyDState.Processes)) {
				ablyDHandler.ablyCommandChannel.Publish(context.Background(), "Error", "Failed to create new process: Max processes reached")
			} else {
				// TODO: Should use msg ID not data, 
				// but this does not currently work https://github.com/ably/ably-go/issues/58
				go ablyDHandler.Accept(data.MessageID, data.Args)
			}
		}
	})
	ablyDHandler.ablyCommandChannel.Subscribe(context.Background(), "stop", func(msg *ably.Message) {
		wg.Done()
	})
	ablyDHandler.log.Access("ablyD", "READY")
}

func (ablyDHandler *AblyDHandler) Accept(messageID string, args []string) {
	allArgs := append(ablyDHandler.config.CommandArgs, args...)

	log := RootLogScope(ablyDHandler.config.LogLevel, Logfunc)

	launched, err := launchCmd(ablyDHandler.command, allArgs, ablyDHandler.config.Env)
	if err != nil {
		ablyDHandler.log.Error("process", "Could not launch process (%s)", err)
		return
	}

	pid := strconv.Itoa(launched.cmd.Process.Pid)

	log.Associate("pid", pid)

	log.Access("session", "CONNECT")
	defer log.Access("session", "DISCONNECT")

	process := NewProcessEndpoint(launched, ablyDHandler.log)
	serverOutputChannel :=  ablyDHandler.config.ChannelPrefix + pid + ":serveroutput"
    channelOutput := ablyDHandler.ablyRealtime.Channels.Get(serverOutputChannel)
	channelInput := ablyDHandler.ablyRealtime.Channels.Get("[?rewind=10]"+ pid + ":serverinput")

	// Enter presence of serverinput
	channelInput.Presence.Enter(context.Background(), "")
	channelOutput.Presence.Enter(context.Background(), "")

	ablyEndpoint := NewAblyEndpoint(channelInput, channelOutput, log)

	newProcessMessage := &NewProcessMessage{MessageID: messageID, Pid: pid, 
	Namespace: ablyDHandler.config.ChannelNamespace, ChannelPrefix: ablyDHandler.config.ChannelPrefix}

	ablyDHandler.ablyCommandChannel.Publish(context.Background(), "new-process", newProcessMessage)

	// Add to our list of active processes
	ablyDHandler.ablyDState.Processes[launched.cmd.Process.Pid] = "Running"
	ablyDHandler.ablyCommandChannel.Presence.Update(context.Background(), ablyDHandler.ablyDState)

	PipeEndpoints(process, ablyEndpoint)

	delete(ablyDHandler.ablyDState.Processes, launched.cmd.Process.Pid)
	ablyDHandler.ablyCommandChannel.Presence.Update(context.Background(), ablyDHandler.ablyDState)
	channelInput.Detach(context.Background())
	channelOutput.Detach(context.Background())
}

