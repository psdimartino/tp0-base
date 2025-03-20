echo "Validacion del echo server con netcat"

PORT=$(grep '^SERVER_PORT' "./server/config.ini" | awk -F'=' '{print $2}' | tr -d ' ')
IP=$(grep '^SERVER_IP' "./server/config.ini" | awk -F'=' '{print $2}' | tr -d ' ')

echo "Server hostname: $IP. Resolviendo IP,"

#if ! [[ "$IP" =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
#    IP=$(getent hosts "$IP" | awk '{print $1}')
#fi

echo "Server: $IP:$PORT"

echo "Enviando mensaje <Hola mundo>"



RESPONSE=$(echo "$MESSAGE" | nc -w 2 "$IP" "$PORT")

if [[ "$RESPONSE" == "$MESSAGE" ]]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi