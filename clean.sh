docker rm $(docker ps -a -q) 2>/dev/null
docker rmi $(docker images -q) 2>/dev/null
docker volume rm $(docker volume ls -q) 2>/dev/null
docker system prune -a --volumes
sudo rm -rf /home/barcher/data/go_volume
sudo rm -rf /home/barcher/data/mariadb_volume
