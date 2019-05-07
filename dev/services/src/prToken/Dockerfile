# User golang based off al Alpine
FROM golang:latest

LABEL "vendor": "PavedRoad.io" \
      "microservice": "prToken" \
      "description": "Stores OAUTH access tokens" \
      "version": "0.0.1" \
      "env": "dev"

MAINTAINER "support@pavedroad.io"

# Build paths for placing kevlar microservice
ENV ms prToken
ENV kevlar /kevlar
ENV kevlarbin $kevlar/$ms


# make working directory, move to it, and copy in prToken ms
RUN mkdir ${kevlar}
WORKDIR ${kevlar}
Add $ms $kevlar

EXPOSE 8081
CMD ["/bin/sh", "-c", "$kevlarbin"]
