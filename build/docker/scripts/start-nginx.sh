#!/bin/sh

# Function to handle SIGTERM and terminate the script
terminate() {
	echo "Received SIGTERM, terminating..."
	exit 0
}

# Trap SIGTERM to call the terminate function
trap terminate SIGTERM

# Wait for the app to become ready
until curl -s http://app:$SERVER_PORT/health; do
	echo "Waiting for app on port: $SERVER_PORT..."
	sleep 2
done

# Inject Environment Variables into Nginx Config
envsubst '$SERVER_PORT' </etc/nginx/nginx.conf.template >/etc/nginx/nginx.conf

# Start Nginx
exec nginx -g 'daemon off;'
