FROM --platform=$BUILDPLATFORM golang:1-alpine AS build
ADD --chmod=755 https://github.com/zhangyoufu/actions/raw/refs/heads/main/go-crossbuild.sh /go-crossbuild.sh
ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT
RUN --mount=target=/mnt CGO_ENABLED=0 /go-crossbuild.sh -C /mnt -o /hfs ./cmd/hfs

FROM scratch
COPY --from=build /hfs /
VOLUME ["/htdocs"]
ENTRYPOINT ["/hfs"]
CMD ["-addr", ":8000", "-dotfile", "-root", "/htdocs"]
