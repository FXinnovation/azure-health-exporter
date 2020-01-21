FROM golang:1.12 as builder
WORKDIR /azure-health-exporter/
COPY . .
RUN make getpromu test build

FROM ubuntu:18.04
COPY --from=builder /azure-health-exporter/azure-health-exporter /azure-health-exporter
ADD ./resources /resources
RUN /resources/build && rm -rf /resources
USER ahe
EXPOSE 9613
WORKDIR /opt/azure-health-exporter
ENTRYPOINT  [ "/opt/azure-health-exporter/azure-health-exporter" ]

LABEL maintainer="FXinnovation CloudToolDevelopment <CloudToolDevelopment@fxinnovation.com>" \
      "org.label-schema.name"="azure-health-exporter" \
      "org.label-schema.base-image.name"="docker.io/library/ubuntu" \
      "org.label-schema.base-image.version"="18.04" \
      "org.label-schema.description"="azure-health-exporter in a container" \
      "org.label-schema.url"="https://github.com/FXinnovation/azure-health-exporter" \
      "org.label-schema.vcs-url"="https://github.com/FXinnovation/azure-health-exporter" \
      "org.label-schema.vendor"="FXinnovation" \
      "org.label-schema.schema-version"="1.0.0-rc.1" \
      "org.label-schema.usage"="Please see README.md"