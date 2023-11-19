This is readme.txt file

Steps to execute the given source code
=============================================
1) To build the source code on the server
go build LogIngestor_QueryInterface.go

2) To run the executable server
./LogIngestor_QueryInterface

3) Trigger the request using curl or postmain.
Step to test using curl
curl -X POST -H "Content-Type: application/json" -d '{  "level": "error" }' http://localhost:3000/query