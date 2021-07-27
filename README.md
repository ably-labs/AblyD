# AblyD

This is a basic demo which takes [websocketd](https://github.com/joewalnes/websocketd), which makes command-line commands on a device available to other devices via WebSockets, and provides this functionality through [Ably](https://www.ably.com).

## Setup

Firstly you need to get an Ably API key. You can sign up for an account with [Ably](https://www.ably.com/) and access your API key from the [app dashboard](https://www.ably.com/accounts/any/apps/any/app_keys). 

Now all you need to do is run the main go file, passing in the command line command you want to be run as a parameter! For example, to use a bash file `count.sh`, which is included in this repo's `examples` folder, just run:

```bash
~ $ go run . --apikey=YOUR_API_KEY ./examples/bash/count.sh
```

The program is now running, waiting for a message in the `command` channel in Ably to be sent. The message's data field should match the structure, with the value `start` for the message's name:

```json
{
  "MessageID": "unique string value",
  "Args": [ "some", "additional", "args", "for", "the", "programs" ]
}
```

Once the server receives a message in the `command` channel, it will start up an instance of the program, using the `Args` you specify in the message as well. Once the program has started up, the server will send a message onto the `command` channel of structure:

```json
{
	"MessageID": "unique string value matching the requesting MessageID",
	"Pid": "unique identifier for the process"
}
```

The client can identify the instance which has started for them by the `MessageID`, then use the `Pid` to connect to an input and an output channel for the process. These will be of structure `{Pid}:serverinput` and `{Pid}:serveroutput`.

Subscribing to the `serveroutput` channel will allow the client to receive any stdout messages from the server. The client can also publish messages into the `serverinput` channel which will be passed into the stdin of the process.

This will continue until the program naturally terminates, resulting in the process dying, or the client submits a message to the `serverinput` with data `KILL`.

## Checking Current State of Processes and an AblyD Instance

AblyD makes use of Ably Presence to identify what AblyD instances exist, and what processes are running on each instance. If you check the presence set of the `command` channel, you'll see each currently active instance present with the following attached data:

```json
{
	"MaxInstances": 20,
	"Instances": {
		"3490348": "Running",
		"Another PID": "Running"
	}
}
```

This indicates the most processes you can have, and the currently active processes. The server will also enter the presence set of any `serverinput` channels to indicate it is actively listening to them.

## Testing

You can use the `/examples/bash/count.html` file to easily test this out. Replace the `INSERT_API_KEY_HERE` with the same API key used in your main Go function, and load the webpage. If you press the publish button on that page, you should see counting coming from the server!