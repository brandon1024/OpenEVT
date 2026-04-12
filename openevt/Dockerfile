ARG BUILD_FROM=ghcr.io/home-assistant/base:latest
ARG BUILD_VERSION
FROM brandon1024/openevt:$BUILD_VERSION AS builder

FROM $BUILD_FROM

COPY --from=builder /usr/bin/openevt /usr/bin/openevt

COPY run.sh /
RUN chmod a+x /run.sh

CMD [ "/run.sh" ]
