# usage: update.sh <id> <task>

curl --header "Content-Type: application/json" \
  --request PUT \
  --data '{"task":"'$2'"}' \
  http://localhost:8081/todos/$1
echo