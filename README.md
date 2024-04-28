### Building and running your application

Rename the `.env.example` file to `.env` and fill in the DRPC api key.
If you need, set up your own stream name.
The default one is `ethereum_events`.

When you're ready, start your application by running:
`docker compose up --build`.

NATS JetStream should be up and running on port 4222 (NATS address: `nats://localhost:4222`).
