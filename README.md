# AblyD

This is a basic demo which takes [websocketd](https://github.com/joewalnes/websocketd), which makes command-line commands on a device available to other devices via WebSockets, and provides this functionality through [Ably](https://www.ably.com).

## Setup

Firstly you need to get an Ably API key. You can sign up for an account with [Ably](https://www.ably.com/) and access your API key from the [app dashboard](https://www.ably.com/accounts/any/apps/any/app_keys). 

Now all you need to do is run the main go file, passing in the command line command you want to be run as a parameter! For example, to use a bash file `count.sh`, which is included in this repo's `examples` folder, just run:

```bash
~ $ ./ablyD ./examples/bash/count.sh
```

Make sure you've specified your [your API key](https://www.ably.com/accounts/any/apps/any/app_keys) in your environment, or by passing it in as `ABLYD_API_KEY=ABC123 ./ablyD ./examples/bash/count.sh`.

The program is now running, waiting for a message in the `ablyd:command` channel in Ably to be sent. For example, to start a process with curl, you would send:

```bash
curl -X POST https://rest.ably.io/channels/ablyd:command/messages \
  -u "${API_KEY}" \
  -H 'Content-Type: application/json' \
  --data \
  '{
    "name": "start",
    "data": {
      "MessageID": "unique string value",
      "Args": [ "some", "additional", "args", "for", "the", "programs" ]
    },
    "format":"json"
  }'
```

Once the server receives a message in the `ablyd:command` channel, it will start up an instance of the program, using the `Args` you specify in the message as well. Once the program has started up, the server will send a message onto the `command` channel of structure:

```json
{
	"MessageID": "unique string value matching the requesting MessageID",
	"Pid": "process_id",
	"Namespace": "ablyd",
	"ChannelPrefix": "ablyd:server_id"
}
```

The client can identify the process which has started for them by the `MessageID`, then use the `Prefix` to connect to an input and an output channel for the process. These will be of structure `{Prefix}{Pid}:serverinput` and `{Prefix}{Pid}:serveroutput`.

Subscribing to the `serveroutput` channel will allow the client to receive any stdout messages from the server. The client can also publish messages into the `serverinput` channel which will be passed into the stdin of the process.

This will continue until the program naturally terminates, resulting in the process dying, or the client submits a message to the `serverinput` with data `KILL`.

## Checking Current State of Processes and an AblyD Instance

AblyD makes use of Ably Presence to identify what AblyD instances exist, and what processes are running on each instance. If you check the presence set of the `command` channel, you'll see each currently active process present with the following attached data:

```json
{
	"ServerID": "my-server-id",
	"Namespace": "ablyd",
	"MaxProcesses": 20,
	"Processes": {
		"3490348": "Running",
		"Another PID": "Running"
	}
}
```

This indicates the most processes you can have, and the currently active processes. The server will also enter the presence set of any `serverinput` channels to indicate it is actively listening to them.

## Interacting with an AblyD Instance

To simplify the process of interacting with an AblyD instance, there is currently a client available for NodeJS on [GitHub](https://github.com/ably-labs/Ablyd-client) and on [npm](https://www.npmjs.com/package/ablyd-client).

## Testing

You can use the `/examples/bash/count.html` file to easily test this out. Replace the `INSERT_API_KEY_HERE` with the same API key used in your main Go function, and load the webpage. If you press the publish button on that page, you should see counting coming from the server!