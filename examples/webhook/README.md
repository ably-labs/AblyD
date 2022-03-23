# Sending messages from AblyD to a webhook

With Ably it's incredibly easy to send messages from channels to any webhook. All you need to do is define a [Rule](https://ably.com/documentation/general/events) which defines which channels to get messages from, and what URL to send them to.

We can do this programatically using Ably's [Control API](https://ably.com/documentation/control-api).  In this sample, we create a Rule programatically which will consume from all potential `serveroutput` channels, and then makes a request to the AblyD command channel to start a process. Assuming you have an AblyD instance running, it will receive this request and start publishing to a new `serveroutput` channel. This will be picked up by the Rule created, and the messages will be sent to the webhook you defined.

You'll need to create a `.env` file matching the values in the `.env.example` file, which also contains details of how to get each required field.

Once you're done using the Rule, you can close it by closing this program with `Ctrl+C`.

# Trying it out

A simple series of steps to try this out would be:

* Start up an instance of AblyD from the base of the directory with `./ablyD ./examples/bash/count.sh`. Make sure you've specified your [your API key](https://www.ably.com/accounts/any/apps/any/app_keys) in your environment, or by passing it in as `ABLYD_API_KEY=ABC123 ./ablyD ./examples/bash/count.sh`.
* Create a .env file in this folder with the required details
* [Create a requestbin endpoint](https://requestbin.com/) as a test endpoint, and add that to the .env file
* Run `npm install` then `node webhook.js`, and the above process should occur!
