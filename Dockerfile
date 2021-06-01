# Build the ci-operator binary
ARG GOVERSION=1.15.5
FROM golang:$GOVERSION as builder
ARG GOVERSION=1.15.5
ARG VCS_REF
ARG BUILD_DATE
ARG VERSION
ENV GO111MODULE="on" \
    CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
COPY vendor/ vendor
# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY internal/ internal/
COPY pkg/ pkg/

# Build
RUN  go build \
     -ldflags="-X 'main.Version=${VERSION}' -X 'main.Revision=${VCS_REF}' -X 'main.GoVersion=go${GOVERSION}' -X 'main.Built=${BUILD_DATE}' -X 'main.OsArch=${GOOS}/${GOARCH}'" \
     -mod=vendor \
     -a -o ci-operator main.go
RUN chown 1001:1001 ci-operator

# Use distroless as minimal base image to package the ci-operator binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/base:nonroot
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
USER 1001:1001

ENTRYPOINT ["/ci-operator"]
