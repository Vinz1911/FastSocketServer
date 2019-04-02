# get golang container -> latest build
FROM golang:latest

# create workdir path
WORKDIR /go/src/socket_server
# copy all bundle resources to the workdir path
COPY . .

# call go get to fetch all frameworks/librarys
# install those librarys
# compile golang executable
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build socket_server.go

# expose port
EXPOSE 3333

# run the application
CMD ["socket_server"]