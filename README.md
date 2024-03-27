# Logsync

## What is this?
Logsync is a simple, self-hostable solution to sync files between two or more clients. It is specifically designed for LogSeq qraphs, but can basically be used for anything.

## Why does it exist?
Because it always is a pain for me to sync my LogSeq graphs between two devices in a simple and secure way. \

## Why not just use Dropbox, Onedrive...?
Those are probably the easier and better solutions, but I also wanted a simple solution for end-to-end encryption. While there
definitely are ways to do this with the given providers, it sounded like a nice challenge to develop for myself.

## How does it work?
Logsync is a simple client/server application. The server component just stores the (encrypted) files for a specific graph and
remembers all the changes in a SQLite database. It provides the JSON API endpoints to check for changes since a specific time and to get the current file content.

The client first tries to download any changes from the server that happened since the last sync and then stores a representation of those files to the disk.
On the next sync it will compare the locally stored graph with the real graph to find changes and send them to the server.

## How to set up?
### Server
You can build the server yourself for the operating system you are using. There is also a Dockerfile available. Check the 
Check the [server docs](server/Readme.md) for the options to configure.

### Client
Check the [client docs](client/Readme.md) for the available options. You can configure the client with a yaml file like this:
```yaml
sync:
  graphs:
    /home/soerenchrist/graphs/Personal
  once: true
  interval: -1
server:
  host: http://myserver.com:8080
encryption:
  enabled: true
  key: "MySuperSecureEncryptionKey"
```
Put the yaml file next to the executable or store the values in environment variables.
For the first client you connect to the server, the graph directory should already have some data, that will then be synced.
The next clients should point to empty directories. They will then fetch the current state from the server.

## Todos
- Conflict handling
- Client logging
- User interface / CLI

## Contributions
Feel free to contribute