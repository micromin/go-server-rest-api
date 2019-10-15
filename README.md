#A Simple Go Restful server with Auth
Go simple (gos) is a simple GO server for adding todos. 
You can register and login and start adding your tasks.

##How to Run
First run a docker container for mysql by running

`docker run --name mysql57 -p 3306:3306 -e MYSQL_ROOT_PASSWORD=1234 -d mysql/mysql-server:5.7`

Then, enter to the mysql instance by running:

`docker exec -it mysql57 mysql -uroot -p`, after running enter `1234` as password. 

In order to create the required user, database and tables run:

```
CREATE USER 'gos' IDENTIFIED BY '1234';
grant all on *.* to 'gos'@'%' identified by '1234';
FLUSH PRIVILEGES;
```
```
CREATE DATABASE gos CHARACTER SET utf8 COLLATE utf8_general_ci;
use gos;
```
```
CREATE TABLE GOS_USER (
user_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
name VARCHAR(100) NOT NULL,
email VARCHAR(200),
password VARCHAR(100),
last_login int(10),
failed_login_attempt int(8),
date_created int(10),
date_updated int(10)) ENGINE=InnoDB;
```
```
CREATE TABLE GOS_TASK (
task_id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
user_id BIGINT UNSIGNED,
title VARCHAR(200) NOT NULL,
description TEXT,
date_created int(10),
date_updated int(10),
due_date int(10),
date_complete int(10),
FOREIGN KEY (user_id) REFERENCES GOS_USER(user_id)) ENGINE=InnoDB;
```
And finally run `go run main.go`.

The server should be running on port 8080.