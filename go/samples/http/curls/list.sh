# usage: list.sh

curl http://localhost:8081/todos?search=$1
echo