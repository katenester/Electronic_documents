FROM golang:latest
LABEL authors="katy"

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh
# install go-migrate
#RUN GO111MODULE=on go get -v github.com/golang-migrate/migrate/v4/cmd/migrate@latest
# build go app
RUN go mod download
RUN go build -o doc-app ./cmd/main.go

CMD ["./doc-app"]
#CMD migrate -path ./schema -database 'postgres://postgres:Katy314@db/postgres?sslmode=disable' up && ./todo-app