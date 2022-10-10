APP_NAME?=app
APP_PATH?=/go/src/$(shell grep ^module go.mod | awk '{print $$2}')
APP_DIR?=$(notdir ${APP_PATH})
APP_BINARY?=${APP_DIR}-${APP_NAME}