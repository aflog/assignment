# assignment-messagebird

This repository contains a simple API that listens to message POST requests and
sends them over to MessageBird through their REST API go library. It takes care
of validating the message input and concatenating the messages that exceed 160
characters.

## Usage
### Run in Docker
To run the API inside a docker:
`APIKEY={{APIKEY}} make run-locally`

Substitute {{APIKEY}} for your MessageBird API key.

### Send a message
Once a service is running you can send a message posting a request with
`{"recipient":"+31612345678","originator":"MessageBird","message":"This is a test message."}`

All fields are required and recipient phone number has to be a valid
international number in E.164 format without spaces.

## Response
The MessageBird API limit (imaginary) is set to 1rps, for this reason a decision
was taken not to wait for the MessageBird API response for returning the 200
HTTP status, Therefore receiving 200 HTTP status means that message data
received in the request are valid and the message was passed to the queue to be
send to MessageBird. This is order to avoid long connections in case of many
concatenated messages and simultaneous API calls.
