create user otabek;
alter user otabek with password 'otabek123';
create database test;
alter database test owner to otabek;
create table users(
    username varchar(255) unique not null,
    password varchar(255) not null,
    firstname varchar(50) not null,
    lastname varchar(50) not null
);