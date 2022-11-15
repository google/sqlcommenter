curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"task":"'$1'"}' \
  http://localhost:8080/todos
echo