FROM golang:1.13 as modules
ADD go.mod /m/
RUN cd /m && go mod download && mkdir -p /go/pkg/

FROM golang:1.13 as builder
COPY --from=modules /go/pkg/ /go/pkg
RUN mkdir -p /opt/resource/ 
WORKDIR /opt/resource/
COPY . .
WORKDIR /opt/resource/cmd/umserver
RUN go build -o /opt/services/user-manager .

FROM alpine:3.7
COPY --from=builder /opt/services/user-manager /opt/services/user-manager
CMD /opt/services/user-manager