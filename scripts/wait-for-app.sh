#!/bin/sh

# Wait for the app to become ready
until curl -s http://app:$PORT/health; do
	echo "Waiting for app on port: $PORT..."
	sleep 2
done

# Start Nginx
exec nginx -g 'daemon off;'
