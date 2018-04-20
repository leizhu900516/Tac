/*
SQLyog  v12.2.6 (64 bit)
MySQL - 5.7.20-log : Database - tac
*********************************************************************
*/


/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
CREATE DATABASE /*!32312 IF NOT EXISTS*/`tac` /*!40100 DEFAULT CHARACTER SET utf8 */;

USE `tac`;

/*Table structure for table `backgroundtask` */

DROP TABLE IF EXISTS `backgroundtask`;

CREATE TABLE `backgroundtask` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `pid` int(11) DEFAULT NULL COMMENT '任务id',
  `taskname` varchar(100) DEFAULT NULL COMMENT '任务名称',
  `status` tinyint(4) DEFAULT '0' COMMENT '状态1正常0失败',
  `author` varchar(100) DEFAULT NULL COMMENT '作者',
  `ipaddress` varchar(15) DEFAULT NULL COMMENT '主机ip',
  `addtimes` int(11) DEFAULT NULL COMMENT '添加时间',
  `url` varchar(200) DEFAULT NULL COMMENT 'svn地址',
  `svnuser` varchar(100) DEFAULT NULL COMMENT 'svn账户',
  `svnpasswd` varchar(100) DEFAULT NULL COMMENT 'svn密码',
  `svn_number` varchar(20) DEFAULT NULL COMMENT 'svn版本号',
  `action_cmd` varchar(200) DEFAULT NULL COMMENT '执行命令',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8;

/*Data for the table `backgroundtask` */

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;


CREATE TABLE `userinfo` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(100) DEFAULT NULL COMMENT '用户名',
  `password` varchar(100) DEFAULT NULL COMMENT '密码',
  `addtimes` int(10) DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='用户表'

insert into `userinfo` (`id`, `username`, `password`, `addtimes`) values('1','root','123456',NULL);
