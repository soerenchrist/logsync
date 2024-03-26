# Logsync Server

## Config
There are two options to configure the server:
1. Create a config.yaml file in one of these locations
   1. next to the executable
   2. in ~/.config/logsync/config.yaml
   3. in ~/.logsync
2. Pass the option as environment variable prefixed with LOGSYNC_

You don't have to configure anything, sane defaults are provided

### Options
#### server.host (LOGSYNC_SERVER_HOST)
Hostname of the server \
default: localhost

#### server.apitoken (LOGSYNC_SERVER_APITOKEN)
Optional api token. Clients will have to send the token in the header as X-Api-Token, otherwise they will receive 401 - Unauthorized

#### server.port (LOGSYNC_SERVER_PORT)
Port of the server \
default: 3000

#### files.path (LOGSYNC_FILES_PATH)
Path to the directory where the files are stored. Will be created if it not exists \
default: ./files/

#### db.path (LOGSYNC_DB_PATH)
Path to the database file (sqlite). Will be created, should not exist. \
default: ./logsync.db

#### logging.level (LOGSYNC_LOGGING_LEVEL)
Log level (debug, info, warn, error). Any other given value will be interpreted as "info" \
default: info
