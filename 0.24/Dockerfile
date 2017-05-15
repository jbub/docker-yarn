FROM node:7.9-alpine
MAINTAINER Juraj Bubniak <juraj.bubniak@gmail.com>

ENV YARN_VERSION=0.24.4
ENV PATH /root/.yarn/bin:$PATH

RUN apk --no-cache add gnupg curl bash binutils tar \
  && touch /root/.bashrc \
  && curl -o- -L https://yarnpkg.com/install.sh | bash -s -- --version ${YARN_VERSION} \
  && apk del gnupg curl tar binutils

ENTRYPOINT [ "yarn" ]