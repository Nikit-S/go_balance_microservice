start:
	sudo mkdir -p /home/barcher/data/go_volume
	sudo mkdir -p /home/barcher/data/mariadb_volume
	docker-compose -f ./srcs/docker-compose.yml up -d

stop:
	docker-compose -f ./srcs/docker-compose.yml down

clean:
	./clean.sh
