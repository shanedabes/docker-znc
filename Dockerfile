FROM golang:alpine as build-env

ADD . /src
RUN cd /src && go build -o zncconfer

FROM linuxserver/znc:znc-1.7.5-ls28
COPY --from=build-env /src/zncconfer /bin/
COPY run /etc/services.d/znc/
