# goNaverWorksBot

This is a robust Naver Works Bot application built with Go. It provides the essential infrastructure for secure authentication and real-time message handling via webhooks.

## Getting Started
### Prerequisites
- Go (version 1.25+ is recommended)
- A Naver Works Developer Account
- ngrok (for exposing your local server for webhook testing)

### Setup & Run
1. Configuration
    - To begin, copy the sample environment file and update it with your actual credentials.
    Copy the sample file:
    ```bash
    cp .env.sample .env
    ```
    - Place your Key: Ensure your RSA Private Key file is saved at the path specified in .env.

2. Run the Bot
    Execute the `main.go` file to start the server:
    ```bash
    go run main.go
    ```

3. Set up Webhook
    To allow Naver Works to send messages to your local bot, you need to expose your server:
    ```bash
    ngrok http 8080
    ```
    Copy the generated HTTPS Forwarding URL (e.g., https://xxxxxx.ngrok-free.app).

    In the Naver Works Bot Console, set the Webhook URL to your forwarded address plus the handler path:
    https://xxxxxx.ngrok-free.app/webhook

Your bot is now ready to receive messages!
