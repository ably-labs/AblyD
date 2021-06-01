package libablyd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"os"
	"github.com/ably/ably-go/ably"
)

var ScriptNotFoundError = errors.New("script not found")


type AblyDHandler struct {
	conn *ably.RealtimeClient
	config *Config
	*URLInfo // TODO: I cannot find where it's used except in one single place as URLInfo.FilePath

	command string
}

// TODO: Replace s
func NewAblyDHandler(commandName string, config *Config, log *LogScope) (ablyDHandler *AblyDHandler, err error) {
	client, err := ably.NewRealtimeClient(ably.NewClientOptions(config.AblyKey))
	ablyDHandler = &AblyDHandler{conn: client, config: config}
	ablyDHandler.command = commandName
	// if s.Config.UsingScriptDir {
	// 	ablyDHandler.command = ablyDHandler.URLInfo.FilePath
	// }
	log.Associate("command", ablyDHandler.command)

	return ablyDHandler, nil
}

func (ablyHandler *AblyDHandler) Accept(log *LogScope, channelName string) {
	log.Access("session", "CONNECT")
	defer log.Access("session", "DISCONNECT")
	// launched contains input, output, etc
	launched, err := launchCmd(ablyHandler.command, ablyHandler.config.CommandArgs, ablyHandler.config.Env) // Removed command args, Env at end
	if err != nil {
		log.Error("process", "Could not launch process (%s)", err)
		return
	}

	log.Associate("pid", strconv.Itoa(launched.cmd.Process.Pid))

    // false was binary
	binary := ablyHandler.config.Binary
	process := NewProcessEndpoint(launched, binary, log)

    channelOutput := ablyHandler.conn.Channels.Get(channelName + ":serveroutput")
	channelInput := ablyHandler.conn.Channels.Get(channelName + ":serverinput")

	ablyEndpoint := NewAblyEndpoint(channelInput, channelOutput, log)

	PipeEndpoints(process, ablyEndpoint)
}


// URLInfo - structure carrying information about current request and it's mapping to filesystem
type URLInfo struct {
	ScriptPath string
	PathInfo   string
	FilePath   string
}

// GetURLInfo is a function that parses path and provides URL info according to libwebsocketd.Config fields
func GetURLInfo(path string) (*URLInfo, error) {
	// if !config.UsingScriptDir {
	// 	return &URLInfo{"/", path, ""}, nil
	// }

	parts := strings.Split(path[1:], "/")
	urlInfo := &URLInfo{}

	for i, part := range parts {
		urlInfo.ScriptPath = strings.Join([]string{urlInfo.ScriptPath, part}, "/")
		// urlInfo.FilePath = filepath.Join(config.ScriptDir, urlInfo.ScriptPath)
		isLastPart := i == len(parts)-1
		statInfo, err := os.Stat(urlInfo.FilePath)

		// not a valid path
		if err != nil {
			return nil, ScriptNotFoundError
		}

		// at the end of url but is a dir
		if isLastPart && statInfo.IsDir() {
			return nil, ScriptNotFoundError
		}

		// we've hit a dir, carry on looking
		if statInfo.IsDir() {
			continue
		}

		// no extra args
		if isLastPart {
			return urlInfo, nil
		}

		// build path info from extra parts of url
		urlInfo.PathInfo = "/" + strings.Join(parts[i+1:], "/")
		return urlInfo, nil
	}
	panic(fmt.Sprintf("GetURLInfo cannot parse path %#v", path))
}

