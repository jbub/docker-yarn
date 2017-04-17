FROM node:7.9-alpine
MAINTAINER Juraj Bubniak <contact@jbub.eu>

ENV YARN_VERSION=0.21.3
ENV PATH /root/.yarn/bin:$PATH

RUN apk --no-cache add gnupg curl bash binutils tar \
  && touch /root/.bashrc \
  && curl -o- -L https://yarnpkg.com/install.sh | bash -s -- --version ${YARN_VERSION} \
  && apk del gnupg curl tar binutils

ENTRYPOINT [ "yarn" ]