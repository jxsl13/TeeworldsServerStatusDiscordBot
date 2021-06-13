# TeeworldsServerStatusDiscordBot


## Requirements

- Docker
- Docker Compose

or

- Go compiler

(can also be used as standalone executable)

## Configuration:
Fill in the values in the `sample.env` file and rename it to `.env`
```
DISCORD_TOKEN=ODUy....<TOKEN>
DISCORD_CHANNEL_ID=7188....<CHANNEL ID>
DICORD_OWNER=username#1234
TEEWORLDS_SERVERS=92.42.44.64:8303,89.163.148.121:8305,89.163.148.121:8303,89.163.148.121:8306,89.163.148.121:8304
REFRESH_INTERVAL=60s
CUSTOM_FLAGS=flag-><:custom_emoji:123456>
```

## Deploy

### Docker Compose
Execute the first target in the `Makefile` with the following command.
```
make
```

### Standalone executable
Compile the project
```
make build
```

Start the application
````
./TeeworldsServerStatusDiscordBot --env-file .env
```


