/*
 Navicat Premium Data Transfer

 Source Server         : local-2
 Source Server Type    : MySQL
 Source Server Version : 80032 (8.0.32)
 Source Host           : 192.168.99.68:3306
 Source Schema         : icheck

 Target Server Type    : MySQL
 Target Server Version : 80032 (8.0.32)
 File Encoding         : 65001

 Date: 06/06/2023 16:17:18
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for check_historys
-- ----------------------------
DROP TABLE IF EXISTS `check_historys`;
CREATE TABLE `check_historys` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `donode_id` int DEFAULT NULL,
  `node_id` int DEFAULT NULL,
  `ip` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ping_long` float DEFAULT NULL,
  `node_long` float DEFAULT NULL,
  `api_long` float DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `index_ip` (`ip`)
) ENGINE=InnoDB AUTO_INCREMENT=2033 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for check_nodes
-- ----------------------------
DROP TABLE IF EXISTS `check_nodes`;
CREATE TABLE `check_nodes` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `host` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `list_path` text COLLATE utf8mb4_unicode_ci,
  `limit_wait` int DEFAULT '20',
  `node_type` varchar(50) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `req_encode` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `req_encode_key` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `is_show` tinyint DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Table structure for docheck_nodes
-- ----------------------------
DROP TABLE IF EXISTS `docheck_nodes`;
CREATE TABLE `docheck_nodes` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `ip` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `city` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
