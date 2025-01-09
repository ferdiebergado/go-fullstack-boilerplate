#!/bin/sh

# Inject Environment Variables into Nginx Config
envsubst '$SERVER_PORT' </etc/nginx/nginx.conf.template >/etc/nginx/nginx.conf

# Start Nginx
exec nginx -g 'daemon off;'
