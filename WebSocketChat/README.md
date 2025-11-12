# WebSocketChat

Simple chat demo using gorilla/websocket.

- Server: serves static `./static` and `/ws` endpoint for websocket connections.
- Client: `static/index.html` is a minimal browser client.

Run locally:

```bash
cd WebSocketChat
go run .
# open http://localhost:8081/
```

Docker: you can containerize similarly to other modules in this workspace.
