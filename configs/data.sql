use
`sophie`;

-- ----------------------------
-- 初始化-部门表数据
-- ----------------------------
insert into sys_dept
values (100, 0, '0', 'sophie', 0, '雪菲', '15888888888', 'sophie@qq.com', '0', '0', 'admin', sysdate(), '', null,
        null, '');
insert into sys_dept
values (101, 100, '0,100', '深圳总公司', 1, '雪菲', '15888888888', 'sophie@qq.com', '0', '0', 'admin', sysdate(), '',
        null, null, '');
insert into sys_dept
values (102, 100, '0,100', '长沙分公司', 2, '雪菲', '15888888888', 'sophie@qq.com', '0', '0', 'admin', sysdate(), '',
        null, null, '');
insert into sys_dept
values (103, 101, '0,100,101', '研发部门', 1, '雪菲', '15888888888', 'sophie@qq.com', '0', '0', 'admin', sysdate(), '',
        null, null, '');
insert into sys_dept
values (104, 101, '0,100,101', '市场部门', 2, '雪菲', '15888888888', 'sophie@qq.com', '0', '0', 'admin', sysdate(), '',
        null, null, '');
insert into sys_dept
values (105, 101, '0,100,101', '测试部门', 3, '雪菲', '15888888888', 'sophie@qq.com', '0', '0', 'admin', sysdate(), '',
        null, null, '');
insert into sys_dept
values (106, 101, '0,100,101', '财务部门', 4, '雪菲', '15888888888', 'sophie@qq.com', '0', '0', 'admin', sysdate(), '',
        null, null, '');
insert into sys_dept
values (107, 101, '0,100,101', '运维部门', 5, '雪菲', '15888888888', 'sophie@qq.com', '0', '0', 'admin', sysdate(), '',
        null, null, '');
insert into sys_dept
values (108, 102, '0,100,102', '市场部门', 1, '雪菲', '15888888888', 'sophie@qq.com', '0', '0', 'admin', sysdate(), '',
        null, null, '');
insert into sys_dept
values (109, 102, '0,100,102', '财务部门', 2, '雪菲', '15888888888', 'sophie@qq.com', '0', '0', 'admin', sysdate(), '',
        null, null, '');

-- ----------------------------
-- 初始化-用户信息表数据
-- ----------------------------
insert into sys_user
values (1, 103, 'admin', '雪菲', '00', 'sophie@163.com', '15888888888', '1', '',
        '$2a$10$7JB720yubVSZvUI0rEqK/.VqGOZTH.ulu33dHOiBE8ByOhJIrdAu2', '0', '0', '127.0.0.1', sysdate(), 'admin',
        sysdate(), '', null, '管理员', null);
insert into sys_user
values (2, 105, 'sophie', '雪菲', '00', 'sophie@qq.com', '15666666666', '1', '',
        '$2a$10$7JB720yubVSZvUI0rEqK/.VqGOZTH.ulu33dHOiBE8ByOhJIrdAu2', '0', '0', '127.0.0.1', sysdate(), 'admin',
        sysdate(), '', null, '测试员', null);

-- ----------------------------
-- 初始化-岗位信息表数据
-- ----------------------------
insert into sys_post
values (1, 'ceo', '董事长', 1, '0', 'admin', sysdate(), '', null, '', null);
insert into sys_post
values (2, 'se', '项目经理', 2, '0', 'admin', sysdate(), '', null, '', null);
insert into sys_post
values (3, 'hr', '人力资源', 3, '0', 'admin', sysdate(), '', null, '', null);
insert into sys_post
values (4, 'user', '普通员工', 4, '0', 'admin', sysdate(), '', null, '', null);


-- ----------------------------
-- 初始化-角色信息表数据
-- ----------------------------
insert into sys_role
values ('1', '超级管理员', 'admin', 1, 1, 1, 1, '0', '0', 'admin', sysdate(), '', null, '超级管理员', null);
insert into sys_role
values ('2', '普通角色', 'common', 2, 2, 1, 1, '0', '0', 'admin', sysdate(), '', null, '普通角色', null);


-- ----------------------------
-- 初始化-菜单信息表数据
-- ----------------------------
-- 一级菜单
insert into sys_menu
values ('1', '系统管理', '0', '1', 'system', null, '', 1, 0, 'M', '0', '0', '', 'system', 'admin', sysdate(), '', null,
        '系统管理目录', null);
insert into sys_menu
values ('2', '系统监控', '0', '2', 'monitor', null, '', 1, 0, 'M', '0', '0', '', 'monitor', 'admin', sysdate(), '',
        null, '系统监控目录', null);
insert into sys_menu
values ('3', '系统工具', '0', '3', 'tool', null, '', 1, 0, 'M', '0', '0', '', 'tool', 'admin', sysdate(), '', null,
        '系统工具目录', null);
insert into sys_menu
values ('4', 'GitHub', '0', '4', 'https://github.com/user823', null, '', 0, 0, 'M', '0', '0', '', 'guide', 'admin', sysdate(),
        '', null, 'github官网地址', null);
-- 二级菜单
insert into sys_menu
values ('100', '用户管理', '1', '1', 'user', 'system/user/index', '', 1, 0, 'C', '0', '0', 'system:user:list', 'user',
        'admin', sysdate(), '', null, '用户管理菜单', null);
insert into sys_menu
values ('101', '角色管理', '1', '2', 'role', 'system/role/index', '', 1, 0, 'C', '0', '0', 'system:role:list',
        'peoples', 'admin', sysdate(), '', null, '角色管理菜单', null);
insert into sys_menu
values ('102', '菜单管理', '1', '3', 'menu', 'system/menu/index', '', 1, 0, 'C', '0', '0', 'system:menu:list',
        'tree-table', 'admin', sysdate(), '', null, '菜单管理菜单', null);
insert into sys_menu
values ('103', '部门管理', '1', '4', 'dept', 'system/dept/index', '', 1, 0, 'C', '0', '0', 'system:dept:list', 'tree',
        'admin', sysdate(), '', null, '部门管理菜单', null);
insert into sys_menu
values ('104', '岗位管理', '1', '5', 'post', 'system/post/index', '', 1, 0, 'C', '0', '0', 'system:post:list', 'post',
        'admin', sysdate(), '', null, '岗位管理菜单', null);
insert into sys_menu
values ('105', '字典管理', '1', '6', 'dict', 'system/dict/index', '', 1, 0, 'C', '0', '0', 'system:dict:list', 'dict',
        'admin', sysdate(), '', null, '字典管理菜单', null);
insert into sys_menu
values ('106', '参数设置', '1', '7', 'config', 'system/config/index', '', 1, 0, 'C', '0', '0', 'system:config:list',
        'edit', 'admin', sysdate(), '', null, '参数设置菜单', null);
insert into sys_menu
values ('107', '通知公告', '1', '8', 'notice', 'system/notice/index', '', 1, 0, 'C', '0', '0', 'system:notice:list',
        'message', 'admin', sysdate(), '', null, '通知公告菜单', null);
insert into sys_menu
values ('108', '日志管理', '1', '9', 'log', '', '', 1, 0, 'M', '0', '0', '', 'log', 'admin', sysdate(), '', null,
        '日志管理菜单', null);
insert into sys_menu
values ('109', '在线用户', '2', '1', 'online', 'monitor/online/index', '', 1, 0, 'C', '0', '0', 'monitor:online:list',
        'online', 'admin', sysdate(), '', null, '在线用户菜单', null);
insert into sys_menu
values ('110', '定时任务', '2', '2', 'job', 'monitor/job/index', '', 1, 0, 'C', '0', '0', 'monitor:job:list', 'job',
        'admin', sysdate(), '', null, '定时任务菜单', null);
insert into sys_menu
values ('111', 'tracing', '2', '3', 'http://localhost:16686', '', '', 0, 0, 'C', '0', '0',
        'monitor:tracing:list', 'sentinel', 'admin', sysdate(), '', null, '链路追踪菜单', null);
insert into sys_menu
values ('114', '表单构建', '3', '1', 'build', 'tool/build/index', '', 1, 0, 'C', '0', '0', 'tool:build:list', 'build',
        'admin', sysdate(), '', null, '表单构建菜单', null);
insert into sys_menu
values ('115', '代码生成', '3', '2', 'gen', 'tool/gen/index', '', 1, 0, 'C', '0', '0', 'tool:gen:list', 'code', 'admin',
        sysdate(), '', null, '代码生成菜单', null);
insert into sys_menu
values ('116', '系统接口', '3', '3', 'http://localhost:8082/swagger/index.html', '', '', 0, 0, 'C', '0', '0',
        'tool:swagger:list', 'swagger', 'admin', sysdate(), '', null, '系统接口菜单', null);
-- 三级菜单
insert into sys_menu
values ('500', '操作日志', '108', '1', 'operlog', 'system/operlog/index', '', 1, 0, 'C', '0', '0',
        'system:operlog:list', 'form', 'admin', sysdate(), '', null, '操作日志菜单', null);
insert into sys_menu
values ('501', '登录日志', '108', '2', 'logininfor', 'system/logininfor/index', '', 1, 0, 'C', '0', '0',
        'system:logininfor:list', 'logininfor', 'admin', sysdate(), '', null, '登录日志菜单', null);
-- 用户管理按钮
insert into sys_menu
values ('1000', '用户查询', '100', '1', '', '', '', 1, 0, 'F', '0', '0', 'system:user:query', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1001', '用户新增', '100', '2', '', '', '', 1, 0, 'F', '0', '0', 'system:user:add', '#', 'admin', sysdate(), '',
        null, '', null);
insert into sys_menu
values ('1002', '用户修改', '100', '3', '', '', '', 1, 0, 'F', '0', '0', 'system:user:edit', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1003', '用户删除', '100', '4', '', '', '', 1, 0, 'F', '0', '0', 'system:user:remove', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1004', '用户导出', '100', '5', '', '', '', 1, 0, 'F', '0', '0', 'system:user:export', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1005', '用户导入', '100', '6', '', '', '', 1, 0, 'F', '0', '0', 'system:user:import', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1006', '重置密码', '100', '7', '', '', '', 1, 0, 'F', '0', '0', 'system:user:resetPwd', '#', 'admin',
        sysdate(), '', null, '', null);
-- 角色管理按钮
insert into sys_menu
values ('1007', '角色查询', '101', '1', '', '', '', 1, 0, 'F', '0', '0', 'system:role:query', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1008', '角色新增', '101', '2', '', '', '', 1, 0, 'F', '0', '0', 'system:role:add', '#', 'admin', sysdate(), '',
        null, '', null);
insert into sys_menu
values ('1009', '角色修改', '101', '3', '', '', '', 1, 0, 'F', '0', '0', 'system:role:edit', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1010', '角色删除', '101', '4', '', '', '', 1, 0, 'F', '0', '0', 'system:role:remove', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1011', '角色导出', '101', '5', '', '', '', 1, 0, 'F', '0', '0', 'system:role:export', '#', 'admin', sysdate(),
        '', null, '', null);
-- 菜单管理按钮
insert into sys_menu
values ('1012', '菜单查询', '102', '1', '', '', '', 1, 0, 'F', '0', '0', 'system:menu:query', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1013', '菜单新增', '102', '2', '', '', '', 1, 0, 'F', '0', '0', 'system:menu:add', '#', 'admin', sysdate(), '',
        null, '', null);
insert into sys_menu
values ('1014', '菜单修改', '102', '3', '', '', '', 1, 0, 'F', '0', '0', 'system:menu:edit', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1015', '菜单删除', '102', '4', '', '', '', 1, 0, 'F', '0', '0', 'system:menu:remove', '#', 'admin', sysdate(),
        '', null, '', null);
-- 部门管理按钮
insert into sys_menu
values ('1016', '部门查询', '103', '1', '', '', '', 1, 0, 'F', '0', '0', 'system:dept:query', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1017', '部门新增', '103', '2', '', '', '', 1, 0, 'F', '0', '0', 'system:dept:add', '#', 'admin', sysdate(), '',
        null, '', null);
insert into sys_menu
values ('1018', '部门修改', '103', '3', '', '', '', 1, 0, 'F', '0', '0', 'system:dept:edit', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1019', '部门删除', '103', '4', '', '', '', 1, 0, 'F', '0', '0', 'system:dept:remove', '#', 'admin', sysdate(),
        '', null, '', null);
-- 岗位管理按钮
insert into sys_menu
values ('1020', '岗位查询', '104', '1', '', '', '', 1, 0, 'F', '0', '0', 'system:post:query', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1021', '岗位新增', '104', '2', '', '', '', 1, 0, 'F', '0', '0', 'system:post:add', '#', 'admin', sysdate(), '',
        null, '', null);
insert into sys_menu
values ('1022', '岗位修改', '104', '3', '', '', '', 1, 0, 'F', '0', '0', 'system:post:edit', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1023', '岗位删除', '104', '4', '', '', '', 1, 0, 'F', '0', '0', 'system:post:remove', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1024', '岗位导出', '104', '5', '', '', '', 1, 0, 'F', '0', '0', 'system:post:export', '#', 'admin', sysdate(),
        '', null, '', null);
-- 字典管理按钮
insert into sys_menu
values ('1025', '字典查询', '105', '1', '#', '', '', 1, 0, 'F', '0', '0', 'system:dict:query', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1026', '字典新增', '105', '2', '#', '', '', 1, 0, 'F', '0', '0', 'system:dict:add', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1027', '字典修改', '105', '3', '#', '', '', 1, 0, 'F', '0', '0', 'system:dict:edit', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1028', '字典删除', '105', '4', '#', '', '', 1, 0, 'F', '0', '0', 'system:dict:remove', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1029', '字典导出', '105', '5', '#', '', '', 1, 0, 'F', '0', '0', 'system:dict:export', '#', 'admin', sysdate(),
        '', null, '', null);
-- 参数设置按钮
insert into sys_menu
values ('1030', '参数查询', '106', '1', '#', '', '', 1, 0, 'F', '0', '0', 'system:config:query', '#', 'admin',
        sysdate(), '', null, '', null);
insert into sys_menu
values ('1031', '参数新增', '106', '2', '#', '', '', 1, 0, 'F', '0', '0', 'system:config:add', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1032', '参数修改', '106', '3', '#', '', '', 1, 0, 'F', '0', '0', 'system:config:edit', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1033', '参数删除', '106', '4', '#', '', '', 1, 0, 'F', '0', '0', 'system:config:remove', '#', 'admin',
        sysdate(), '', null, '', null);
insert into sys_menu
values ('1034', '参数导出', '106', '5', '#', '', '', 1, 0, 'F', '0', '0', 'system:config:export', '#', 'admin',
        sysdate(), '', null, '', null);
-- 通知公告按钮
insert into sys_menu
values ('1035', '公告查询', '107', '1', '#', '', '', 1, 0, 'F', '0', '0', 'system:notice:query', '#', 'admin',
        sysdate(), '', null, '', null);
insert into sys_menu
values ('1036', '公告新增', '107', '2', '#', '', '', 1, 0, 'F', '0', '0', 'system:notice:add', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1037', '公告修改', '107', '3', '#', '', '', 1, 0, 'F', '0', '0', 'system:notice:edit', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1038', '公告删除', '107', '4', '#', '', '', 1, 0, 'F', '0', '0', 'system:notice:remove', '#', 'admin',
        sysdate(), '', null, '', null);
-- 操作日志按钮
insert into sys_menu
values ('1039', '操作查询', '500', '1', '#', '', '', 1, 0, 'F', '0', '0', 'system:operlog:query', '#', 'admin',
        sysdate(), '', null, '', null);
insert into sys_menu
values ('1040', '操作删除', '500', '2', '#', '', '', 1, 0, 'F', '0', '0', 'system:operlog:remove', '#', 'admin',
        sysdate(), '', null, '', null);
insert into sys_menu
values ('1041', '日志导出', '500', '3', '#', '', '', 1, 0, 'F', '0', '0', 'system:operlog:export', '#', 'admin',
        sysdate(), '', null, '', null);
-- 登录日志按钮
insert into sys_menu
values ('1042', '登录查询', '501', '1', '#', '', '', 1, 0, 'F', '0', '0', 'system:logininfor:query', '#', 'admin',
        sysdate(), '', null, '', null);
insert into sys_menu
values ('1043', '登录删除', '501', '2', '#', '', '', 1, 0, 'F', '0', '0', 'system:logininfor:remove', '#', 'admin',
        sysdate(), '', null, '', null);
insert into sys_menu
values ('1044', '日志导出', '501', '3', '#', '', '', 1, 0, 'F', '0', '0', 'system:logininfor:export', '#', 'admin',
        sysdate(), '', null, '', null);
insert into sys_menu
values ('1045', '账户解锁', '501', '4', '#', '', '', 1, 0, 'F', '0', '0', 'system:logininfor:unlock', '#', 'admin',
        sysdate(), '', null, '', null);
-- 在线用户按钮
insert into sys_menu
values ('1046', '在线查询', '109', '1', '#', '', '', 1, 0, 'F', '0', '0', 'monitor:online:query', '#', 'admin',
        sysdate(), '', null, '', null);
insert into sys_menu
values ('1047', '批量强退', '109', '2', '#', '', '', 1, 0, 'F', '0', '0', 'monitor:online:batchLogout', '#', 'admin',
        sysdate(), '', null, '', null);
insert into sys_menu
values ('1048', '单条强退', '109', '3', '#', '', '', 1, 0, 'F', '0', '0', 'monitor:online:forceLogout', '#', 'admin',
        sysdate(), '', null, '', null);
-- 定时任务按钮
insert into sys_menu
values ('1049', '任务查询', '110', '1', '#', '', '', 1, 0, 'F', '0', '0', 'monitor:job:query', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1050', '任务新增', '110', '2', '#', '', '', 1, 0, 'F', '0', '0', 'monitor:job:add', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1051', '任务修改', '110', '3', '#', '', '', 1, 0, 'F', '0', '0', 'monitor:job:edit', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1052', '任务删除', '110', '4', '#', '', '', 1, 0, 'F', '0', '0', 'monitor:job:remove', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1053', '状态修改', '110', '5', '#', '', '', 1, 0, 'F', '0', '0', 'monitor:job:changeStatus', '#', 'admin',
        sysdate(), '', null, '', null);
insert into sys_menu
values ('1054', '任务导出', '110', '6', '#', '', '', 1, 0, 'F', '0', '0', 'monitor:job:export', '#', 'admin', sysdate(),
        '', null, '', null);
-- 代码生成按钮
insert into sys_menu
values ('1055', '生成查询', '115', '1', '#', '', '', 1, 0, 'F', '0', '0', 'tool:gen:query', '#', 'admin', sysdate(), '',
        null, '', null);
insert into sys_menu
values ('1056', '生成修改', '115', '2', '#', '', '', 1, 0, 'F', '0', '0', 'tool:gen:edit', '#', 'admin', sysdate(), '',
        null, '', null);
insert into sys_menu
values ('1057', '生成删除', '115', '3', '#', '', '', 1, 0, 'F', '0', '0', 'tool:gen:remove', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1058', '导入代码', '115', '2', '#', '', '', 1, 0, 'F', '0', '0', 'tool:gen:import', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1059', '预览代码', '115', '4', '#', '', '', 1, 0, 'F', '0', '0', 'tool:gen:preview', '#', 'admin', sysdate(),
        '', null, '', null);
insert into sys_menu
values ('1060', '生成代码', '115', '5', '#', '', '', 1, 0, 'F', '0', '0', 'tool:gen:code', '#', 'admin', sysdate(), '',
        null, '', null);


-- ----------------------------
-- 初始化-用户和角色关联表数据
-- ----------------------------
insert into sys_user_role
values ('1', '1', null);
insert into sys_user_role
values ('2', '2', null);


-- ----------------------------
-- 初始化-角色和菜单关联表数据
-- ----------------------------
insert into sys_role_menu
values ('2', '1', null);
insert into sys_role_menu
values ('2', '2', null);
insert into sys_role_menu
values ('2', '3', null);
insert into sys_role_menu
values ('2', '4', null);
insert into sys_role_menu
values ('2', '100', null);
insert into sys_role_menu
values ('2', '101', null);
insert into sys_role_menu
values ('2', '102', null);
insert into sys_role_menu
values ('2', '103', null);
insert into sys_role_menu
values ('2', '104', null);
insert into sys_role_menu
values ('2', '105', null);
insert into sys_role_menu
values ('2', '106', null);
insert into sys_role_menu
values ('2', '107', null);
insert into sys_role_menu
values ('2', '108', null);
insert into sys_role_menu
values ('2', '109', null);
insert into sys_role_menu
values ('2', '110', null);
insert into sys_role_menu
values ('2', '111', null);
insert into sys_role_menu
values ('2', '112', null);
insert into sys_role_menu
values ('2', '113', null);
insert into sys_role_menu
values ('2', '114', null);
insert into sys_role_menu
values ('2', '115', null);
insert into sys_role_menu
values ('2', '116', null);
insert into sys_role_menu
values ('2', '500', null);
insert into sys_role_menu
values ('2', '501', null);
insert into sys_role_menu
values ('2', '1000', null);
insert into sys_role_menu
values ('2', '1001', null);
insert into sys_role_menu
values ('2', '1002', null);
insert into sys_role_menu
values ('2', '1003', null);
insert into sys_role_menu
values ('2', '1004', null);
insert into sys_role_menu
values ('2', '1005', null);
insert into sys_role_menu
values ('2', '1006', null);
insert into sys_role_menu
values ('2', '1007', null);
insert into sys_role_menu
values ('2', '1008', null);
insert into sys_role_menu
values ('2', '1009', null);
insert into sys_role_menu
values ('2', '1010', null);
insert into sys_role_menu
values ('2', '1011', null);
insert into sys_role_menu
values ('2', '1012', null);
insert into sys_role_menu
values ('2', '1013', null);
insert into sys_role_menu
values ('2', '1014', null);
insert into sys_role_menu
values ('2', '1015', null);
insert into sys_role_menu
values ('2', '1016', null);
insert into sys_role_menu
values ('2', '1017', null);
insert into sys_role_menu
values ('2', '1018', null);
insert into sys_role_menu
values ('2', '1019', null);
insert into sys_role_menu
values ('2', '1020', null);
insert into sys_role_menu
values ('2', '1021', null);
insert into sys_role_menu
values ('2', '1022', null);
insert into sys_role_menu
values ('2', '1023', null);
insert into sys_role_menu
values ('2', '1024', null);
insert into sys_role_menu
values ('2', '1025', null);
insert into sys_role_menu
values ('2', '1026', null);
insert into sys_role_menu
values ('2', '1027', null);
insert into sys_role_menu
values ('2', '1028', null);
insert into sys_role_menu
values ('2', '1029', null);
insert into sys_role_menu
values ('2', '1030', null);
insert into sys_role_menu
values ('2', '1031', null);
insert into sys_role_menu
values ('2', '1032', null);
insert into sys_role_menu
values ('2', '1033', null);
insert into sys_role_menu
values ('2', '1034', null);
insert into sys_role_menu
values ('2', '1035', null);
insert into sys_role_menu
values ('2', '1036', null);
insert into sys_role_menu
values ('2', '1037', null);
insert into sys_role_menu
values ('2', '1038', null);
insert into sys_role_menu
values ('2', '1039', null);
insert into sys_role_menu
values ('2', '1040', null);
insert into sys_role_menu
values ('2', '1041', null);
insert into sys_role_menu
values ('2', '1042', null);
insert into sys_role_menu
values ('2', '1043', null);
insert into sys_role_menu
values ('2', '1044', null);
insert into sys_role_menu
values ('2', '1045', null);
insert into sys_role_menu
values ('2', '1046', null);
insert into sys_role_menu
values ('2', '1047', null);
insert into sys_role_menu
values ('2', '1048', null);
insert into sys_role_menu
values ('2', '1049', null);
insert into sys_role_menu
values ('2', '1050', null);
insert into sys_role_menu
values ('2', '1051', null);
insert into sys_role_menu
values ('2', '1052', null);
insert into sys_role_menu
values ('2', '1053', null);
insert into sys_role_menu
values ('2', '1054', null);
insert into sys_role_menu
values ('2', '1055', null);
insert into sys_role_menu
values ('2', '1056', null);
insert into sys_role_menu
values ('2', '1057', null);
insert into sys_role_menu
values ('2', '1058', null);
insert into sys_role_menu
values ('2', '1059', null);
insert into sys_role_menu
values ('2', '1060', null);


-- ----------------------------
-- 初始化-角色和部门关联表数据
-- ----------------------------
insert into sys_role_dept
values ('2', '100', null);
insert into sys_role_dept
values ('2', '101', null);
insert into sys_role_dept
values ('2', '105', null);


-- ----------------------------
-- 初始化-用户与岗位关联表数据
-- ----------------------------
insert into sys_user_post
values ('1', '1', null);
insert into sys_user_post
values ('2', '2', null);

-- ----------------------------
-- 初始化-字典类型和字典数据
-- ----------------------------

insert into sys_dict_type
values (1, '用户性别', 'sys_user_sex', '0', 'admin', sysdate(), '', null, '用户性别列表', null);
insert into sys_dict_type
values (2, '菜单状态', 'sys_show_hide', '0', 'admin', sysdate(), '', null, '菜单状态列表', null);
insert into sys_dict_type
values (3, '系统开关', 'sys_normal_disable', '0', 'admin', sysdate(), '', null, '系统开关列表', null);
insert into sys_dict_type
values (4, '任务状态', 'sys_job_status', '0', 'admin', sysdate(), '', null, '任务状态列表', null);
insert into sys_dict_type
values (5, '任务分组', 'sys_job_group', '0', 'admin', sysdate(), '', null, '任务分组列表', null);
insert into sys_dict_type
values (6, '系统是否', 'sys_yes_no', '0', 'admin', sysdate(), '', null, '系统是否列表', null);
insert into sys_dict_type
values (7, '通知类型', 'sys_notice_type', '0', 'admin', sysdate(), '', null, '通知类型列表', null);
insert into sys_dict_type
values (8, '通知状态', 'sys_notice_status', '0', 'admin', sysdate(), '', null, '通知状态列表', null);
insert into sys_dict_type
values (9, '操作类型', 'sys_oper_type', '0', 'admin', sysdate(), '', null, '操作类型列表', null);
insert into sys_dict_type
values (10, '系统状态', 'sys_common_status', '0', 'admin', sysdate(), '', null, '登录状态列表', null);


insert into sys_dict_data
values (1, 1, '男', '0', 'sys_user_sex', '', '', 'Y', '0', 'admin', sysdate(), '', null, '性别男', null);
insert into sys_dict_data
values (2, 2, '女', '1', 'sys_user_sex', '', '', 'N', '0', 'admin', sysdate(), '', null, '性别女', null);
insert into sys_dict_data
values (3, 3, '未知', '2', 'sys_user_sex', '', '', 'N', '0', 'admin', sysdate(), '', null, '性别未知', null);
insert into sys_dict_data
values (4, 1, '显示', '0', 'sys_show_hide', '', 'primary', 'Y', '0', 'admin', sysdate(), '', null, '显示菜单', null);
insert into sys_dict_data
values (5, 2, '隐藏', '1', 'sys_show_hide', '', 'danger', 'N', '0', 'admin', sysdate(), '', null, '隐藏菜单', null);
insert into sys_dict_data
values (6, 1, '正常', '0', 'sys_normal_disable', '', 'primary', 'Y', '0', 'admin', sysdate(), '', null, '正常状态',
        null);
insert into sys_dict_data
values (7, 2, '停用', '1', 'sys_normal_disable', '', 'danger', 'N', '0', 'admin', sysdate(), '', null, '停用状态',
        null);
insert into sys_dict_data
values (8, 1, '正常', '0', 'sys_job_status', '', 'primary', 'Y', '0', 'admin', sysdate(), '', null, '正常状态', null);
insert into sys_dict_data
values (9, 2, '暂停', '1', 'sys_job_status', '', 'danger', 'N', '0', 'admin', sysdate(), '', null, '停用状态', null);
insert into sys_dict_data
values (10, 1, '默认', 'DEFAULT', 'sys_job_group', '', '', 'Y', '0', 'admin', sysdate(), '', null, '默认分组', null);
insert into sys_dict_data
values (11, 2, '系统', 'SYSTEM', 'sys_job_group', '', '', 'N', '0', 'admin', sysdate(), '', null, '系统分组', null);
insert into sys_dict_data
values (12, 1, '是', 'Y', 'sys_yes_no', '', 'primary', 'Y', '0', 'admin', sysdate(), '', null, '系统默认是', null);
insert into sys_dict_data
values (13, 2, '否', 'N', 'sys_yes_no', '', 'danger', 'N', '0', 'admin', sysdate(), '', null, '系统默认否', null);
insert into sys_dict_data
values (14, 1, '通知', '1', 'sys_notice_type', '', 'warning', 'Y', '0', 'admin', sysdate(), '', null, '通知', null);
insert into sys_dict_data
values (15, 2, '公告', '2', 'sys_notice_type', '', 'success', 'N', '0', 'admin', sysdate(), '', null, '公告', null);
insert into sys_dict_data
values (16, 1, '正常', '0', 'sys_notice_status', '', 'primary', 'Y', '0', 'admin', sysdate(), '', null, '正常状态',
        null);
insert into sys_dict_data
values (17, 2, '关闭', '1', 'sys_notice_status', '', 'danger', 'N', '0', 'admin', sysdate(), '', null, '关闭状态',
        null);
insert into sys_dict_data
values (18, 99, '其他', '0', 'sys_oper_type', '', 'info', 'N', '0', 'admin', sysdate(), '', null, '其他操作', null);
insert into sys_dict_data
values (19, 1, '新增', '1', 'sys_oper_type', '', 'info', 'N', '0', 'admin', sysdate(), '', null, '新增操作', null);
insert into sys_dict_data
values (20, 2, '修改', '2', 'sys_oper_type', '', 'info', 'N', '0', 'admin', sysdate(), '', null, '修改操作', null);
insert into sys_dict_data
values (21, 3, '删除', '3', 'sys_oper_type', '', 'danger', 'N', '0', 'admin', sysdate(), '', null, '删除操作', null);
insert into sys_dict_data
values (22, 4, '授权', '4', 'sys_oper_type', '', 'primary', 'N', '0', 'admin', sysdate(), '', null, '授权操作', null);
insert into sys_dict_data
values (23, 5, '导出', '5', 'sys_oper_type', '', 'warning', 'N', '0', 'admin', sysdate(), '', null, '导出操作', null);
insert into sys_dict_data
values (24, 6, '导入', '6', 'sys_oper_type', '', 'warning', 'N', '0', 'admin', sysdate(), '', null, '导入操作', null);
insert into sys_dict_data
values (25, 7, '强退', '7', 'sys_oper_type', '', 'danger', 'N', '0', 'admin', sysdate(), '', null, '强退操作', null);
insert into sys_dict_data
values (26, 8, '生成代码', '8', 'sys_oper_type', '', 'warning', 'N', '0', 'admin', sysdate(), '', null, '生成操作',
        null);
insert into sys_dict_data
values (27, 9, '清空数据', '9', 'sys_oper_type', '', 'danger', 'N', '0', 'admin', sysdate(), '', null, '清空操作',
        null);
insert into sys_dict_data
values (28, 1, '成功', '0', 'sys_common_status', '', 'primary', 'N', '0', 'admin', sysdate(), '', null, '正常状态',
        null);
insert into sys_dict_data
values (29, 2, '失败', '1', 'sys_common_status', '', 'danger', 'N', '0', 'admin', sysdate(), '', null, '停用状态',
        null);


-- ----------------------------
-- 初始化-参数配置
-- ----------------------------

insert into sys_config
values (1, '主框架页-默认皮肤样式名称', 'sys.index.skinName', 'skin-blue', 'Y', 'admin', sysdate(), '', null,
        '蓝色 skin-blue、绿色 skin-green、紫色 skin-purple、红色 skin-red、黄色 skin-yellow', null);
insert into sys_config
values (2, '用户管理-账号初始密码', 'sys.user.initPassword', '123456', 'Y', 'admin', sysdate(), '', null,
        '初始化密码 123456', null);
insert into sys_config
values (3, '主框架页-侧边栏主题', 'sys.index.sideTheme', 'theme-dark', 'Y', 'admin', sysdate(), '', null,
        '深色主题theme-dark，浅色主题theme-light', null);
insert into sys_config
values (4, '账号自助-是否开启用户注册功能', 'sys.account.registerUser', 'false', 'Y', 'admin', sysdate(), '', null,
        '是否开启注册用户功能（true开启，false关闭）', null);
insert into sys_config
values (5, '用户登录-黑名单列表', 'sys.login.blackIPList', '', 'Y', 'admin', sysdate(), '', null,
        '设置登录IP黑名单限制，多个匹配项以;分隔，支持匹配（*通配、网段）', null);

-- ----------------------------
-- 初始化-定时任务
-- ----------------------------

insert into sys_job
values (1, '系统默认（无参）', 'DEFAULT', 'TestFunc()', '0/1 * * * * ?', '3', '1', '1', 'admin', sysdate(), '',
        null, '', null);
insert into sys_job
values (2, '系统默认（多参）', 'DEFAULT', 'TestFuncWithParams(\'sophie\', true, 2000L, 316.50D, 100)',
        '0/20 * * * * ?', '3', '1', '1', 'admin', sysdate(), '', null, '', null);


