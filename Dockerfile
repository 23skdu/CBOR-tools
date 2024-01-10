FROM alpine:edge
LABEL author=23skdu@users.noreply.github.com
RUN set -x \
    && apk update && apk upgrade && apk add --no-cache \
    build-base perl-utils perl-canary-stability libcbor libcbor-dev \
    && cpan install CBOR::Free && rm -rf /root/.cpan && apk del build-base libcbor-dev 
COPY tool /
CMD ["/usr/bin/perl","-e", "print 'works';"] 
