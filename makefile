gen:
	protoc proto/*.proto --go_out=. --go-grpc_out=.
clear:
	rm -r proto/*.go
