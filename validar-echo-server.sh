#!/bin/bash
#echo "Hosts disponibles:"
#
#cat /etc/hosts

echo "Validacion del echo server con netcat"

PORT=$(grep '^SERVER_PORT' "./server/config.ini" | awk -F'=' '{print $2}' | tr -d ' ')
IP=$(grep '^SERVER_IP' "./server/config.ini" | awk -F'=' '{print $2}' | tr -d ' ')
MESSAGE="Hola mundo"
NETWORK_NAME=$(docker network ls --format '{{.Name}}' | grep testing_net)

echo "Server: $IP:$PORT"
echo "La red del server es: $NETWORK_NAME"

echo "Construyendo contenedor para ejecutar netcat"

docker build -f ./netcat/Dockerfile -t netcat .

echo "Enviando mensaje <$MESSAGE>"
RESPONSE=$(docker run --rm --network="$NETWORK_NAME" netcat sh -c "echo 'Hola mundo' | nc -w 2 '$IP' '$PORT'")
echo "Respuesta del servidor: $RESPONSE"

if [ "$RESPONSE" = "$MESSAGE" ]; then
  echo "action: test_echo_server | result: success"
else
  echo "action: test_echo_server | result: fail"
fi