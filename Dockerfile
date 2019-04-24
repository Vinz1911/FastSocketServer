# get golang container -> latest build
FROM golang:latest

# create workdir path
WORKDIR /go/src/FastSocketServer
# copy all bundle resources to the workdir path
COPY . .

# call go get to fetch all frameworks/librarys
# install those librarys
# compile golang executable
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build FastSocketServer.go

# expose port
EXPOSE 3333

# run the application
CMD ["FastSocketServer"]