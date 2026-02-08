#!/bin/sh

PUID=${PUID:-1000}
PGID=${PGID:-1000}

# Update the 'refrain' user/group to match the requested IDs
groupmod -o -g "$PGID" refrain
usermod -o -u "$PUID" refrain

# Change ownership of volume directories if they are mounted
chown -R refrain:refrain /data /config

exec su-exec refrain "$@"