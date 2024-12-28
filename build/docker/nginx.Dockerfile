FROM nginx:1.27.2-alpine3.20

COPY ./web/static/ /usr/share/nginx/html/
COPY ./build/docker/scripts/start-nginx.sh /start-nginx.sh
