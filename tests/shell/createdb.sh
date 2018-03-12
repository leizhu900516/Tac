




CREATE TABLE `backgroundtask` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `pid` int(11) DEFAULT NULL COMMENT '任务id',
  `taskname` varchar(100) DEFAULT NULL COMMENT '任务名称',
  `status` tinyint(4) DEFAULT NULL COMMENT '状态1正常0失败',
  `author` varchar(100) DEFAULT NULL COMMENT '作者',
  `ipaddress` varchar(15) DEFAULT NULL COMMENT '主机ip',
  `addtimes` int(11) DEFAULT NULL COMMENT '添加时间',
  `url` varchar(200) DEFAULT NULL COMMENT 'svn地址',
  `svnuser` varchar(100) DEFAULT NULL COMMENT 'svn账户',
  `svnpasswd` varchar(100) DEFAULT NULL COMMENT 'svn密码',
  `svn_number` varchar(20) DEFAULT NULL COMMENT 'svn版本号',
  `action_cmd` varchar(200) DEFAULT NULL COMMENT '执行命令',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8