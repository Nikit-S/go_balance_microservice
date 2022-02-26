FROM debian:buster

SHELL ["/bin/bash", "-c"]

RUN	apt-get -y update
RUN	apt-get -y upgrade
RUN	apt-get -y install apt-utils
RUN	apt-get -y install vim
RUN apt-get -y install curl
RUN curl -O https://dl.google.com/go/go1.12.7.linux-amd64.tar.gz
RUN tar xvf go1.12.7.linux-amd64.tar.gz
RUN chown -R root:root ./go
RUN mv go /usr/local

COPY .profile /root/
ADD srcs ./srcs

RUN source /root/.profile
ENV PATH="/usr/local/go/bin:${PATH}"

EXPOSE 80

WORKDIR "./srcs/"
CMD go build -mod=vendor main.go; ./main
