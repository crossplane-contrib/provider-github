FROM BASEIMAGE
RUN apk --no-cache add ca-certificates bash

ARG ARCH
ARG TINI_VERSION

ADD provider /usr/local/bin/crossplane-github-provider

EXPOSE 8080
ENTRYPOINT ["crossplane-github-provider"]
