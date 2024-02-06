CURRENT_DIR=$(shell pwd)

APP=$(shell basename ${CURRENT_DIR})
APP_CMD_DIR=${CURRENT_DIR}/cmd

build:
	CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o ${CURRENT_DIR}/bin/${APP} ${CURRENT_DIR}/$(main)


# make migration name=users_ddl
migration:
	migrate create -ext sql -dir ${CURRENT_DIR}/migrations -seq -digits 2 $(name)