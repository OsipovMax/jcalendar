FROM golang:1.19.1

ARG APP_PATH
WORKDIR $APP_PATH

# copy files for build
COPY . .

# env
ENV DOCKER_ENV=yes
ENV APP_PATH=$APP_PATH

CMD ["make", "run"]
