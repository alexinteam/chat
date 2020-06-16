START DB
============

``docker-compose up -d``
 
RUN MIGRATIONS
============

``docker-compose exec database psql -U postgres -f /var/www/migrations/0001_init.sql chat``
 
BULD APP
============
go build -o chatserv

RUN APP
===============

`run`   ./chatserv

`navigate`   http://localhost:8080/

`USERS LIST`
 
 |  USERNAME | PASSWORD  |
 |-------|--------|
 |        |        |
 | admin  | admin  |
 | moder  | moder  |
 | user1  | user1  |
 | user2  | user2  |