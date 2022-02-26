chown -R mysql:mysql /var/lib/mysql
if ! [ -d "var/lib/mysql/avito" ]; then
service	mysql		start
mysql -u root -p -password=$MYSQL_ROOT_PASSWORD -e "CREATE DATABASE IF NOT EXISTS $MYSQL_AV_NAME DEFAULT CHARACTER SET utf8;"
mysql -u root -p -password=$MYSQL_ROOT_PASSWORD "avito" < init.sql
mysql -u root -p -password=$MYSQL_ROOT_PASSWORD -e "CREATE USER '$MYSQL_USER_NAME'@'%' IDENTIFIED BY '$MYSQL_ROOT_PASSWORD';"
mysql -u root -p -password=$MYSQL_ROOT_PASSWORD -e "GRANT ALL PRIVILEGES ON *.* TO '$MYSQL_USER_NAME'@'%';"
#mysql -u root -p -password=$MYSQL_ROOT_PASSWORD -e "FLUSH PRIVILEGES;"
mysqladmin -u root password $MYSQL_ROOT_PASSWORD
mysqladmin -u barcher password $MYSQL_ROOT_PASSWORD
service mysql stop
fi

mysqld_safe
