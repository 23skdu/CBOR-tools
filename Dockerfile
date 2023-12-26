FROM alpine:edge
LABEL author=23skdu@users.noreply.github.com
RUN set -x \
    && apk update && apk upgrade && apk add --no-cache \
    perl-utils perl-canary-stability libcbor libcbor-dev\
    && cpan install CBOR::XS  && rm -rf /root/.cpan 
COPY tool /bin
CMD ["/usr/bin/perl","-e", "print 'works';"] 
