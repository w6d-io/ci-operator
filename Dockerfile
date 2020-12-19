# Build the ci-operator binary
ARG GOVERSION=1.14.6
FROM golang:$GOVERSION as builder
ARG GOVERSION=1.14.6
ARG VCS_REF
ARG BUILD_DATE
ARG VERSION
ENV GO111MODULE="on" \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY util/ util/

# Build
RUN  go build \
     -ldflags="-X 'main.Version=${VERSION}' -X 'main.Revision=${VCS_REF}' -X 'main.GoVersion=go${GOVERSION}' -X 'main.Built=${BUILD_DATE}' -X 'main.OsArch=${GOOS}/${GOARCH}'" \
     -mod=vendor \
     -a -o ci-operator main.go

# Use distroless as minimal base image to package the ci-operator binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
ARG VCS_REF
ARG BUILD_DATE
ARG VERSION
ARG PROJECT_URL
ARG USER_EMAIL="david.alexandre@w6d.io"
ARG USER_NAME="David ALEXANDRE"
LABEL maintainer="${USER_NAME} <${USER_EMAIL}>" \
        io.w6d.ci.vcs-ref=$VCS_REF       \
        io.w6d.ci.vcs-url=$PROJECT_URL   \
        io.w6d.ci.build-date=$BUILD_DATE \
        io.w6d.ci.version=$VERSION
WORKDIR /
COPY --from=builder /workspace/ci-operator .
RUN chown 1001:1001 /usr/local/bin/ci-operator
USER 1001:1001

ENTRYPOINT ["/usr/local/bin/ci-operator"]
