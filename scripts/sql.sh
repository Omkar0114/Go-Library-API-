// create database

create database if not exists library;

use library;

create table if not exists book (
    id int not null auto_increment,
    name varchar(255) not null,
    isbn varchar(255) not null,
    primary key (id)
);

insert into book (name, isbn) values ('The Kubernetes store', 'isbn-1234');