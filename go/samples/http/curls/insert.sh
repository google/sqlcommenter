# usage: insert.sh <task>

curl --header "Content-Type: application/json" \
  --request POST \
  --data '{"task":"'$1'"}' \
  http://localhost:8081/todos
echo