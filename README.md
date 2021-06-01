# AblyD

This is a basic demo which takes [websocketd](https://github.com/joewalnes/websocketd), which makes command-line commands on a device available to other devices via WebSockets, and provides this functionality through [Ably](https://www.ably.com).

## Setup

Firstly you need to replace the **YOUR_API_KEY** text in `.env.example` with your Ably API key and re-name the file to `.env`. You can sign up for an account with [Ably](https://www.ably.com/) and access your API key from the [app dashboard](https://www.ably.com/accounts/any/apps/any/app_keys). We keep the API key in `.env` and ignore it in `.gitignore` to avoid accidentally sharing the API key.

Now all you need to do is run the main go file, passing in the command line command you want to be run as a parameter! For example, to use a bash file `count.sh`, which is included in this repo's `examples` folder, just run:

```bash
~ $ go run . ./examples/bash/count.sh
```

The program is now running, waiting for a message in the `commands` channel in Ably to be messaged with a namespace to use for publishing the outputs from the script. You can use the `/examples/bash/count.html` file to easily test this out. Replace the `INSERT_API_KEY_HERE` with the same API key used in your main Go function, and load the webpage. If you press the publish button on that page, you should see counting coming from the server!