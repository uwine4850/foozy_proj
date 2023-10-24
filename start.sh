#!/bin/bash

sudo docker exec -i mysql mysql --defaults-extra-file=/schema/mysql.cnf -e "CREATE DATABASE IF NOT EXISTS foozy_proj;"
sudo docker exec -i mysql mysql --defaults-extra-file=./schema/mysql.cnf < ./schema/profile.sql

#sudo docker compose run --rm node npm install
