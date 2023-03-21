cert:
	cd cert; sh ./generator.sh; cd ..

gen:
	protoc --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative ./proto/*.proto

clean:
	rm -f pb/proto/*.go

.PHONY: gen clean cert
