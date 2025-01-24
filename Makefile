run:
	go run ./...

check:
	curl http://localhost:1323 	


createUser:
	curl -i -X POST 'http://localhost:1323/users' \
	-H 'Content-Type: application/json' \
	-H 'x-api-key: enjoy' \
	-d '{"name": "Dan","email": "dan@text.com"}'

getUser:
	curl -i http://localhost:1323/users/1

updateUser:
	curl -i -X PUT 'http://localhost:1323/users/1' \
	-H 'Content-Type: application/json' \
	-d '{"name": "ko","email": "ko@text.com"}'

updateUser/fail:
	curl -i -X PUT 'http://localhost:1323/users/9999' \
	-H 'Content-Type: application/json' \
	-d '{"name": "ko","email": "ko@text.com"}'


deleteUser:
	curl -i -X DELETE http://localhost:1323/users/1		