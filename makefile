test:
	go test ./vault/... -v

test-cov:
	if ! [ -d "coverage" ]; then \
		echo "Creating covorage folder" ; \
		mkdir coverage; \
	fi
	go test ./vault/... -cover -v -coverprofile ./coverage/c.out && go tool cover -html=./coverage/c.out -o ./coverage/coverage.html