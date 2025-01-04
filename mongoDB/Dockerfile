FROM mongo:6.0

WORKDIR /data/db

COPY mongod.conf /etc/mongod.conf

COPY init.sh /docker-entrypoint-initdb.d/init.sh
COPY schema.json /docker-entrypoint-initdb.d/schema.json

RUN chmod +x /docker-entrypoint-initdb.d/init.sh

EXPOSE 27017

CMD ["mongod", "--config", "/etc/mongod.conf"]
