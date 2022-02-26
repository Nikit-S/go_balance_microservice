FROM debian:buster

RUN	apt-get -y update
RUN	apt-get -y upgrade
RUN	apt-get -y install apt-utils
RUN	apt-get -y install vim


RUN		apt-get install -y mariadb-server
COPY	ignite.sh	.

COPY ./50-server.cnf /etc/mysql/mariadb.conf.d/
COPY ./init.sql .
RUN chmod 444 /etc/mysql/mariadb.conf.d/50-server.cnf

EXPOSE 3306

CMD ./ignite.sh
