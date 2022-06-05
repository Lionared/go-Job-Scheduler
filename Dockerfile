FROM alpine:3.15

WORKDIR /dist

COPY dist/goscheduler-docker  .

EXPOSE 20001

CMD ["/dist/goscheduler-docker", "-h", "0.0.0.0", "-p", "20001"]