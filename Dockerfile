FROM golang:1.21.3-alpine

WORKDIR /spycat
COPY . . 
RUN go mod download
RUN go build -o spycat-ms

CMD [ "/spycat/spycat-ms" ]