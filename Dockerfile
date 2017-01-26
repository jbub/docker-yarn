FROM node:7.4-alpine
MAINTAINER Juraj Bubniak <contact@jbub.eu>

ENV PATH /root/.yarn/bin:$PATH

RUN apk --no-cache add curl bash binutils tar \
  && rm -rf /var/cache/apk/* \
  && /bin/bash \
  && touch ~/.bashrc \
  && curl -o- -L https://yarnpkg.com/install.sh | bash \
  && apk del curl tar binutils

ENTRYPOINT [ "yarn" ]