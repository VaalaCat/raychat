# Raychat

trun your raycast pro to a OpenAI API compatible api server

## Installation

you can use docker or run it directly

### easy way

if you have [Docker](https://www.docker.com) installed, you can run this command to start a server

```bash
docker run -dit --name raychat \
	-p 8080:8080 \
	-e EMAIL='your_email' \
	-e PASSWORD='your_password' \
	-e CLIENT_ID='your_client_id' \
	-e CLIENT_SECRET='your_client_secret' \
	-e EXTERNAL_TOKEN='your_fake_openai_token' \ # you can provide multi token like "token_a,token_b", token splitted with comma
	--restart always \
	vaalacat/raychat:latest
```

or if you already have a token, you can run this command

```bash
docker run -dit --name raychat \
	-p 8080:8080 \
	-e TOKEN='your_token' \
	-e EXTERNAL_TOKEN='your_fake_openai_token' \ # you can provide multi token like "token_a,token_b", token splitted with comma
	--restart always \
	vaalacat/raychat:latest
```

then you can use `http://localhost:8080/v1/chat/completions` to test your server, arm and amd64 are both supported

### common way

0. clone this repo and `cd` into it

1. Install [Raycast](https://raycast.com)

2. Set Fiddler to capture traffic, get the `ClientID` and `ClientSecret`

3. Put the `ClientID` and `ClientSecret` into a `.env` file, fill Email and Password

4. Run `go run main.go` your server will start at `http://localhost:8080`

you can use `http://localhost:8080/v1/chat/completions` to test your server

