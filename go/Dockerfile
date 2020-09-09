############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

# Install git.
RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/github.com/aosh50/momenton/go
COPY . .

# Fetch dependencies. Using go get.
RUN go get -d -v
# Build the binary.
RUN go build -o /go/bin/momenton

############################
# STEP 2 build a small image
############################
FROM alpine
# Copy our static executable.
COPY --from=builder /go/bin/momenton /go/bin/momenton

# Run the binary.
ENTRYPOINT ["/go/bin/momenton"]