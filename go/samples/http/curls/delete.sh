# usage: delete.sh <id>

curl --header "Content-Type: application/json" \
  --request DELETE \
  http://localhost:8080/todos/$1
echo