#!/bin/bash

MONGO_USER=${MONGO_INITDB_ROOT_USERNAME}
MONGO_PASS=${MONGO_INITDB_ROOT_PASSWORD}
MONGO_DB=${MONGO_INITDB_DATABASE:-ecommerce}
MONGO_HOST=localhost
MONGO_PORT=27017

until mongosh --host $MONGO_HOST --port $MONGO_PORT -u $MONGO_USER -p $MONGO_PASS --authenticationDatabase admin --eval "quit()" > /dev/null 2>&1; do
  echo "Esperando a MongoDB..."
  sleep 2
done

echo "MongoDB está listo. Creando la colección..."

mongosh --host $MONGO_HOST --port $MONGO_PORT -u $MONGO_USER -p $MONGO_PASS --authenticationDatabase admin <<EOF
use $MONGO_DB;

db.createCollection("products", {
  validator: { \$jsonSchema: $(cat /docker-entrypoint-initdb.d/schema.json) }
});

EOF

echo "Colección creada con éxito."
