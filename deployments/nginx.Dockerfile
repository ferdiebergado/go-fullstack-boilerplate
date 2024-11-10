FROM nginx:1.27.2-alpine3.20

COPY configs/nginx.conf /etc/nginx/nginx.conf
COPY web/static /usr/share/nginx/html
