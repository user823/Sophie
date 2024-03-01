use sophie;
drop table if exists t_user;
create table t_user(
                       u_id int not null auto_increment,
                       username varchar(10) collate utf8mb4_general_ci default null,
                       d_id int not null,
                       d_name varchar(10),
                       primary key(u_id),
                       unique key(username)
) engine = InnoDB default charset = utf8mb4 collate = utf8mb4_general_ci;

drop table if exists t_role;
create table t_role(
                       r_id int not null auto_increment,
                       rolename varchar(10) collate utf8mb4_general_ci default null,
                       primary key(r_id)
) engine = InnoDB default charset = utf8mb4 collate = utf8mb4_general_ci;


drop table if exists t_ur;
create table t_ur(
                     u_id int,
                     r_id int,
                     primary key(u_id, r_id)
) engine = InnoDB default charset = utf8mb4 collate = utf8mb4_general_ci;

drop table if exists t_dept;
create table t_dept(
                       d_id int not null auto_increment,
                       deptname varchar(10) collate utf8mb4_general_ci default null,
                       primary key(d_id)
) engine = InnoDB default charset = utf8mb4 collate = utf8mb4_general_ci;