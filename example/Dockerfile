# get golang container -> latest build
FROM golang:latest

# create workdir path
WORKDIR /go/src/server
# copy all bundle resources to the workdir path
COPY . .

# call go get to fetch all frameworks/librarys
# install those librarys
# compile golang executable
RUN go get -d -v ./...
RUN go install -v ./...
RUN go build server.go

# expose port
EXPOSE 8081

# run the application
CMD ["server"]