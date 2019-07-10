# User golang based off al Alpine
FROM golang:latest

LABEL "vendor": "PavedRoad.io" \
      "microservice": "prUserIdMapper" \
      "description": "Maps 3rd party credentials to pr user UUID" \
      "version": "0.0.1" \
      "env": "dev"

MAINTAINER "support@pavedroad.io"

# Build paths for placing kevlar microservice
ENV ms prUserIdMapper
ENV pavedroad /pavedroad
ENV pavedroadbin $pavedroad/$ms

# make working directory, move to it, and copy in prToken ms
RUN mkdir ${pavedroad}
WORKDIR ${pavedroad}
COPY $ms $pavedroad

EXPOSE 8081
CMD ["/bin/sh", "-c", "$pavedroadbin"]
