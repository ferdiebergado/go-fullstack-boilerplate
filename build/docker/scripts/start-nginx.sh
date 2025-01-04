#!/bin/sh

# Function to handle os signals and terminate the script
terminate() {
	echo "Received stop signal, terminating..."
	exit 0
}

# Trap os signals to call the terminate function
trap terminate SIGTERM SIGINT SIGQUIT

# Wait for the app to become ready
until curl -s http://app:$SERVER_PORT/api/health; do
	echo "Waiting for app on port: $SERVER_PORT..."
	sleep 2
done

# Inject Environment Variables into Nginx Config
envsubst '$SERVER_PORT' </etc/nginx/nginx.conf.template >/etc/nginx/nginx.conf

# Start Nginx
exec nginx -g 'daemon off;'
