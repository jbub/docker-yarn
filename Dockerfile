FROM node:7.5-alpine
MAINTAINER Juraj Bubniak <contact@jbub.eu>

ENV PATH /root/.yarn/bin:$PATH

RUN apk --no-cache add curl bash binutils tar \
  && touch /root/.bashrc \
  && curl -o- -L https://yarnpkg.com/install.sh | bash \
  && apk del curl tar binutils

ENTRYPOINT [ "yarn" ]