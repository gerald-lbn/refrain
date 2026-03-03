# Refrain

A self-hosted lyrics fetcher for your music library. Refrain watches your music directories for new files, reads track metadata, and automatically downloads matching `.lrc` (synced) or `.txt` (plain) lyrics files from [LRCLIB](https://lrclib.net).

## Features

- Automatic lyrics downloading alongside your music files
- Synced lyrics (`.lrc`) preferred over plain text (`.txt`)
- Real-time file watching via fsnotify (no more periodic scanning)
- Initial full scan at startup, then instant detection of new files
- Parallel track processing with configurable worker count

## Quick Start (Docker)

### Docker Compose

See [`docker-compose.yaml`](docker-compose.yaml) for an example.

### Docker Run

```bash
docker run -d \
  --name refrain \
  -e PUID=1000 \
  -e PGID=1000 \
  -e REFRAIN_LIBRARIES=/music \
  -v /path/to/your/music:/music \
  ghcr.io/gerald-lbn/refrain:latest
```

### Multiple Libraries

Mount each library as a separate volume and list them comma-separated:

```yaml
volumes:
  - /home/user/Music:/music
  - /home/user/Jazz:/jazz
environment:
  - REFRAIN_LIBRARIES=/music,/jazz
```

## Configuration

All configuration is done via environment variables.

| Variable | Default | Description |
|---|---|---|
| `REFRAIN_LIBRARIES` | *(required)* | Comma-separated paths to music directories |
| `REFRAIN_LOG_LEVEL` | `info` | Log level: `debug`, `info`, `warn`, `error` |
| `REFRAIN_APP_WORKERS` | `5` | Number of parallel workers for lyrics fetching |
| `PUID` | `1000` | User ID for file permissions |
| `PGID` | `1000` | Group ID for file permissions |

## Development

```bash
# Run locally
make run

# Run tests
make test

# Run tests with coverage
make coverage

# Build binary
make build
```
