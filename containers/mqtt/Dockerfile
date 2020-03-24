FROM eclipse-mosquitto:latest

COPY ./mosquitto/config/mosquitto.conf ./mosquitto/config

EXPOSE 1883

ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["/usr/sbin/mosquitto", "-c", "/mosquitto/config/mosquitto.conf"]