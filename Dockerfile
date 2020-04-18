FROM alpine
MAINTAINER Julian Kleinhans <julian.kleinhans@aoe.com>

RUN apk --no-cache add ca-certificates && update-ca-certificates
RUN mkdir -p /usr/bin
ADD http-monitoring /usr/bin/
RUN chmod +x /usr/bin/http-monitoring
ENV PATH $PATH:/usr/bin

EXPOSE 1877

CMD ["/usr/bin/http-monitoring"]