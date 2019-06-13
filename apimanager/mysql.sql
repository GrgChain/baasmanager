
create database baas_api;

use baas_api;

-- auto-generated definition
create table user
(
  id       int auto_increment primary key,
  account  varchar(30)  not null,
  password varchar(100) not null,
  avatar   varchar(200) null,
  name     varchar(20)  not null,
  created      bigint not null,
  constraint user_account_uindex
  unique (account)
)  ENGINE=InnoDB  DEFAULT CHARSET=utf8 comment '用户表';

-- auto-generated definition
create table role
(
  rkey        varchar(20)  not null primary key,
  name        varchar(40)  not null,
  description varchar(200) null
)ENGINE=InnoDB  DEFAULT CHARSET=utf8 comment '角色表';

-- auto-generated definition
create table user_role
(
  user_id  int         not null,
  role_key varchar(20) not null,
  primary key (user_id, role_key)
) ENGINE=InnoDB  DEFAULT CHARSET=utf8 comment '用户角色表';

-- auto-generated definition
create table chain
(
	id           int auto_increment primary key,
	name         varchar(64) not null,
	user_account varchar(100) not null,
	description  varchar(255) null,
	consensus    varchar(10) not null,
	peers_orgs   varchar(100) not null,
	order_count  int not null,
	peer_count   int not null,
	tls_enabled  varchar(5) not null,
	status       int default '0' null,
	created      bigint not null
) ENGINE=InnoDB DEFAULT CHARSET=utf8 comment 'chain表';

-- auto-generated definition
create table channel
(
  id           int auto_increment primary key,
  chain_id     int          not null,
  orgs         varchar(255) not null,
  channel_name varchar(64)  not null,
  user_account varchar(100) not null,
  created      bigint not null
) ENGINE=InnoDB DEFAULT CHARSET=utf8 comment 'channel表';


-- auto-generated definition
create table chaincode
(
  id             int auto_increment primary key,
  chaincode_name varchar(64)  not null,
  channel_id     int          not null,
  user_account   varchar(100) not null,
  created        bigint       not null,
  version        varchar(10)  null,
  status         int default '0' null,
  github_path    varchar(256) null,
  args           varchar(500)    not null,
  policy         varchar(200)    not null
)ENGINE=InnoDB DEFAULT CHARSET=utf8 comment 'chaincode表';


-- admin 123456
INSERT INTO baas_api.user (id, account, password, avatar, name, created) VALUES (1, 'admin', 'pbkdf2_sha256$180000$JEavgdkTBzU3$3pIgoygm1QBtgbfEHeWZ7H4O2rEIgkgLxYV48mE+J4M=', 'https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif', 'admin', 1557977199);

INSERT INTO baas_api.role (rkey, name, description) VALUES ('admin', '管理员', '超级管理员,拥有所有权限');
INSERT INTO baas_api.role (rkey, name, description) VALUES ('user', '用户', '普通用户');

INSERT INTO baas_api.user_role (user_id, role_key) VALUES (1, 'admin');