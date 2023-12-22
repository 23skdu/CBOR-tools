FROM tool:0.1
LABEL author=23skdu@users.noreply.github.com
RUN set -x \
    && mkdir /tests \
    && apk update && apk upgrade && apk add --no-cache \
    perl-critic perl-test-most \
    && cpan install TAP::Parser && rm -rf ~/.cpan 
ADD t/* /tests
ENTRYPOINT ["/bin/bash","-c","prove -v /tests/*"]] 
