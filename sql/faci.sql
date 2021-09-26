drop database faci;
create database faci;
use faci;

create table if not exists hour_data (
	id int auto_increment,
	create_time datetime unique not null comment "时间",
	capital_b decimal(20, 4) not null comment "锁仓量B",
	lowcase_b decimal(20, 4) not null comment "可流通量b",
	cfil_to_fil decimal(6, 4) not null comment "cfil to fil",
	count_drawns_fil decimal(20, 4) not null comment "累计已提取FIL",
	primary key(id)
);

create table if not exists 5_mins_data (
	id int auto_increment,
	primary key(id)
);
