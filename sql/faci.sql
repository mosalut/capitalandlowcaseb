#create database faci;
use faci;

create table if not exists hour_data (
	id int auto_increment,
	create_date date unique not null comment "日期",
	apy varchar(32) not null comment "APY",
	cfil_to_fil varchar(32) not null comment "cfil to fil",
	capital_b varchar(32) not null comment "锁仓量B",
	lowcase_b varchar(32) not null comment "可流通量b",
	loss varchar(32) not null comment "损耗值",
	primary key(id)
);

create table if not exists 5_mins_data (
	id int auto_increment,
	primary key(id)
);
