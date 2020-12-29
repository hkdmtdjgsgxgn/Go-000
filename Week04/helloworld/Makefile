build:
	cd cmd/greeter_server && go build -o ../../bin/greeter_server && cd ../../ \
		&& cd cmd/greeter_client && go build -o ../../bin/greeter_client && cd ../../

run:
	./bin/greeter_server

test1:
	./bin/greeter_client
test2:
	./bin/greeter_client gopher
