-- ----------------------------
-- 创建各个表单的初始化数据
-- ----------------------------

-- ----------------------------
-- 初始化-部门表数据
-- ----------------------------
insert into sys_dept values(100, '', '0', 'admin', sysdate(), '', null, 0,   '0',         '雪菲科技',   0, '雪菲', '15888888888', 'xf@qq.com', '0');
insert into sys_dept values(101, '', '0', 'admin', sysdate(), '', null,100, '0,100',      '深圳总公司', 1, '雪菲', '15888888888', 'xf@qq.com', '0');
insert into sys_dept values(102, '', '0', 'admin', sysdate(), '', null,100, '0,100',      '长沙分公司', 2, '雪菲', '15888888888', 'xf@qq.com', '0');
insert into sys_dept values(103, '', '0', 'admin', sysdate(), '', null,101, '0,100,101',  '研发部门',   1, '雪菲', '15888888888', 'xf@qq.com', '0');
insert into sys_dept values(104, '', '0', 'admin', sysdate(), '', null,101, '0,100,101',  '市场部门',   2, '雪菲', '15888888888', 'xf@qq.com', '0');
insert into sys_dept values(105, '', '0', 'admin', sysdate(), '', null,101, '0,100,101',  '测试部门',   3, '雪菲', '15888888888', 'xf@qq.com', '0');
insert into sys_dept values(106, '', '0', 'admin', sysdate(), '', null,101, '0,100,101',  '财务部门',   4, '雪菲', '15888888888', 'xf@qq.com', '0');
insert into sys_dept values(107, '', '0', 'admin', sysdate(), '', null,101, '0,100,101',  '运维部门',   5, '雪菲', '15888888888', 'xf@qq.com', '0');
insert into sys_dept values(108, '', '0', 'admin', sysdate(), '', null,102, '0,100,102',  '市场部门',   1, '雪菲', '15888888888', 'xf@qq.com', '0');
insert into sys_dept values(109, '', '0', 'admin', sysdate(), '', null,102, '0,100,102',  '财务部门',   2, '雪菲', '15888888888', 'xf@qq.com', '0');

-- ----------------------------
-- 初始化-用户信息表数据
-- ----------------------------
insert into sys_user values(1, '', '0', 'admin', sysdate(), '', null, 103, 'admin', '雪菲', '00', 'xf@163.com', '15888888888', '1', '', '$2a$10$7JB720yubVSZvUI0rEqK/.VqGOZTH.ulu33dHOiBE8ByOhJIrdAu2', '0', '127.0.0.1', sysdate(), '管理员');
insert into sys_user values(2, '', '0', 'admin', sysdate(), '', null, 105, 'xf',    '雪菲', '00', 'xf@qq.com',  '15666666666', '1', '', '$2a$10$7JB720yubVSZvUI0rEqK/.VqGOZTH.ulu33dHOiBE8ByOhJIrdAu2', '0', '127.0.0.1', sysdate(), '测试员');

-- ----------------------------
-- 初始化-岗位信息表数据
-- ----------------------------
insert into sys_post values(1, '', '0', 'admin', sysdate(), '', null,'ceo',  '董事长',    1,  '');
insert into sys_post values(2, '', '0', 'admin', sysdate(), '', null,'se',   '项目经理',  2,  '');
insert into sys_post values(3, '', '0', 'admin', sysdate(), '', null,'hr',   '人力资源',  3,  '');
insert into sys_post values(4, '', '0', 'admin', sysdate(), '', null,'user', '普通员工',  4,  '');

-- ----------------------------
-- 初始化-角色信息表数据
-- ----------------------------
insert into sys_role values('1', '', '0', 'admin', sysdate(), '', null,'超级管理员',  'admin',  1, 1, 1, 1, '0', '超级管理员');
insert into sys_role values('2', '', '0', 'admin', sysdate(), '', null,'普通角色',    'common', 2, 2, 1, 1, '0', '普通角色');

-- ----------------------------
-- 初始化-菜单信息表数据
-- ----------------------------
-- 一级菜单
insert into sys_menu values('1', '', '0', 'admin', sysdate(), '', null, '系统管理', '0', '1', 'system',           null, '', 1, 0, 'M', '', 'system',  '系统管理目录');
insert into sys_menu values('2', '', '0', 'admin', sysdate(), '', null, '系统监控', '0', '2', 'monitor',          null, '', 1, 0, 'M', '', 'monitor', '系统监控目录');
insert into sys_menu values('3', '', '0', 'admin', sysdate(), '', null, '系统工具', '0', '3', 'tool',             null, '', 1, 0, 'M', '', 'tool',    '系统工具目录');
insert into sys_menu values('4', '', '0', 'admin', sysdate(), '', null, '雪菲官网', '0', '4', 'http://ruoyi.vip', null, '', 0, 0, 'M', '', 'guide',   '雪菲官网地址');
-- 二级菜单
insert into sys_menu values('100',  '', '0', 'admin', sysdate(), '', null, '用户管理',       '1',   '1', 'user',       'system/user/index',                 '', 1, 0, 'C', '0', 'system:user:list',        'user',          '用户管理菜单');
insert into sys_menu values('101',  '', '0', 'admin', sysdate(), '', null, '角色管理',       '1',   '2', 'role',       'system/role/index',                 '', 1, 0, 'C', '0', 'system:role:list',        'peoples',       '角色管理菜单');
insert into sys_menu values('102',  '', '0', 'admin', sysdate(), '', null, '菜单管理',       '1',   '3', 'menu',       'system/menu/index',                 '', 1, 0, 'C', '0', 'system:menu:list',        'tree-table',    '菜单管理菜单');
insert into sys_menu values('103',  '', '0', 'admin', sysdate(), '', null, '部门管理',       '1',   '4', 'dept',       'system/dept/index',                 '', 1, 0, 'C', '0', 'system:dept:list',        'tree',          '部门管理菜单');
insert into sys_menu values('104',  '', '0', 'admin', sysdate(), '', null, '岗位管理',       '1',   '5', 'post',       'system/post/index',                 '', 1, 0, 'C', '0', 'system:post:list',        'post',          '岗位管理菜单');
insert into sys_menu values('105',  '', '0', 'admin', sysdate(), '', null, '字典管理',       '1',   '6', 'dict',       'system/dict/index',                 '', 1, 0, 'C', '0', 'system:dict:list',        'dict',          '字典管理菜单');
insert into sys_menu values('106',  '', '0', 'admin', sysdate(), '', null, '参数设置',       '1',   '7', 'config',     'system/config/index',               '', 1, 0, 'C', '0', 'system:config:list',      'edit',          '参数设置菜单');
insert into sys_menu values('107',  '', '0', 'admin', sysdate(), '', null, '通知公告',       '1',   '8', 'notice',     'system/notice/index',               '', 1, 0, 'C', '0', 'system:notice:list',      'message',       '通知公告菜单');
insert into sys_menu values('108',  '', '0', 'admin', sysdate(), '', null, '日志管理',       '1',   '9', 'log',        '',                                  '', 1, 0, 'M', '0', '',                        'log',           '日志管理菜单');
insert into sys_menu values('109',  '', '0', 'admin', sysdate(), '', null, '在线用户',       '2',   '1', 'online',     'monitor/online/index',              '', 1, 0, 'C', '0', 'monitor:online:list',     'online',        '在线用户菜单');
insert into sys_menu values('110',  '', '0', 'admin', sysdate(), '', null, '定时任务',       '2',   '2', 'job',        'monitor/job/index',                 '', 1, 0, 'C', '0', 'monitor:job:list',        'job',           '定时任务菜单');
insert into sys_menu values('111',  '', '0', 'admin', sysdate(), '', null, 'Sentinel控制台', '2',   '3', 'http://localhost:8718',        '',                '', 0, 0, 'C', '0', 'monitor:sentinel:list',   'sentinel',      '流量控制菜单');
insert into sys_menu values('112',  '', '0', 'admin', sysdate(), '', null, 'Nacos控制台',    '2',   '4', 'http://localhost:8848/nacos',  '',                '', 0, 0, 'C', '0', 'monitor:nacos:list',      'nacos',         '服务治理菜单');
insert into sys_menu values('113',  '', '0', 'admin', sysdate(), '', null, 'Admin控制台',    '2',   '5', 'http://localhost:9100/login',  '',                '', 0, 0, 'C', '0', 'monitor:server:list',     'server',        '服务监控菜单');
insert into sys_menu values('114',  '', '0', 'admin', sysdate(), '', null, '表单构建',       '3',   '1', 'build',      'tool/build/index',                  '', 1, 0, 'C', '0', 'tool:build:list',         'build',         '表单构建菜单');
insert into sys_menu values('115',  '', '0', 'admin', sysdate(), '', null, '代码生成',       '3',   '2', 'gen',        'tool/gen/index',                    '', 1, 0, 'C', '0', 'tool:gen:list',           'code',          '代码生成菜单');
insert into sys_menu values('116',  '', '0', 'admin', sysdate(), '', null, '系统接口',       '3',   '3', 'http://localhost:8080/swagger-ui/index.html', '', '', 0, 0, 'C', '0', 'tool:swagger:list',       'swagger',       '系统接口菜单');
-- 三级菜单
insert into sys_menu values('500',  '', '0', 'admin', sysdate(), '', null, '操作日志', '108', '1', 'operlog',    'system/operlog/index',    '', 1, 0, 'C', '0', 'system:operlog:list',    'form',          '操作日志菜单');
insert into sys_menu values('501',  '', '0', 'admin', sysdate(), '', null, '登录日志', '108', '2', 'logininfor', 'system/logininfor/index', '', 1, 0, 'C', '0', 'system:logininfor:list', 'logininfor',    '登录日志菜单');
-- 用户管理按钮
insert into sys_menu values('1000', '', '0', 'admin', sysdate(), '', null, '用户查询', '100', '1',  '', '', '', 1, 0, 'F', '0', 'system:user:quexf',          '#', 'admin', '');
insert into sys_menu values('1001', '', '0', 'admin', sysdate(), '', null, '用户新增', '100', '2',  '', '', '', 1, 0, 'F', '0', 'system:user:add',            '#', 'admin', '');
insert into sys_menu values('1002', '', '0', 'admin', sysdate(), '', null, '用户修改', '100', '3',  '', '', '', 1, 0, 'F', '0', 'system:user:edit',           '#', 'admin', '');
insert into sys_menu values('1003', '', '0', 'admin', sysdate(), '', null, '用户删除', '100', '4',  '', '', '', 1, 0, 'F', '0', 'system:user:remove',         '#', 'admin', '');
insert into sys_menu values('1004', '', '0', 'admin', sysdate(), '', null, '用户导出', '100', '5',  '', '', '', 1, 0, 'F', '0', 'system:user:export',         '#', 'admin', '');
insert into sys_menu values('1005', '', '0', 'admin', sysdate(), '', null, '用户导入', '100', '6',  '', '', '', 1, 0, 'F', '0', 'system:user:import',         '#', 'admin', '');
insert into sys_menu values('1006', '', '0', 'admin', sysdate(), '', null, '重置密码', '100', '7',  '', '', '', 1, 0, 'F', '0', 'system:user:resetPwd',       '#', 'admin', '');
-- 角色管理按钮
insert into sys_menu values('1007', '', '0', 'admin', sysdate(), '', null, '角色查询', '101', '1',  '', '', '', 1, 0, 'F', '0', 'system:role:quexf',          '#', 'admin', '');
insert into sys_menu values('1008', '', '0', 'admin', sysdate(), '', null, '角色新增', '101', '2',  '', '', '', 1, 0, 'F', '0', 'system:role:add',            '#', 'admin', '');
insert into sys_menu values('1009', '', '0', 'admin', sysdate(), '', null, '角色修改', '101', '3',  '', '', '', 1, 0, 'F', '0', 'system:role:edit',           '#', 'admin', '');
insert into sys_menu values('1010', '', '0', 'admin', sysdate(), '', null, '角色删除', '101', '4',  '', '', '', 1, 0, 'F', '0', 'system:role:remove',         '#', 'admin', '');
insert into sys_menu values('1011', '', '0', 'admin', sysdate(), '', null, '角色导出', '101', '5',  '', '', '', 1, 0, 'F', '0', 'system:role:export',         '#', 'admin', '');
-- 菜单管理按钮
insert into sys_menu values('1012', '', '0', 'admin', sysdate(), '', null, '菜单查询', '102', '1',  '', '', '', 1, 0, 'F', '0', 'system:menu:quexf',          '#', '');
insert into sys_menu values('1013', '', '0', 'admin', sysdate(), '', null, '菜单新增', '102', '2',  '', '', '', 1, 0, 'F', '0', 'system:menu:add',            '#', '');
insert into sys_menu values('1014', '', '0', 'admin', sysdate(), '', null, '菜单修改', '102', '3',  '', '', '', 1, 0, 'F', '0', 'system:menu:edit',           '#', '');
insert into sys_menu values('1015', '', '0', 'admin', sysdate(), '', null, '菜单删除', '102', '4',  '', '', '', 1, 0, 'F', '0', 'system:menu:remove',         '#', '');
-- 部门管理按钮
insert into sys_menu values('1016', '', '0', 'admin', sysdate(), '', null, '部门查询', '103', '1',  '', '', '', 1, 0, 'F', '0', 'system:dept:quexf',          '#', 'admin', '');
insert into sys_menu values('1017', '', '0', 'admin', sysdate(), '', null, '部门新增', '103', '2',  '', '', '', 1, 0, 'F', '0', 'system:dept:add',            '#', 'admin', '');
insert into sys_menu values('1018', '', '0', 'admin', sysdate(), '', null, '部门修改', '103', '3',  '', '', '', 1, 0, 'F', '0', 'system:dept:edit',           '#', 'admin', '');
insert into sys_menu values('1019', '', '0', 'admin', sysdate(), '', null, '部门删除', '103', '4',  '', '', '', 1, 0, 'F', '0', 'system:dept:remove',         '#', 'admin', '');
-- 岗位管理按钮
insert into sys_menu values('1020', '', '0', 'admin', sysdate(), '', null, '岗位查询', '104', '1',  '', '', '', 1, 0, 'F', '0', 'system:post:quexf',          '#', '');
insert into sys_menu values('1021', '', '0', 'admin', sysdate(), '', null, '岗位新增', '104', '2',  '', '', '', 1, 0, 'F', '0', 'system:post:add',            '#', '');
insert into sys_menu values('1022', '', '0', 'admin', sysdate(), '', null, '岗位修改', '104', '3',  '', '', '', 1, 0, 'F', '0', 'system:post:edit',           '#', '');
insert into sys_menu values('1023', '', '0', 'admin', sysdate(), '', null, '岗位删除', '104', '4',  '', '', '', 1, 0, 'F', '0', 'system:post:remove',         '#', '');
insert into sys_menu values('1024', '', '0', 'admin', sysdate(), '', null, '岗位导出', '104', '5',  '', '', '', 1, 0, 'F', '0', 'system:post:export',         '#', '');
-- 字典管理按钮
insert into sys_menu values('1025', '', '0', 'admin', sysdate(), '', null, '字典查询', '105', '1', '#', '', '', 1, 0, 'F', '0', 'system:dict:quexf',          '#', '');
insert into sys_menu values('1026', '', '0', 'admin', sysdate(), '', null, '字典新增', '105', '2', '#', '', '', 1, 0, 'F', '0', 'system:dict:add',            '#', '');
insert into sys_menu values('1027', '', '0', 'admin', sysdate(), '', null, '字典修改', '105', '3', '#', '', '', 1, 0, 'F', '0', 'system:dict:edit',           '#', '');
insert into sys_menu values('1028', '', '0', 'admin', sysdate(), '', null, '字典删除', '105', '4', '#', '', '', 1, 0, 'F', '0', 'system:dict:remove',         '#', '');
insert into sys_menu values('1029', '', '0', 'admin', sysdate(), '', null, '字典导出', '105', '5', '#', '', '', 1, 0, 'F', '0', 'system:dict:export',         '#', '');
-- 参数设置按钮
insert into sys_menu values('1030', '', '0', 'admin', sysdate(), '', null, '参数查询', '106', '1', '#', '', '', 1, 0, 'F', '0', 'system:config:quexf',        '#', '');
insert into sys_menu values('1031', '', '0', 'admin', sysdate(), '', null, '参数新增', '106', '2', '#', '', '', 1, 0, 'F', '0', 'system:config:add',          '#', '');
insert into sys_menu values('1032', '', '0', 'admin', sysdate(), '', null, '参数修改', '106', '3', '#', '', '', 1, 0, 'F', '0', 'system:config:edit',         '#', '');
insert into sys_menu values('1033', '', '0', 'admin', sysdate(), '', null, '参数删除', '106', '4', '#', '', '', 1, 0, 'F', '0', 'system:config:remove',       '#', '');
insert into sys_menu values('1034', '', '0', 'admin', sysdate(), '', null, '参数导出', '106', '5', '#', '', '', 1, 0, 'F', '0', 'system:config:export',       '#', '');
-- 通知公告按钮
insert into sys_menu values('1035', '', '0', 'admin', sysdate(), '', null, '公告查询', '107', '1', '#', '', '', 1, 0, 'F', '0', 'system:notice:quexf',        '#', 'admin', '');
insert into sys_menu values('1036', '', '0', 'admin', sysdate(), '', null, '公告新增', '107', '2', '#', '', '', 1, 0, 'F', '0', 'system:notice:add',          '#', 'admin', '');
insert into sys_menu values('1037', '', '0', 'admin', sysdate(), '', null, '公告修改', '107', '3', '#', '', '', 1, 0, 'F', '0', 'system:notice:edit',         '#', 'admin', '');
insert into sys_menu values('1038', '', '0', 'admin', sysdate(), '', null, '公告删除', '107', '4', '#', '', '', 1, 0, 'F', '0', 'system:notice:remove',       '#', 'admin', '');
-- 操作日志按钮
insert into sys_menu values('1039', '', '0', 'admin', sysdate(), '', null, '操作查询', '500', '1', '#', '', '', 1, 0, 'F', '0', 'system:operlog:quexf',       '#', '');
insert into sys_menu values('1040', '', '0', 'admin', sysdate(), '', null, '操作删除', '500', '2', '#', '', '', 1, 0, 'F', '0', 'system:operlog:remove',      '#', '');
insert into sys_menu values('1041', '', '0', 'admin', sysdate(), '', null, '日志导出', '500', '3', '#', '', '', 1, 0, 'F', '0', 'system:operlog:export',      '#', '');
-- 登录日志按钮
insert into sys_menu values('1042', '', '0', 'admin', sysdate(), '', null, '登录查询', '501', '1', '#', '', '', 1, 0, 'F', '0', 'system:logininfor:quexf',    '#', 'admin', '');
insert into sys_menu values('1043', '', '0', 'admin', sysdate(), '', null, '登录删除', '501', '2', '#', '', '', 1, 0, 'F', '0', 'system:logininfor:remove',   '#', 'admin', '');
insert into sys_menu values('1044', '', '0', 'admin', sysdate(), '', null, '日志导出', '501', '3', '#', '', '', 1, 0, 'F', '0', 'system:logininfor:export',   '#', 'admin', '');
insert into sys_menu values('1045', '', '0', 'admin', sysdate(), '', null, '账户解锁', '501', '4', '#', '', '', 1, 0, 'F', '0', 'system:logininfor:unlock',   '#', 'admin', '');
-- 在线用户按钮
insert into sys_menu values('1046', '', '0', 'admin', sysdate(), '', null, '在线查询', '109', '1', '#', '', '', 1, 0, 'F', '0', 'monitor:online:quexf',       '#', '');
insert into sys_menu values('1047', '', '0', 'admin', sysdate(), '', null, '批量强退', '109', '2', '#', '', '', 1, 0, 'F', '0', 'monitor:online:batchLogout', '#', '');
insert into sys_menu values('1048', '', '0', 'admin', sysdate(), '', null, '单条强退', '109', '3', '#', '', '', 1, 0, 'F', '0', 'monitor:online:forceLogout', '#', '');
-- 定时任务按钮
insert into sys_menu values('1049', '', '0', 'admin', sysdate(), '', null, '任务查询', '110', '1', '#', '', '', 1, 0, 'F', '0', 'monitor:job:quexf',          '#', '');
insert into sys_menu values('1050', '', '0', 'admin', sysdate(), '', null, '任务新增', '110', '2', '#', '', '', 1, 0, 'F', '0', 'monitor:job:add',            '#', '');
insert into sys_menu values('1051', '', '0', 'admin', sysdate(), '', null, '任务修改', '110', '3', '#', '', '', 1, 0, 'F', '0', 'monitor:job:edit',           '#', '');
insert into sys_menu values('1052', '', '0', 'admin', sysdate(), '', null, '任务删除', '110', '4', '#', '', '', 1, 0, 'F', '0', 'monitor:job:remove',         '#', '');
insert into sys_menu values('1053', '', '0', 'admin', sysdate(), '', null, '状态修改', '110', '5', '#', '', '', 1, 0, 'F', '0', 'monitor:job:changeStatus',   '#', '');
insert into sys_menu values('1054', '', '0', 'admin', sysdate(), '', null, '任务导出', '110', '6', '#', '', '', 1, 0, 'F', '0', 'monitor:job:export',         '#', '');
-- 代码生成按钮
insert into sys_menu values('1055', '', '0', 'admin', sysdate(), '', null, '生成查询', '115', '1', '#', '', '', 1, 0, 'F', '0', 'tool:gen:quexf',             '#', 'admin', '');
insert into sys_menu values('1056', '', '0', 'admin', sysdate(), '', null, '生成修改', '115', '2', '#', '', '', 1, 0, 'F', '0', 'tool:gen:edit',              '#', 'admin', '');
insert into sys_menu values('1057', '', '0', 'admin', sysdate(), '', null, '生成删除', '115', '3', '#', '', '', 1, 0, 'F', '0', 'tool:gen:remove',            '#', 'admin', '');
insert into sys_menu values('1058', '', '0', 'admin', sysdate(), '', null, '导入代码', '115', '2', '#', '', '', 1, 0, 'F', '0', 'tool:gen:import',            '#', 'admin', '');
insert into sys_menu values('1059', '', '0', 'admin', sysdate(), '', null, '预览代码', '115', '4', '#', '', '', 1, 0, 'F', '0', 'tool:gen:preview',           '#', 'admin', '');
insert into sys_menu values('1060', '', '0', 'admin', sysdate(), '', null, '生成代码', '115', '5', '#', '', '', 1, 0, 'F', '0', 'tool:gen:code',              '#', 'admin', '');

-- ----------------------------
-- 初始化-用户和角色关联表数据
-- ----------------------------
insert into sys_user_role values ('1', '1');
insert into sys_user_role values ('2', '2');

-- ----------------------------
-- 初始化-角色和菜单关联表数据
-- ----------------------------
insert into sys_role_menu values ('2', '1');
insert into sys_role_menu values ('2', '2');
insert into sys_role_menu values ('2', '3');
insert into sys_role_menu values ('2', '4');
insert into sys_role_menu values ('2', '100');
insert into sys_role_menu values ('2', '101');
insert into sys_role_menu values ('2', '102');
insert into sys_role_menu values ('2', '103');
insert into sys_role_menu values ('2', '104');
insert into sys_role_menu values ('2', '105');
insert into sys_role_menu values ('2', '106');
insert into sys_role_menu values ('2', '107');
insert into sys_role_menu values ('2', '108');
insert into sys_role_menu values ('2', '109');
insert into sys_role_menu values ('2', '110');
insert into sys_role_menu values ('2', '111');
insert into sys_role_menu values ('2', '112');
insert into sys_role_menu values ('2', '113');
insert into sys_role_menu values ('2', '114');
insert into sys_role_menu values ('2', '115');
insert into sys_role_menu values ('2', '116');
insert into sys_role_menu values ('2', '500');
insert into sys_role_menu values ('2', '501');
insert into sys_role_menu values ('2', '1000');
insert into sys_role_menu values ('2', '1001');
insert into sys_role_menu values ('2', '1002');
insert into sys_role_menu values ('2', '1003');
insert into sys_role_menu values ('2', '1004');
insert into sys_role_menu values ('2', '1005');
insert into sys_role_menu values ('2', '1006');
insert into sys_role_menu values ('2', '1007');
insert into sys_role_menu values ('2', '1008');
insert into sys_role_menu values ('2', '1009');
insert into sys_role_menu values ('2', '1010');
insert into sys_role_menu values ('2', '1011');
insert into sys_role_menu values ('2', '1012');
insert into sys_role_menu values ('2', '1013');
insert into sys_role_menu values ('2', '1014');
insert into sys_role_menu values ('2', '1015');
insert into sys_role_menu values ('2', '1016');
insert into sys_role_menu values ('2', '1017');
insert into sys_role_menu values ('2', '1018');
insert into sys_role_menu values ('2', '1019');
insert into sys_role_menu values ('2', '1020');
insert into sys_role_menu values ('2', '1021');
insert into sys_role_menu values ('2', '1022');
insert into sys_role_menu values ('2', '1023');
insert into sys_role_menu values ('2', '1024');
insert into sys_role_menu values ('2', '1025');
insert into sys_role_menu values ('2', '1026');
insert into sys_role_menu values ('2', '1027');
insert into sys_role_menu values ('2', '1028');
insert into sys_role_menu values ('2', '1029');
insert into sys_role_menu values ('2', '1030');
insert into sys_role_menu values ('2', '1031');
insert into sys_role_menu values ('2', '1032');
insert into sys_role_menu values ('2', '1033');
insert into sys_role_menu values ('2', '1034');
insert into sys_role_menu values ('2', '1035');
insert into sys_role_menu values ('2', '1036');
insert into sys_role_menu values ('2', '1037');
insert into sys_role_menu values ('2', '1038');
insert into sys_role_menu values ('2', '1039');
insert into sys_role_menu values ('2', '1040');
insert into sys_role_menu values ('2', '1041');
insert into sys_role_menu values ('2', '1042');
insert into sys_role_menu values ('2', '1043');
insert into sys_role_menu values ('2', '1044');
insert into sys_role_menu values ('2', '1045');
insert into sys_role_menu values ('2', '1046');
insert into sys_role_menu values ('2', '1047');
insert into sys_role_menu values ('2', '1048');
insert into sys_role_menu values ('2', '1049');
insert into sys_role_menu values ('2', '1050');
insert into sys_role_menu values ('2', '1051');
insert into sys_role_menu values ('2', '1052');
insert into sys_role_menu values ('2', '1053');
insert into sys_role_menu values ('2', '1054');
insert into sys_role_menu values ('2', '1055');
insert into sys_role_menu values ('2', '1056');
insert into sys_role_menu values ('2', '1057');
insert into sys_role_menu values ('2', '1058');
insert into sys_role_menu values ('2', '1059');
insert into sys_role_menu values ('2', '1060');

-- ----------------------------
-- 初始化-角色和部门关联表数据
-- ----------------------------
insert into sys_role_dept values ('2', '100');
insert into sys_role_dept values ('2', '101');
insert into sys_role_dept values ('2', '105');

-- ----------------------------
-- 初始化-用户与岗位关联表数据
-- ----------------------------
insert into sys_user_post values ('1', '1');
insert into sys_user_post values ('2', '2');


-- ----------------------------
-- 初始化-字典类型表数据
-- ----------------------------
insert into sys_dict_type values(1,  '', 'admin', sysdate(), '', null,'用户性别', 'sys_user_sex',        '用户性别列表');
insert into sys_dict_type values(2,  '', 'admin', sysdate(), '', null,'菜单状态', 'sys_show_hide',       '菜单状态列表');
insert into sys_dict_type values(3,  '', 'admin', sysdate(), '', null,'系统开关', 'sys_normal_disable',  '系统开关列表');
insert into sys_dict_type values(4,  '', 'admin', sysdate(), '', null,'任务状态', 'sys_job_status',      '任务状态列表');
insert into sys_dict_type values(5,  '', 'admin', sysdate(), '', null,'任务分组', 'sys_job_group',       '任务分组列表');
insert into sys_dict_type values(6,  '', 'admin', sysdate(), '', null,'系统是否', 'sys_yes_no',          '系统是否列表');
insert into sys_dict_type values(7,  '', 'admin', sysdate(), '', null,'通知类型', 'sys_notice_type',     '通知类型列表');
insert into sys_dict_type values(8,  '', 'admin', sysdate(), '', null,'通知状态', 'sys_notice_status',   '通知状态列表');
insert into sys_dict_type values(9,  '', 'admin', sysdate(), '', null,'操作类型', 'sys_oper_type',       '操作类型列表');
insert into sys_dict_type values(10, '', 'admin', sysdate(), '', null,'系统状态', 'sys_common_status',   '登录状态列表');


-- ----------------------------
-- 初始化-字典数据表数据
-- ----------------------------
insert into sys_dict_data values(1,  '', '0', 'admin', sysdate(), '', null, 1,  '男',       '0',       'sys_user_sex',        '',   '',        'Y', '性别男');
insert into sys_dict_data values(2,  '', '0', 'admin', sysdate(), '', null, 2,  '女',       '1',       'sys_user_sex',        '',   '',        'N', '性别女');
insert into sys_dict_data values(3,  '', '0', 'admin', sysdate(), '', null, 3,  '未知',     '2',       'sys_user_sex',        '',   '',        'N', '性别未知');
insert into sys_dict_data values(4,  '', '0', 'admin', sysdate(), '', null, 1,  '显示',     '0',       'sys_show_hide',       '',   'primaxf', 'Y', '显示菜单');
insert into sys_dict_data values(5,  '', '0', 'admin', sysdate(), '', null, 2,  '隐藏',     '1',       'sys_show_hide',       '',   'danger',  'N', '隐藏菜单');
insert into sys_dict_data values(6,  '', '0', 'admin', sysdate(), '', null, 1,  '正常',     '0',       'sys_normal_disable',  '',   'primaxf', 'Y', '正常状态');
insert into sys_dict_data values(7,  '', '0', 'admin', sysdate(), '', null, 2,  '停用',     '1',       'sys_normal_disable',  '',   'danger',  'N', '停用状态');
insert into sys_dict_data values(8,  '', '0', 'admin', sysdate(), '', null, 1,  '正常',     '0',       'sys_job_status',      '',   'primaxf', 'Y', '正常状态');
insert into sys_dict_data values(9,  '', '0', 'admin', sysdate(), '', null, 2,  '暂停',     '1',       'sys_job_status',      '',   'danger',  'N', '停用状态');
insert into sys_dict_data values(10, '', '0', 'admin', sysdate(), '', null, 1,  '默认',     'DEFAULT', 'sys_job_group',       '',   '',        'Y', '默认分组');
insert into sys_dict_data values(11, '', '0', 'admin', sysdate(), '', null, 2,  '系统',     'SYSTEM',  'sys_job_group',       '',   '',        'N', '系统分组');
insert into sys_dict_data values(12, '', '0', 'admin', sysdate(), '', null, 1,  '是',       'Y',       'sys_yes_no',          '',   'primaxf', 'Y', '系统默认是');
insert into sys_dict_data values(13, '', '0', 'admin', sysdate(), '', null, 2,  '否',       'N',       'sys_yes_no',          '',   'danger',  'N', '系统默认否');
insert into sys_dict_data values(14, '', '0', 'admin', sysdate(), '', null, 1,  '通知',     '1',       'sys_notice_type',     '',   'warning', 'Y', '通知');
insert into sys_dict_data values(15, '', '0', 'admin', sysdate(), '', null, 2,  '公告',     '2',       'sys_notice_type',     '',   'success', 'N', '公告');
insert into sys_dict_data values(16, '', '0', 'admin', sysdate(), '', null, 1,  '正常',     '0',       'sys_notice_status',   '',   'primaxf', 'Y', '正常状态');
insert into sys_dict_data values(17, '', '0', 'admin', sysdate(), '', null, 2,  '关闭',     '1',       'sys_notice_status',   '',   'danger',  'N', '关闭状态');
insert into sys_dict_data values(18, '', '0', 'admin', sysdate(), '', null, 99, '其他',     '0',       'sys_oper_type',       '',   'info',    'N', '其他操作');
insert into sys_dict_data values(19, '', '0', 'admin', sysdate(), '', null, 1,  '新增',     '1',       'sys_oper_type',       '',   'info',    'N', '新增操作');
insert into sys_dict_data values(20, '', '0', 'admin', sysdate(), '', null, 2,  '修改',     '2',       'sys_oper_type',       '',   'info',    'N', '修改操作');
insert into sys_dict_data values(21, '', '0', 'admin', sysdate(), '', null, 3,  '删除',     '3',       'sys_oper_type',       '',   'danger',  'N', '删除操作');
insert into sys_dict_data values(22, '', '0', 'admin', sysdate(), '', null, 4,  '授权',     '4',       'sys_oper_type',       '',   'primaxf', 'N', '授权操作');
insert into sys_dict_data values(23, '', '0', 'admin', sysdate(), '', null, 5,  '导出',     '5',       'sys_oper_type',       '',   'warning', 'N', '导出操作');
insert into sys_dict_data values(24, '', '0', 'admin', sysdate(), '', null, 6,  '导入',     '6',       'sys_oper_type',       '',   'warning', 'N', '导入操作');
insert into sys_dict_data values(25, '', '0', 'admin', sysdate(), '', null, 7,  '强退',     '7',       'sys_oper_type',       '',   'danger',  'N', '强退操作');
insert into sys_dict_data values(26, '', '0', 'admin', sysdate(), '', null, 8,  '生成代码', '8',       'sys_oper_type',       '',   'warning', 'N', '生成操作');
insert into sys_dict_data values(27, '', '0', 'admin', sysdate(), '', null, 9,  '清空数据', '9',       'sys_oper_type',       '',   'danger',  'N', '清空操作');
insert into sys_dict_data values(28, '', '0', 'admin', sysdate(), '', null, 1,  '成功',     '0',       'sys_common_status',   '',   'primaxf', 'N', '正常状态');
insert into sys_dict_data values(29, '', '0', 'admin', sysdate(), '', null, 2,  '失败',     '1',       'sys_common_status',   '',   'danger',  'N', '停用状态');


-- ----------------------------
-- 初始化-参数配置表
-- ----------------------------
insert into sys_config values(1, '', '0', 'admin', sysdate(), '', null, '主框架页-默认皮肤样式名称',     'sys.index.skinName',       'skin-blue',     'Y', '蓝色 skin-blue、绿色 skin-green、紫色 skin-purple、红色 skin-red、黄色 skin-yellow' );
insert into sys_config values(2, '', '0', 'admin', sysdate(), '', null, '用户管理-账号初始密码',         'sys.user.initPassword',    '123456',        'Y', '初始化密码 123456' );
insert into sys_config values(3, '', '0', 'admin', sysdate(), '', null, '主框架页-侧边栏主题',           'sys.index.sideTheme',      'theme-dark',    'Y', '深色主题theme-dark，浅色主题theme-light' );
insert into sys_config values(4, '', '0', 'admin', sysdate(), '', null, '账号自助-是否开启用户注册功能', 'sys.account.registerUser', 'false',         'Y', '是否开启注册用户功能（true开启，false关闭）');
insert into sys_config values(5, '', '0', 'admin', sysdate(), '', null, '用户登录-黑名单列表',           'sys.login.blackIPList',    '',              'Y', '设置登录IP黑名单限制，多个匹配项以;分隔，支持匹配（*通配、网段）');


-- ----------------------------
-- 初始化-定时任务调度表
-- ----------------------------
insert into sys_job values(1, '', '1', 'admin', sysdate(), '', null,  '系统默认（无参）', 'DEFAULT', 'xfTask.xfNoParams',        '0/10 * * * * ?', '3', '1', '');
insert into sys_job values(2, '', '1', 'admin', sysdate(), '', null,  '系统默认（有参）', 'DEFAULT', 'xfTask.xfParams(\'xf\')',  '0/15 * * * * ?', '3', '1', '');
insert into sys_job values(3, '', '1', 'admin', sysdate(), '', null,  '系统默认（多参）', 'DEFAULT', 'xfTask.xfMultipleParams(\'xf\', true, 2000L, 316.50D, 100)',  '0/20 * * * * ?', '3', '1', '');
