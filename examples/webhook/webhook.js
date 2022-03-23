const axiosDefault = require('axios');
const Ably = require('ably');
require('dotenv').config()

// This should create a webhook rule, then make a request to the AblyD command channel to start going

// Access keys allow for programatic Ably Account changes, like adding apps
// Access keys can be found and created in https://ably.com/users/access_tokens
const ABLY_CONTROL_KEY = process.env.ABLY_CONTROL_KEY;

// Ably API key allows for interactions with an Ably App's channels
// Get an Ably API key from https://www.ably.com/accounts/any/apps/any/api_keys
const ABLY_API_KEY = process.env.ABLY_API_KEY;
const APPID = ABLY_API_KEY.substr(0, ABLY_API_KEY.indexOf('.'));
const NAMESPACE = process.env.NAMESPACE; // The namespace your AblyD instance is using for channels. Default is ablyd.

const WEBHOOK_URL = process.env.WEBHOOK_URL;

let axios = axiosDefault.create({
  baseURL: 'https://control.ably.net/v1/apps/',
  timeout: 10000,
  headers: {
	    'Content-Type': 'application/json',
	    'Authorization': 'Bearer ' + ABLY_CONTROL_KEY
	  }
});

// Connect to Ably and get the command channel
let ably = new Ably.Realtime(ABLY_API_KEY);
let commandChannel = ably.channels.get(`${NAMESPACE}:command`);

let ruleID; // Defined once we create the rule

function createWebhookRule(callback) {
	let data = {
	  "status": "enabled",
	  "ruleType": "http",
	  "requestMode": "single",
	  "source": {
	  	// All serveroutput messages will go to the endpoint
	    "channelFilter": `${NAMESPACE}:.*:.*:serveroutput`,
	    "type": "channel.message"
	  },
	  "target": {
	    "url": WEBHOOK_URL,
	    "format": "json"
	  }
	}

	axios.post(`${APPID}/rules`, data)
	.then(function (response) {
		ruleID = response.data.id;
		console.log("Rule created! Making request to AblyD to start a new process.")
		callback();
	})
}

function deleteWebhookRule() {
	axios.delete(`${APPID}/rules/${ruleID}`, "")
	.then(function (response) {
		console.log("Webhook rule deleted");
		process.exit();
	})
	.catch(function (err) {
		console.log(err);
		process.exit();
	})
}

function startAblyDProcess() {
	createWebhookRule(makeStartRequest);
}

function makeStartRequest() {
	let MESSAGEID = Math.random().toString();
	let data = {
    "MessageID": MESSAGEID,
    // "ServerID": "SERVER_ID" // Uncomment and specify an active server's ID to only have that server process the request
  };
	commandChannel.publish("start", data);
	console.log("Process started! Any messages output by an AblyD process should now be sent to the webhook endpoint");
	console.log("Use Ctrl+C to quit this program and delete the webhook rule");
}

startAblyDProcess();

function exitHandler() {
	console.log("Deleting webhook rule");
	deleteWebhookRule();
}

// Delete the webhook Rule we created before closing
process.on('SIGINT', exitHandler.bind());