# Raychat

trun your raycast pro to a OpenAI API compatible api server

## Installation

0. clone this repo and `cd` into it

1. Install [Raycast](https://raycast.com)

2. Set Fiddler to capture traffic, get the `ClientID` and `ClientSecret`

3. Put the `ClientID` and `ClientSecret` into a `.env` file, fill Email and Password

4. Run `go run main.go` your server will start at `http://localhost:8080`

you can use `http://localhost:8080/v1/chat/completions` to test your server