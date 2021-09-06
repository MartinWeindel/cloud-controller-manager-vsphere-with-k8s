#####################      builder       #####################
FROM eu.gcr.io/gardener-project/3rd/golang:1.16.7 AS builder

WORKDIR /build
COPY . .

RUN make build

#############      cloud-controller-manager-vsphere-with-k8s     #############
FROM eu.gcr.io/gardener-project/3rd/alpine:3.13.5 AS cloud-controller-manager

COPY --from=builder /build/bin/cloud-controller-manager-vsphere-with-k8s /bin/vspherewk8s-cloud-controller-manager

WORKDIR /

ENTRYPOINT ["/bin/vspherewk8s-cloud-controller-manager"]
