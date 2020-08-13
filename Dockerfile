FROM golang:latest

ADD main .

EXPOSE 8000

CMD ./main