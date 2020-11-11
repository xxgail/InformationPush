/*
Navicat MySQL Data Transfer

Source Server         : local
Source Server Version : 50553
Source Host           : localhost:3306
Source Database       : device

Target Server Type    : MYSQL
Target Server Version : 50553
File Encoding         : 65001

Date: 2020-11-11 17:17:31
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for device
-- ----------------------------
DROP TABLE IF EXISTS `device`;
CREATE TABLE `device` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `channel` varchar(255) DEFAULT NULL,
  `device_token` varchar(255) DEFAULT NULL,
  `app_id` varchar(255) DEFAULT NULL,
  `group_id` int(11) DEFAULT NULL,
  `ip` varchar(255) DEFAULT NULL,
  `uid` int(11) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=16 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for platform
-- ----------------------------
DROP TABLE IF EXISTS `platform`;
CREATE TABLE `platform` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `app_id` varchar(255) DEFAULT NULL,
  `group_id` int(11) DEFAULT NULL,
  `channel` enum('oppo','vivo','mz','ios','hw','mi') DEFAULT 'mi',
  `value` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for platform_param
-- ----------------------------
DROP TABLE IF EXISTS `platform_param`;
CREATE TABLE `platform_param` (
  `app_id` varchar(255) NOT NULL,
  `hw_appId` varchar(255) DEFAULT NULL,
  `hw_clientSecret` varchar(255) DEFAULT NULL,
  `iOS_keyId` varchar(255) DEFAULT NULL,
  `iOS_teamId` varchar(255) DEFAULT NULL,
  `iOS_bundleId` varchar(255) DEFAULT NULL,
  `iOS_authTokenPath` varchar(255) DEFAULT NULL,
  `iOS_authToken` varchar(255) DEFAULT NULL,
  `mi_appSecret` varchar(255) DEFAULT NULL,
  `mi_restrictedPackageName` varchar(255) DEFAULT NULL,
  `mz_appId` varchar(255) DEFAULT NULL,
  `mz_appSecret` varchar(255) DEFAULT NULL,
  `oppo_appKey` varchar(255) DEFAULT NULL,
  `oppo_masterSecret` varchar(255) DEFAULT NULL,
  `vi_appId` varchar(255) DEFAULT NULL,
  `vi_appKey` varchar(255) DEFAULT NULL,
  `vi_appSecret` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`app_id`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
