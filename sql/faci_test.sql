#drop database faci_test;
create database faci_test;
use faci_test;

create table if not exists hour_data (
	id int auto_increment,
	create_time datetime unique not null comment "时间戳",
	lowcase_b decimal(36, 18) not null comment "可流通量b",
	count_drawns_fil decimal(36, 18) not null comment "累计已提取FIL",
	primary key(id)
);

create table if not exists 5_mins_data (
	id int auto_increment,
	create_time datetime unique not null comment "时间",
	cfil_to_fil decimal(6, 4) not null comment "cfil to fil",
	primary key(id)
);

create table if not exists fil_node (
	id int auto_increment,
	node_name varchar(16) not null comment "节点名称",
	address varchar(128) not null comment "owner 地址",
	balance decimal(36, 18) not null comment "账户总余额",
	worker_balance decimal(36, 18) not null comment "worker余额",
	quality_adj_power decimal(36, 18) not null comment "有效算力",
	available_balance decimal(36, 18) not null comment "可用余额",
	pledge decimal(36, 18) not null comment "扇区抵押",
	vestingFunds decimal(36, 18) not null comment "存储服务锁仓",
	singletT decimal(36, 18) not null comment "单T",
	hour_data_id int not null,
	primary key(id),
	foreign key(hour_data_id) references hour_data(id) on update cascade on delete cascade
);

drop procedure if exists insertInitData;

delimiter $$
create procedure insertInitData()
begin
	declare i int default 0;
	declare second int default 3600;
	declare current_second bigint default current_timestamp;
	while i < 24 do
		insert into hour_data values(null, from_unixtime(unix_timestamp(current_second) div second * second - second * i), 0, 0);
		set i = i + 1;
	end while;

	set i = 0;
	set second = 300;
	while i < 288 do
		insert into 5_mins_data values(null, from_unixtime(unix_timestamp(current_second) div second * second - second * i), 0);
		set i = i + 1;
	end while;
end $$

delimiter ;
call insertInitData();
drop procedure insertInitData;
