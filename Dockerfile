FROM alpine:edge
LABEL author=23skdu@users.noreply.github.com
RUN set -x \
    && apk update && apk upgrade && apk add --no-cache \
    perl-utils perl-lwp-protocol-https 
COPY tool /bin
CMD ["/usr/bin/perl","-e", "print 'works';"] 
