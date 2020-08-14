FROM golang:latest

WORKDIR /go/src/ImageStoreService

COPY . .

 

RUN go get -d -v ./...

RUN go install -v ./...

 
ADD main .

EXPOSE 8000

CMD ./main