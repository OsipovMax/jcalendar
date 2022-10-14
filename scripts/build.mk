build-app:
build-app:
	go build -trimpath ${LDFLAGS} -o bin/${APP_BINARY} cmd/${APP_NAME}/*.go