# Logsync Server

## Config

There are two options to configure the server:

1. Create a config.yaml file in one of these locations
    1. next to the executable
    2. in ~/.config/logsync/config.yaml
    3. in ~/.logsync
2. Pass the option as environment variable prefixed with LOGSYNC_CLIENT_

### Options

#### sync.graphs (LOGSYNC_CLIENT_SYNC_GRAPHS)

__required__ \
Provide the paths to all graph directories as an array. The directory should already exist and specifies the graph name.

#### sync.once (LOGSYNC_CLIENT_SYNC_ONCE)

If set to true, the client will sync once and then quit \
default: false

#### sync.interval (LOGSYNC_CLIENT_SYNC_INTERVAL)

Interval that specifies, how often the sync should be performed in seconds \
default: 60

#### encryption.enabled (LOGSYNC_CLIENT_ENCRYPTION_ENABLED)

If set to true, the encryption.key will be taken to encrypt and decrypt the files end to end on the client. \
default: false

#### encryption.key (LOGSYNC_CLIENT_ENCRYPTION_KEY)

Specify the key for the aes encryption. Is mandatory, if the encryption is enabled \
default: ""

#### server.host (LOGSYNC_CLIENT_SERVER_HOST)

__required__ \
Url where the server is located in the format http(s)://server:port

#### server.apitoken (LOGSYNC_CLIENT_SERVER_APITOKEN)
Specify the apitoken for the server, if needed. 
