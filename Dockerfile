FROM --platform=$BUILDPLATFORM golang:1-alpine AS build
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
RUN --mount=target=/mnt ["/mnt/build.sh"]

FROM scratch
COPY --from=build /hfs /
VOLUME ["/htdocs"]
ENTRYPOINT ["/hfs"]
CMD ["-addr", ":8000", "-dotfile", "-root", "/htdocs"]
