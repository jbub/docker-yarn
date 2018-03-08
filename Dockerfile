
FROM node:9.7.1-alpine
MAINTAINER Juraj Bubniak <juraj.bubniak@gmail.com>

ENV YARN_VERSION=1.5.1
ENV PATH /root/.yarn/bin:$PATH

RUN apk --no-cache add gnupg curl bash binutils tar \
  && touch /root/.bashrc \
  && curl -o- -L https://yarnpkg.com/install.sh | bash -s -- --version ${YARN_VERSION} \
  && apk del gnupg curl tar binutils

ENTRYPOINT [ "yarn" ]