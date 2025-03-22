#!/bin/sh

WG_CONTAINER="wireguard"

echo "Restarting WireGuard in docker_container $WG_CONTAINER..."

docker restart $WG_CONTAINER

if [ $? -eq 0 ]; then
    echo "WireGuard sussesfully restarted $WG_CONTAINER."
else
    echo "Error: Cannot restart WireGuard."
    exit 1
fi
