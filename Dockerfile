############################
# STEP 1 build executable binary
############################
FROM golang@sha256:e9f6373299678506eaa6e632d5a8d7978209c430aa96c785e5edcb1eebf4885e AS builder
# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates
# Create appuser.
ENV USER=appuser
ENV UID=10001 
# See https://stackoverflow.com/a/55757473/12429735RUN 
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"
#WORKDIR $GOPATH/src/github.com/admyavuzs/golang-hello-app
WORKDIR $GOPATH/src/github.com/admyavuzs/hello-world
COPY . .
# Fetch dependencies.
# Using go get.
RUN go get -d -v
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /server
############################
# STEP 2 build a small image
############################
FROM scratch
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
# Copy our static executable.
COPY --from=builder /server /server
USER appuser:appuser
WORKDIR /root/
# Run the hello binary.
EXPOSE 8081
ENTRYPOINT [ "/server" ]