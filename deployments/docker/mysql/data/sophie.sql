CREATE DATABASE IF NOT EXISTS `sophie`;
USE `sophie`;

-- ----------------------------
-- 1、部门表
-- ----------------------------
drop table if EXISTS sys_dept;

CREATE table sys_dept (
    id bigint(20) not null auto_increment comment '部门id',
    extend_shadow text comment '扩展字段',
    status char(1) default '0' comment '部门状态（0正常 1停用）',
    created_by varchar(64) default '' comment '创建者',
    created_at dateTime comment '创建时间',
    updated_by varchar(64) default '' comment '创建时间',
    updated_at dateTime comment '更新时间',
    parent_id bigint(20) default 0 comment '父部门id',
    ancestors varchar(50) default '' comment '祖级列表',
    dept_name varchar(30) default '' comment '部门名称',
    order_num int(4) default 0 comment '显示顺序',
    leader varchar(20) default null comment '负责人',
    phone varchar(11) default null comment '联系电话',
    email varchar(50) default null comment '电子邮件',
    del_flag char(1) default '0' comment '删除标志（0代表存在， 1代表正常）',
    primary key (id),
    CONSTRAINT `fk_dept_dept` FOREIGN key (`parent_id`) REFERENCES `sys_dept` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) engine = innodb auto_increment = 200 comment = '部门表';

-- ----------------------------
-- 2、用户信息表
-- ----------------------------
drop table if exists sys_user;

create table sys_user (
    id bigint(20) not null auto_increment comment '用户id',
    extend_shadow text comment '扩展字段',
    status char(1) default '0' comment '帐号状态（0正常 1停用）',
    created_by varchar(64) default '' comment '创建者',
    created_at dateTime comment '创建时间',
    updated_by varchar(64) default '' comment '创建时间',
    updated_at dateTime comment '更新时间',
    dept_id bigint(20) default null comment '部门ID',
    user_name varchar(30) not null comment '用户账号',
    nick_name varchar(30) not null comment '用户昵称',
    user_type varchar(2) default '00' comment '用户类型（00系统用户）',
    email varchar(50) default '' comment '用户邮箱',
    phonenumber varchar(11) default '' comment '手机号码',
    sex char(1) default '0' comment '用户性别（0男 1女 2未知）',
    avatar varchar(100) default '' comment '头像地址',
    password varchar(100) default '' comment '密码',
    del_flag char(1) default '0' comment '删除标志（0代表存在 2代表删除）',
    login_ip varchar(128) default '' comment '最后登录IP',
    login_date datetime comment '最后登录时间',
    remark varchar(500) default null comment '备注',
    primary key (id),
    CONSTRAINT `fk_user_dept` FOREIGN KEY (`dept_id`) REFERENCES `sys_dept` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) engine = innodb auto_increment = 100 comment = '用户信息表' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 3、岗位信息表
-- ----------------------------
drop table if exists sys_post;

create table sys_post (
    id bigint(20) not null auto_increment comment '岗位id',
    extend_shadow text comment '扩展字段',
    status char(1) default '0' comment '岗位状态（0正常 1停用）',
    created_by varchar(64) default '' comment '创建者',
    created_at dateTime comment '创建时间',
    updated_by varchar(64) default '' comment '创建时间',
    updated_at dateTime comment '更新时间',
    post_code varchar(64) not null comment '岗位编码',
    post_name varchar(50) not null comment '岗位名称',
    post_sort int(4) not null comment '显示顺序',
    remark varchar(500) default null comment '备注',
    primary key (id)
) engine = innodb comment = '岗位信息表' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 4、角色信息表
-- ----------------------------
drop table if exists sys_role;

create table sys_role (
    id bigint(20) not null auto_increment comment '角色id',
    extend_shadow text comment '扩展字段',
    status char(1) default '0' comment '角色状态（0正常 1停用）',
    created_by varchar(64) default '' comment '创建者',
    created_at dateTime comment '创建时间',
    updated_by varchar(64) default '' comment '创建时间',
    updated_at dateTime comment '更新时间',
    role_name varchar(30) not null comment '角色名称',
    role_key varchar(100) not null comment '角色权限字符串',
    role_sort int(4) not null comment '显示顺序',
    data_scope char(1) default '1' comment '数据范围（1：全部数据权限 2：自定数据权限 3：本部门数据权限 4：本部门及以下数据权限）',
    menu_check_strictly tinyint(1) default 1 comment '菜单树选择项是否关联显示',
    dept_check_strictly tinyint(1) default 1 comment '部门树选择项是否关联显示',
    del_flag char(1) default '0' comment '删除标志（0代表存在 2代表删除）',
    remark varchar(500) default null comment '备注',
    primary key (id)
) engine = innodb auto_increment = 100 comment = '角色信息表' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 5、菜单权限表
-- ----------------------------
drop table if exists sys_menu;

create table sys_menu (
    id bigint(20) not null auto_increment comment '菜单id',
    extend_shadow text comment '扩展字段',
    status char(1) default '0' comment '菜单状态（0正常 1停用）',
    created_by varchar(64) default '' comment '创建者',
    created_at dateTime comment '创建时间',
    updated_by varchar(64) default '' comment '创建时间',
    updated_at dateTime comment '更新时间',
    menu_name varchar(50) not null comment '菜单名称',
    parent_id bigint(20) default 0 comment '父菜单ID',
    order_num int(4) default 0 comment '显示顺序',
    path varchar(200) default '' comment '路由地址',
    component varchar(255) default null comment '组件路径',
    quexf varchar(255) default null comment '路由参数',
    is_frame int(1) default 1 comment '是否为外链（0是 1否）',
    is_cache int(1) default 0 comment '是否缓存（0缓存 1不缓存）',
    menu_type char(1) default '' comment '菜单类型（M目录 C菜单 F按钮）',
    visible char(1) default 0 comment '菜单状态（0显示 1隐藏）',
    perms varchar(100) default null comment '权限标识',
    icon varchar(100) default '#' comment '菜单图标',
    remark varchar(500) default '' comment '备注',
    primary key (id),
    CONSTRAINT `fk_menu_menu` FOREIGN key (`parent_id`) REFERENCES `sys_menu` (`id`) ON DELETE NO ACTION ON UPDATE NO ACTION
) engine = innodb auto_increment = 2000 comment = '菜单权限表' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 6、用户和角色关联表  用户N-1角色
-- ----------------------------
drop table if exists sys_user_role;

create table sys_user_role (
    user_id bigint(20) not null comment '用户ID',
    role_id bigint(20) not null comment '角色ID',
    primary key(user_id, role_id)
) engine = innodb comment = '用户和角色关联表' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 7、角色和菜单关联表  角色1-N菜单
-- ----------------------------
drop table if exists sys_role_menu;

create table sys_role_menu (
    role_id bigint(20) not null comment '角色ID',
    menu_id bigint(20) not null comment '菜单ID',
    primary key(role_id, menu_id)
) engine = innodb comment = '角色和菜单关联表' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 8、角色和部门关联表  角色1-N部门
-- ----------------------------
drop table if exists sys_role_dept;

create table sys_role_dept (
    role_id bigint(20) not null comment '角色ID',
    dept_id bigint(20) not null comment '部门ID',
    primary key(role_id, dept_id)
) engine = innodb comment = '角色和部门关联表' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 9、用户与岗位关联表  用户1-N岗位
-- ----------------------------
drop table if exists sys_user_post;

create table sys_user_post (
    user_id bigint(20) not null comment '用户ID',
    post_id bigint(20) not null comment '岗位ID',
    primary key (user_id, post_id)
) engine = innodb comment = '用户与岗位关联表' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 10、操作日志记录
-- ----------------------------
drop table if exists sys_oper_log;

create table sys_oper_log
(
    id             bigint(20) not null auto_increment comment '日志主键',
    title          varchar(50)   default '' comment '模块标题',
    business_type  int(2) default 0 comment '业务类型（0其它 1新增 2修改 3删除）',
    method         varchar(100)  default '' comment '方法名称',
    request_method varchar(10)   default '' comment '请求方式',
    operator_type  int(1) default 0 comment '操作类别（0其它 1后台用户 2手机端用户）',
    oper_name      varchar(50)   default '' comment '操作人员',
    dept_name      varchar(50)   default '' comment '部门名称',
    oper_url       varchar(255)  default '' comment '请求URL',
    oper_ip        varchar(128)  default '' comment '主机地址',
    oper_location  varchar(255)  default '' comment '操作地点',
    oper_param     varchar(2000) default '' comment '请求参数',
    json_result    varchar(2000) default '' comment '返回参数',
    status         int(1) default 0 comment '操作状态（0正常 1异常）',
    error_msg      varchar(2000) default '' comment '错误消息',
    oper_time      datetime comment '操作时间',
    cost_time      bigint(20) default 0 comment '消耗时间',
    primary key (id)
) engine = innodb auto_increment = 100 comment = '操作日志记录' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 11、字典类型表
-- ----------------------------
drop table if exists sys_dict_type;

create table sys_dict_type (
    id bigint(20) not null auto_increment comment '字典主键',
    extend_shadow text comment '扩展字段',
    status char(1) default '0' comment '状态（0正常 1停用）',
    create_by varchar(64) default '' comment '创建者',
    create_time datetime comment '创建时间',
    update_by varchar(64) default '' comment '更新者',
    update_time datetime comment '更新时间',
    dict_name varchar(100) default '' comment '字典名称',
    dict_type varchar(100) default '' comment '字典类型',
    remark varchar(500) default null comment '备注',
    primary key (id),
    unique (dict_type)
) engine = innodb auto_increment = 100 comment = '字典类型表' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 12、字典数据表
-- ----------------------------
drop table if exists sys_dict_data;

create table sys_dict_data (
    id bigint(20) not null auto_increment comment '字典主键',
    extend_shadow text comment '扩展字段',
    status char(1) default '0' comment '状态（0正常 1停用）',
    create_by varchar(64) default '' comment '创建者',
    create_time datetime comment '创建时间',
    update_by varchar(64) default '' comment '更新者',
    update_time datetime comment '更新时间',
    dict_code bigint(20) not null unique comment '字典编码',
    dict_sort int(4) default 0 comment '字典排序',
    dict_label varchar(100) default '' comment '字典标签',
    dict_value varchar(100) default '' comment '字典键值',
    dict_type varchar(100) default '' comment '字典类型',
    css_class varchar(100) default null comment '样式属性（其他样式扩展）',
    list_class varchar(100) default null comment '表格回显样式',
    is_default char(1) default 'N' comment '是否默认（Y是 N否）',
    remark varchar(500) default null comment '备注',
    primary key (id)
) engine = innodb auto_increment = 100 comment = '字典数据表' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 13、参数配置表
-- ----------------------------
drop table if exists sys_config;

create table sys_config (
    id bigint(20) not null auto_increment comment '参数主键',
    extend_shadow text comment '扩展字段',
    status char(1) default '0' comment '状态（0正常 1停用）',
    create_by varchar(64) default '' comment '创建者',
    create_time datetime comment '创建时间',
    update_by varchar(64) default '' comment '更新者',
    update_time datetime comment '更新时间',
    config_name varchar(100) default '' comment '参数名称',
    config_key varchar(100) default '' comment '参数键名',
    config_value varchar(500) default '' comment '参数键值',
    config_type char(1) default 'N' comment '系统内置（Y是 N否）',
    remark varchar(500) default null comment '备注',
    primary key (id)
) engine = innodb auto_increment = 100 comment = '参数配置表' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 14、系统访问记录
-- ----------------------------
drop table if exists sys_logininfor;

create table sys_logininfor (
    id bigint(20) not null auto_increment comment '访问ID',
    user_name varchar(50) default '' comment '用户账号',
    ipaddr varchar(128) default '' comment '登录IP地址',
    status char(1) default '0' comment '登录状态（0成功 1失败）',
    msg varchar(255) default '' comment '提示信息',
    access_time datetime comment '访问时间',
    primary key (id)
) engine = innodb auto_increment = 100 comment='系统访问记录' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 15、定时任务调度表
-- ----------------------------
drop table if exists sys_job;

create table sys_job (
    id bigint(20) not null auto_increment comment '定时任务主键',
    extend_shadow text comment '扩展字段',
    status char(1) default '0' comment '状态（0正常 1停用）',
    create_by varchar(64) default '' comment '创建者',
    create_time datetime comment '创建时间',
    update_by varchar(64) default '' comment '更新者',
    update_time datetime comment '更新时间',
    job_name varchar(64) default '' comment '任务名称',
    job_group varchar(64) default 'DEFAULT' comment '任务组名',
    invoke_target varchar(500) not null comment '调用目标字符串',
    cron_expression varchar(255) default '' comment 'cron执行表达式',
    misfire_policy varchar(20) default '3' comment '计划执行错误策略（1立即执行 2执行一次 3放弃执行）',
    concurrent char(1) default '1' comment '是否并发执行（0允许 1禁止）',
    remark varchar(500) default '' comment '备注信息',
    primary key (id, job_name, job_group)
) engine = innodb auto_increment = 100 comment = '定时任务调度表' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 16、定时任务调度日志表
-- ----------------------------
drop table if exists sys_job_log;

create table sys_job_log (
    id bigint(20) not null auto_increment comment '任务日志ID',
    job_name varchar(64) not null comment '任务名称',
    job_group varchar(64) not null comment '任务组名',
    invoke_target varchar(500) not null comment '调用目标字符串',
    job_message varchar(500) comment '日志信息',
    status char(1) default '0' comment '执行状态（0正常 1失败）',
    exception_info varchar(2000) default '' comment '异常信息',
    create_time datetime comment '创建时间',
    primary key (id)
) engine = innodb comment = '定时任务调度日志表' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 17、通知公告表
-- ----------------------------
drop table if exists sys_notice;

create table sys_notice (
    id bigint(20) not null auto_increment comment '公告id',
    extend_shadow text comment '扩展字段',
    status char(1) default '0' comment '状态（0正常 1停用）',
    create_by varchar(64) default '' comment '创建者',
    create_time datetime comment '创建时间',
    update_by varchar(64) default '' comment '更新者',
    update_time datetime comment '更新时间',
    notice_title varchar(50) not null comment '公告标题',
    notice_type char(1) not null comment '公告类型（1通知 2公告）',
    notice_content longblob default null comment '公告内容',
    remark varchar(255) default null comment '备注',
    primary key (id)
) engine = innodb auto_increment = 10 comment = '通知公告表' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 18、代码生成业务表
-- ----------------------------
drop table if exists gen_table;

create table gen_table (
    id bigint(20) not null auto_increment comment '业务id',
    extend_shadow text comment '扩展字段',
    status char(1) default '0' comment '状态（0正常 1停用）',
    create_by varchar(64) default '' comment '创建者',
    create_time datetime comment '创建时间',
    update_by varchar(64) default '' comment '更新者',
    update_time datetime comment '更新时间',
    table_name varchar(200) default '' comment '表名称',
    table_comment varchar(500) default '' comment '表描述',
    sub_table_name varchar(64) default null comment '关联子表的表名',
    sub_table_fk_name varchar(64) default null comment '子表关联的外键名',
    class_name varchar(100) default '' comment '实体类名称',
    tpl_categoxf varchar(200) default 'crud' comment '使用的模板（crud单表操作 tree树表操作）',
    package_name varchar(100) comment '生成包路径',
    module_name varchar(30) comment '生成模块名',
    business_name varchar(30) comment '生成业务名',
    function_name varchar(50) comment '生成功能名',
    function_author varchar(50) comment '生成功能作者',
    gen_type char(1) default '0' comment '生成代码方式（0zip压缩包 1自定义路径）',
    gen_path varchar(200) default '/' comment '生成路径（不填默认项目路径）',
    options varchar(1000) comment '其它生成选项',
    remark varchar(500) default null comment '备注',
    primary key (id)
) engine = innodb auto_increment = 1 comment = '代码生成业务表' charset=utf8mb4 collate=utf8mb4_unicode_ci;

-- ----------------------------
-- 19、代码生成业务表字段
-- ----------------------------
drop table if exists gen_table_column;

create table gen_table_column (
    id bigint(20) not null auto_increment comment '代码生成业务id',
    extend_shadow text comment '扩展字段',
    status char(1) default '0' comment '状态（0正常 1停用）',
    create_by varchar(64) default '' comment '创建者',
    create_time datetime comment '创建时间',
    update_by varchar(64) default '' comment '更新者',
    update_time datetime comment '更新时间',
    table_id bigint(20) comment '归属表编号',
    column_name varchar(200) comment '列名称',
    column_comment varchar(500) comment '列描述',
    column_type varchar(100) comment '列类型',
    java_type varchar(500) comment 'JAVA类型',
    java_field varchar(200) comment 'JAVA字段名',
    is_pk char(1) comment '是否主键（1是）',
    is_increment char(1) comment '是否自增（1是）',
    is_required char(1) comment '是否必填（1是）',
    is_insert char(1) comment '是否为插入字段（1是）',
    is_edit char(1) comment '是否编辑字段（1是）',
    is_list char(1) comment '是否列表字段（1是）',
    is_quexf char(1) comment '是否查询字段（1是）',
    quexf_type varchar(200) default 'EQ' comment '查询方式（等于、不等于、大于、小于、范围）',
    html_type varchar(200) comment '显示类型（文本框、文本域、下拉框、复选框、单选框、日期控件）',
    dict_type varchar(200) default '' comment '字典类型',
    sort int comment '排序',
    primary key (id)
) engine = innodb auto_increment = 1 comment = '代码生成业务表字段' charset=utf8mb4 collate=utf8mb4_unicode_ci;