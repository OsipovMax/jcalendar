APP_PROJECT=${APP_DIR}_${MODE}
TMP_DIR_NAME=.tmp

DOCKER_COMPOSE_CMD_DEVEL_DATA_PATH=${PWD}/${TMP_DIR_NAME}
DOCKER_COMPOSE_BUILD_ARGS=\
	GOPATH="${GOPATH}" \
	APP_NAME="${APP_NAME}" \
	APP_PATH="${APP_PATH}" \
	APP_BINARY="${APP_BINARY}" \
	APP_PROJECT="${APP_PROJECT}" \
	DATA_PATH="${DOCKER_COMPOSE_CMD_DEVEL_DATA_PATH}" \
	MODE="${MODE}"

DOCKER_COMPOSE_DIR=${PWD}/build/docker
DOCKER_COMPOSE_ENV=${DOCKER_COMPOSE_DIR}/.env
DOCKER_COMPOSE_ENV_TEST=${DOCKER_COMPOSE_DIR}/.env.test

DOCKER_COMPOSE_CMD=\
	${DOCKER_COMPOSE_BUILD_ARGS} docker-compose \
		-p ${APP_PROJECT} \
		--env-file ${DOCKER_COMPOSE_ENV} \
		-f ${DOCKER_COMPOSE_DIR}/docker-compose.yaml \
		-f ${DOCKER_COMPOSE_DIR}/docker-compose.${MODE}.yaml

DOCKER_COMPOSE_CMD_TEST=\
	${DOCKER_COMPOSE_BUILD_ARGS} docker-compose \
		-p ${APP_PROJECT} \
		--env-file ${DOCKER_COMPOSE_ENV_TEST} \
		-f ${DOCKER_COMPOSE_DIR}/docker-compose.yaml \
		-f ${DOCKER_COMPOSE_DIR}/docker-compose.${MODE}.yaml

DATA_DIRS=\
	./${TMP_DIR_NAME}/kafka \
	./${TMP_DIR_NAME}/postgres \
	./${TMP_DIR_NAME}/redis \
	./${TMP_DIR_NAME}/redis-dump \

docker-test-up:	## Up compose test docker images
docker-test-up: .mode-test .copy-default-env .copy-default-env-test
	${DOCKER_COMPOSE_CMD_TEST} up

docker-test-build:	## Build compose test docker images
docker-test-build: .mode-test .copy-default-env .copy-default-env-test
	${DOCKER_COMPOSE_CMD_TEST} build

docker-test:	## Run compose test
docker-test: cmd?=test
docker-test: .mode-test
	${DOCKER_COMPOSE_CMD_TEST} exec app make ${cmd}

docker-test-down:	## Stop docker-compose to test and remove db
docker-test-down: .mode-test
	${DOCKER_COMPOSE_CMD_TEST} down -v

docker-dev-up:	## Up compose development docker images
docker-dev-up: .mode-dev .copy-default-env .dev-data-dirs
	${DOCKER_COMPOSE_CMD} up

docker-dev-build:	## Build compose development docker images
docker-dev-build: .mode-dev .copy-default-env .dev-data-dirs
	${DOCKER_COMPOSE_CMD} build

docker-dev-down:	## Stop development mode docker-compose
docker-dev-down: .mode-dev
	${DOCKER_COMPOSE_CMD} down -v

docker-dev-reset-force:
	$(MAKE) docker-dev-down
	rm -rf ${DATA_DIRS}

DEV_OPEN_PORTS=\
	-p 3009:3009 \
	-p 8080:8080 \
	-p 8081:8081 \
	-p 8048:8048 \
	-p 6060:6060 \
	-p 9090:9090 \
	-p 2345:2345

docker-dev:	## Run compose in development mode with live reload
docker-dev: cmd?=run-live-reload
docker-dev: .mode-dev .copy-default-env .dev-data-dirs
	${DOCKER_COMPOSE_CMD} run --rm ${DEV_OPEN_PORTS} --name ${MODE}_run_app app make ${cmd}

.copy-default-env:
	@if [ ! -e ${DOCKER_COMPOSE_ENV} ]; then \
  		cp ${DOCKER_COMPOSE_DIR}/.env.sample ${DOCKER_COMPOSE_ENV}; \
  		echo "Copy default .env file"; \
  	fi

.copy-default-env-test:
	@if [ ! -e ${DOCKER_COMPOSE_ENV_TEST} ]; then \
		cp ${DOCKER_COMPOSE_DIR}/.env.sample ${DOCKER_COMPOSE_ENV_TEST}; \
		echo "Copy default .env.test file"; \
	fi

.mode-test:
	$(eval MODE=test)

.mode-dev:
	$(eval MODE=dev)

.dev-data-dirs:
	@mkdir -p ${DATA_DIRS}
