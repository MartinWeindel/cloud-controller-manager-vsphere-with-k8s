#####################      builder       #####################
FROM golang:1.16 AS builder

WORKDIR /go/src/github.com/martinweindel/cloud-controller-manager-vsphere-with-k8s
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
  -o /go/bin/cloud-controller-manager-vsphere-with-k8s \
  -ldflags="-s -w \
    -X github.com/martinweindel/cloud-controller-manager-vsphere-with-k8s/pkg/version.gitVersion=$(cat VERSION) \
    -X github.com/martinweindel/cloud-controller-manager-vsphere-with-k8s/pkg/version.gitCommit=$(git rev-parse --verify HEAD) \
    -X github.com/martinweindel/cloud-controller-manager-vsphere-with-k8s/pkg/version.buildDate=$(date --rfc-3339=seconds | sed 's/ /T/')" \
  cmd/main.go

#############      cloud-controller-manager-vsphere-with-k8s     #############
FROM alpine:3.12 AS cloud-controller-manager

COPY --from=builder /go/bin/cloud-controller-manager-vsphere-with-k8s /cloud-controller-manager-vsphere-with-k8s

WORKDIR /

ENTRYPOINT ["/cloud-controller-manager-vsphere-with-k8s"]
