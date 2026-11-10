FROM busybox
ADD app/index.html /www/index.html
EXPOSE 8005
CMD httpd -p 8005 -h /www; tail -f /dev/null
