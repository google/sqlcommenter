# usage: update_with_traceparent.sh <id> <task>

curl --header "Content-Type: application/json" \
  --header "Traceparent: 00-a7586e0a5bc7934ce028e83bc1d247f2-6b94506168fc7803-01" \
  --request PUT \
  --data '{"task":"'$2'"}' \
  http://localhost:8081/todos/$1
echo