#Dockerfile
FROM ubuntu:14.04
MAINTAINER David Malone <avidmalone@gmail.com>

ENV AWS_ENDPOINT awsendpoint
ENV AWS_ACCESS_KEY_ID awsaccesskeyid
ENV AWS_SECRET_ACCESS_KEY awssecretaccesskey
ENV EMAIL_TEMPLATE_DIR ./
ENV FB_APP_ID fbappid
ENV FB_APP_SECRET fbappsecret

ADD ./build/linux64/cms ./app

EXPOSE 3000

ENTRYPOINT ./app
