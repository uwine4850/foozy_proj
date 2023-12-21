#!/bin/bash

sudo docker exec -i mysql mysql --defaults-extra-file=/schema/mysql.cnf -e "CREATE DATABASE IF NOT EXISTS foozy_proj;"
sudo docker exec -i mysql mysql --defaults-extra-file=./schema/mysql.cnf < ./schema/profile.sql
sudo docker exec -i mysql mysql --defaults-extra-file=./schema/mysql.cnf < ./schema/chat.sql
sudo docker exec -i mysql mysql --defaults-extra-file=./schema/mysql.cnf < ./schema/chat_msg.sql
sudo docker exec -i mysql mysql --defaults-extra-file=./schema/mysql.cnf < ./schema/chat_msg_count.sql

#sudo docker compose run --rm node npm install
