FROM debian:trixie-slim
LABEL author=23skdu@users.noreply.github.com
RUN set -x \
    && apt update && apt upgrade -y && apt install -y build-essential \ 
    && cpan install CBOR::Free && rm -rf /root/.cpan && apt remove -y build-essential 
COPY tool /
CMD ["/usr/bin/perl","-e", "print 'works';"] 
