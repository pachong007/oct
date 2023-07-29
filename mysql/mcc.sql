/*
 Navicat Premium Data Transfer

 Source Server         : loc
 Source Server Type    : MySQL
 Source Server Version : 50737
 Source Host           : 127.0.0.1:3306
 Source Schema         : mcc

 Target Server Type    : MySQL
 Target Server Version : 50737
 File Encoding         : 65001

 Date: 29/07/2023 10:26:52
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for admin_menu
-- ----------------------------
DROP TABLE IF EXISTS `admin_menu`;
CREATE TABLE `admin_menu`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `parent_id` int(11) NOT NULL DEFAULT 0,
  `order` int(11) NOT NULL DEFAULT 0,
  `title` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `uri` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `permission` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 12 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_menu
-- ----------------------------
INSERT INTO `admin_menu` VALUES (1, 0, 1, 'Dashboard', 'fa-bar-chart', '/', NULL, NULL, NULL);
INSERT INTO `admin_menu` VALUES (2, 0, 2, 'Admin', 'fa-tasks', '', NULL, NULL, NULL);
INSERT INTO `admin_menu` VALUES (3, 2, 3, 'Users', 'fa-users', 'auth/users', NULL, NULL, NULL);
INSERT INTO `admin_menu` VALUES (4, 2, 4, 'Roles', 'fa-user', 'auth/roles', NULL, NULL, NULL);
INSERT INTO `admin_menu` VALUES (5, 2, 5, 'Permission', 'fa-ban', 'auth/permissions', NULL, NULL, NULL);
INSERT INTO `admin_menu` VALUES (6, 2, 6, 'Menu', 'fa-bars', 'auth/menu', NULL, NULL, NULL);
INSERT INTO `admin_menu` VALUES (7, 2, 7, 'Operation log', 'fa-history', 'auth/logs', NULL, NULL, NULL);
INSERT INTO `admin_menu` VALUES (8, 0, 0, '爬虫', 'fa-bug', '/paw', NULL, '2023-07-29 09:23:16', '2023-07-29 09:23:16');
INSERT INTO `admin_menu` VALUES (9, 0, 0, '漫画', 'fa-book', '/source_comic', NULL, '2023-07-29 09:23:35', '2023-07-29 09:23:35');
INSERT INTO `admin_menu` VALUES (10, 0, 0, '章节', 'fa-file', '/source_chapter', NULL, '2023-07-29 09:23:53', '2023-07-29 09:23:53');
INSERT INTO `admin_menu` VALUES (11, 0, 0, '错误日志', 'fa-warning', '/fail_info', NULL, '2023-07-29 09:24:20', '2023-07-29 09:24:20');

-- ----------------------------
-- Table structure for admin_operation_log
-- ----------------------------
DROP TABLE IF EXISTS `admin_operation_log`;
CREATE TABLE `admin_operation_log`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `method` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `ip` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `input` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `admin_operation_log_user_id_index`(`user_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 115 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_operation_log
-- ----------------------------
INSERT INTO `admin_operation_log` VALUES (1, 1, 'admin', 'GET', '172.27.0.1', '[]', '2023-07-29 09:22:52', '2023-07-29 09:22:52');
INSERT INTO `admin_operation_log` VALUES (2, 1, 'admin/auth/menu', 'GET', '172.27.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-07-29 09:22:57', '2023-07-29 09:22:57');
INSERT INTO `admin_operation_log` VALUES (3, 1, 'admin/auth/menu', 'POST', '172.27.0.1', '{\"parent_id\":\"0\",\"title\":\"\\u722c\\u866b\",\"icon\":\"fa-bug\",\"uri\":\"\\/paw\",\"roles\":[null],\"permission\":null,\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\"}', '2023-07-29 09:23:16', '2023-07-29 09:23:16');
INSERT INTO `admin_operation_log` VALUES (4, 1, 'admin/auth/menu', 'GET', '172.27.0.1', '[]', '2023-07-29 09:23:17', '2023-07-29 09:23:17');
INSERT INTO `admin_operation_log` VALUES (5, 1, 'admin/auth/menu', 'POST', '172.27.0.1', '{\"parent_id\":\"0\",\"title\":\"\\u6f2b\\u753b\",\"icon\":\"fa-book\",\"uri\":\"\\/source_comic\",\"roles\":[null],\"permission\":null,\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\"}', '2023-07-29 09:23:34', '2023-07-29 09:23:34');
INSERT INTO `admin_operation_log` VALUES (6, 1, 'admin/auth/menu', 'GET', '172.27.0.1', '[]', '2023-07-29 09:23:35', '2023-07-29 09:23:35');
INSERT INTO `admin_operation_log` VALUES (7, 1, 'admin/auth/menu', 'POST', '172.27.0.1', '{\"parent_id\":\"0\",\"title\":\"\\u7ae0\\u8282\",\"icon\":\"fa-file\",\"uri\":\"\\/source_chapter\",\"roles\":[null],\"permission\":null,\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\"}', '2023-07-29 09:23:52', '2023-07-29 09:23:52');
INSERT INTO `admin_operation_log` VALUES (8, 1, 'admin/auth/menu', 'GET', '172.27.0.1', '[]', '2023-07-29 09:23:53', '2023-07-29 09:23:53');
INSERT INTO `admin_operation_log` VALUES (9, 1, 'admin/auth/menu', 'POST', '172.27.0.1', '{\"parent_id\":\"0\",\"title\":\"\\u9519\\u8bef\\u65e5\\u5fd7\",\"icon\":\"fa-warning\",\"uri\":\"\\/fail_info\",\"roles\":[null],\"permission\":null,\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\"}', '2023-07-29 09:24:20', '2023-07-29 09:24:20');
INSERT INTO `admin_operation_log` VALUES (10, 1, 'admin/auth/menu', 'GET', '172.27.0.1', '[]', '2023-07-29 09:24:20', '2023-07-29 09:24:20');
INSERT INTO `admin_operation_log` VALUES (11, 1, 'admin/auth/menu', 'GET', '172.27.0.1', '[]', '2023-07-29 09:24:22', '2023-07-29 09:24:22');
INSERT INTO `admin_operation_log` VALUES (12, 1, 'admin/auth/menu', 'GET', '172.27.0.1', '[]', '2023-07-29 09:25:21', '2023-07-29 09:25:21');
INSERT INTO `admin_operation_log` VALUES (13, 1, 'admin/paw', 'GET', '172.27.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-07-29 09:25:24', '2023-07-29 09:25:24');
INSERT INTO `admin_operation_log` VALUES (14, 1, 'admin/paw.cache', 'POST', '172.27.0.1', '{\"_method\":\"post\",\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\",\"source\":\"kk\"}', '2023-07-29 09:25:25', '2023-07-29 09:25:25');
INSERT INTO `admin_operation_log` VALUES (15, 1, 'admin/source_comic', 'GET', '172.27.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-07-29 09:25:30', '2023-07-29 09:25:30');
INSERT INTO `admin_operation_log` VALUES (16, 1, 'admin/paw', 'GET', '172.27.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-07-29 09:25:37', '2023-07-29 09:25:37');
INSERT INTO `admin_operation_log` VALUES (17, 1, 'admin/paw.cache', 'POST', '172.27.0.1', '{\"_method\":\"post\",\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\",\"source\":\"kk\"}', '2023-07-29 09:25:38', '2023-07-29 09:25:38');
INSERT INTO `admin_operation_log` VALUES (18, 1, 'admin/paw', 'GET', '172.27.0.1', '[]', '2023-07-29 09:25:51', '2023-07-29 09:25:51');
INSERT INTO `admin_operation_log` VALUES (19, 1, 'admin/paw.cache', 'POST', '172.27.0.1', '{\"_method\":\"post\",\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\",\"source\":\"kk\"}', '2023-07-29 09:25:52', '2023-07-29 09:25:52');
INSERT INTO `admin_operation_log` VALUES (20, 1, 'admin/paw.cache', 'POST', '172.27.0.1', '{\"_method\":\"post\",\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\",\"source\":\"kk\"}', '2023-07-29 09:26:22', '2023-07-29 09:26:22');
INSERT INTO `admin_operation_log` VALUES (21, 1, 'admin/paw.cache', 'POST', '172.27.0.1', '{\"_method\":\"post\",\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\",\"source\":\"kk\"}', '2023-07-29 09:26:52', '2023-07-29 09:26:52');
INSERT INTO `admin_operation_log` VALUES (22, 1, 'admin/paw.cache', 'POST', '172.27.0.1', '{\"_method\":\"post\",\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\",\"source\":\"kk\"}', '2023-07-29 09:27:22', '2023-07-29 09:27:22');
INSERT INTO `admin_operation_log` VALUES (23, 1, 'admin/paw.cache', 'POST', '172.27.0.1', '{\"_method\":\"post\",\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\",\"source\":\"kk\"}', '2023-07-29 09:27:52', '2023-07-29 09:27:52');
INSERT INTO `admin_operation_log` VALUES (24, 1, 'admin/paw.cache', 'POST', '172.27.0.1', '{\"_method\":\"post\",\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\",\"source\":\"kk\"}', '2023-07-29 09:28:22', '2023-07-29 09:28:22');
INSERT INTO `admin_operation_log` VALUES (25, 1, 'admin/paw.cache', 'POST', '172.27.0.1', '{\"_method\":\"post\",\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\",\"source\":\"kk\"}', '2023-07-29 09:28:52', '2023-07-29 09:28:52');
INSERT INTO `admin_operation_log` VALUES (26, 1, 'admin/paw.cache', 'POST', '172.27.0.1', '{\"_method\":\"post\",\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\",\"source\":\"kk\"}', '2023-07-29 09:29:22', '2023-07-29 09:29:22');
INSERT INTO `admin_operation_log` VALUES (27, 1, 'admin/paw.cache', 'POST', '172.27.0.1', '{\"_method\":\"post\",\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\",\"source\":\"kk\"}', '2023-07-29 09:29:52', '2023-07-29 09:29:52');
INSERT INTO `admin_operation_log` VALUES (28, 1, 'admin/paw', 'GET', '172.27.0.1', '[]', '2023-07-29 09:30:22', '2023-07-29 09:30:22');
INSERT INTO `admin_operation_log` VALUES (29, 1, 'admin/paw.cache', 'POST', '172.27.0.1', '{\"_method\":\"post\",\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\",\"source\":\"kk\"}', '2023-07-29 09:30:22', '2023-07-29 09:30:22');
INSERT INTO `admin_operation_log` VALUES (30, 1, 'admin/paw.cache', 'POST', '172.27.0.1', '{\"_method\":\"post\",\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\",\"source\":\"kk\"}', '2023-07-29 09:30:30', '2023-07-29 09:30:30');
INSERT INTO `admin_operation_log` VALUES (31, 1, 'admin/paw', 'GET', '172.27.0.1', '[]', '2023-07-29 09:30:46', '2023-07-29 09:30:46');
INSERT INTO `admin_operation_log` VALUES (32, 1, 'admin/paw.cache', 'POST', '172.27.0.1', '{\"_method\":\"post\",\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\",\"source\":\"kk\"}', '2023-07-29 09:30:47', '2023-07-29 09:30:47');
INSERT INTO `admin_operation_log` VALUES (33, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-07-29 09:30:56', '2023-07-29 09:30:56');
INSERT INTO `admin_operation_log` VALUES (34, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"_pjax\":\"#pjax-container\",\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 09:31:04', '2023-07-29 09:31:04');
INSERT INTO `admin_operation_log` VALUES (35, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 09:31:10', '2023-07-29 09:31:10');
INSERT INTO `admin_operation_log` VALUES (36, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 09:31:15', '2023-07-29 09:31:15');
INSERT INTO `admin_operation_log` VALUES (37, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 09:31:17', '2023-07-29 09:31:17');
INSERT INTO `admin_operation_log` VALUES (38, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 09:31:25', '2023-07-29 09:31:25');
INSERT INTO `admin_operation_log` VALUES (39, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 09:31:32', '2023-07-29 09:31:32');
INSERT INTO `admin_operation_log` VALUES (40, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 09:31:43', '2023-07-29 09:31:43');
INSERT INTO `admin_operation_log` VALUES (41, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 09:31:46', '2023-07-29 09:31:46');
INSERT INTO `admin_operation_log` VALUES (42, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 09:32:27', '2023-07-29 09:32:27');
INSERT INTO `admin_operation_log` VALUES (43, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 09:32:54', '2023-07-29 09:32:54');
INSERT INTO `admin_operation_log` VALUES (44, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 09:34:59', '2023-07-29 09:34:59');
INSERT INTO `admin_operation_log` VALUES (45, 1, 'admin/source_chapter/274464/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 09:35:04', '2023-07-29 09:35:04');
INSERT INTO `admin_operation_log` VALUES (46, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 09:37:24', '2023-07-29 09:37:24');
INSERT INTO `admin_operation_log` VALUES (47, 1, 'admin/source_chapter/274466/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 09:37:27', '2023-07-29 09:37:27');
INSERT INTO `admin_operation_log` VALUES (48, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 09:37:30', '2023-07-29 09:37:30');
INSERT INTO `admin_operation_log` VALUES (49, 1, 'admin/source_chapter/274466/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 09:37:32', '2023-07-29 09:37:32');
INSERT INTO `admin_operation_log` VALUES (50, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 09:50:18', '2023-07-29 09:50:18');
INSERT INTO `admin_operation_log` VALUES (51, 1, 'admin/source_chapter/274516/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 09:50:20', '2023-07-29 09:50:20');
INSERT INTO `admin_operation_log` VALUES (52, 1, 'admin/source_chapter/274501/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 09:50:28', '2023-07-29 09:50:28');
INSERT INTO `admin_operation_log` VALUES (53, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"4\",\"_pjax\":\"#pjax-container\"}', '2023-07-29 09:50:35', '2023-07-29 09:50:35');
INSERT INTO `admin_operation_log` VALUES (54, 1, 'admin/source_chapter/274450/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 09:50:38', '2023-07-29 09:50:38');
INSERT INTO `admin_operation_log` VALUES (55, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"4\"}', '2023-07-29 09:56:05', '2023-07-29 09:56:05');
INSERT INTO `admin_operation_log` VALUES (56, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"4\"}', '2023-07-29 09:56:09', '2023-07-29 09:56:09');
INSERT INTO `admin_operation_log` VALUES (57, 1, 'admin/source_chapter/274452/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 09:56:12', '2023-07-29 09:56:12');
INSERT INTO `admin_operation_log` VALUES (58, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"1\",\"_pjax\":\"#pjax-container\"}', '2023-07-29 09:56:18', '2023-07-29 09:56:18');
INSERT INTO `admin_operation_log` VALUES (59, 1, 'admin/source_chapter/274521/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 09:56:23', '2023-07-29 09:56:23');
INSERT INTO `admin_operation_log` VALUES (60, 1, 'admin/source_chapter/274520/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 09:56:29', '2023-07-29 09:56:29');
INSERT INTO `admin_operation_log` VALUES (61, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"_pjax\":\"#pjax-container\",\"page\":\"5\"}', '2023-07-29 09:56:35', '2023-07-29 09:56:35');
INSERT INTO `admin_operation_log` VALUES (62, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"5\"}', '2023-07-29 10:00:33', '2023-07-29 10:00:33');
INSERT INTO `admin_operation_log` VALUES (63, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"1\",\"_pjax\":\"#pjax-container\"}', '2023-07-29 10:00:39', '2023-07-29 10:00:39');
INSERT INTO `admin_operation_log` VALUES (64, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"1\"}', '2023-07-29 10:04:15', '2023-07-29 10:04:15');
INSERT INTO `admin_operation_log` VALUES (65, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"7\",\"_pjax\":\"#pjax-container\"}', '2023-07-29 10:04:21', '2023-07-29 10:04:21');
INSERT INTO `admin_operation_log` VALUES (66, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"_pjax\":\"#pjax-container\",\"page\":\"5\"}', '2023-07-29 10:04:25', '2023-07-29 10:04:25');
INSERT INTO `admin_operation_log` VALUES (67, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"_pjax\":\"#pjax-container\",\"page\":\"2\"}', '2023-07-29 10:04:29', '2023-07-29 10:04:29');
INSERT INTO `admin_operation_log` VALUES (68, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"_pjax\":\"#pjax-container\",\"page\":\"1\"}', '2023-07-29 10:04:32', '2023-07-29 10:04:32');
INSERT INTO `admin_operation_log` VALUES (69, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"_pjax\":\"#pjax-container\",\"page\":\"5\"}', '2023-07-29 10:04:35', '2023-07-29 10:04:35');
INSERT INTO `admin_operation_log` VALUES (70, 1, 'admin/auth/menu', 'GET', '172.27.0.1', '[]', '2023-07-29 10:08:05', '2023-07-29 10:08:05');
INSERT INTO `admin_operation_log` VALUES (71, 1, 'admin/source_comic', 'GET', '172.27.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-07-29 10:08:08', '2023-07-29 10:08:08');
INSERT INTO `admin_operation_log` VALUES (72, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-07-29 10:08:10', '2023-07-29 10:08:10');
INSERT INTO `admin_operation_log` VALUES (73, 1, 'admin/fail_info', 'GET', '172.27.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-07-29 10:08:12', '2023-07-29 10:08:12');
INSERT INTO `admin_operation_log` VALUES (74, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-07-29 10:08:17', '2023-07-29 10:08:17');
INSERT INTO `admin_operation_log` VALUES (75, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"_pjax\":\"#pjax-container\",\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 10:08:22', '2023-07-29 10:08:22');
INSERT INTO `admin_operation_log` VALUES (76, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 10:08:28', '2023-07-29 10:08:28');
INSERT INTO `admin_operation_log` VALUES (77, 1, 'admin/source_chapter/274521/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 10:08:32', '2023-07-29 10:08:32');
INSERT INTO `admin_operation_log` VALUES (78, 1, 'admin/source_chapter/274511/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 10:08:54', '2023-07-29 10:08:54');
INSERT INTO `admin_operation_log` VALUES (79, 1, 'admin/source_chapter/274521/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 10:09:02', '2023-07-29 10:09:02');
INSERT INTO `admin_operation_log` VALUES (80, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 10:09:09', '2023-07-29 10:09:09');
INSERT INTO `admin_operation_log` VALUES (81, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 10:09:14', '2023-07-29 10:09:14');
INSERT INTO `admin_operation_log` VALUES (82, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 10:09:18', '2023-07-29 10:09:18');
INSERT INTO `admin_operation_log` VALUES (83, 1, 'admin/paw', 'GET', '172.27.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-07-29 10:09:23', '2023-07-29 10:09:23');
INSERT INTO `admin_operation_log` VALUES (84, 1, 'admin/paw.cache', 'POST', '172.27.0.1', '{\"_method\":\"post\",\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\",\"source\":\"kk\"}', '2023-07-29 10:09:24', '2023-07-29 10:09:24');
INSERT INTO `admin_operation_log` VALUES (85, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"_pjax\":\"#pjax-container\"}', '2023-07-29 10:09:27', '2023-07-29 10:09:27');
INSERT INTO `admin_operation_log` VALUES (86, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"_pjax\":\"#pjax-container\",\"id\":null,\"comic_id\":null,\"title\":null,\"source\":null,\"is_free\":null,\"status\":null,\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"create_at\":{\"start\":null,\"end\":null}}', '2023-07-29 10:09:33', '2023-07-29 10:09:33');
INSERT INTO `admin_operation_log` VALUES (87, 1, 'admin/paw.cache', 'POST', '172.27.0.1', '{\"_method\":\"post\",\"_token\":\"AlPRdpShT81aozkX3umjrOClDU9U9DI8opbzvNQX\",\"source\":\"kk\"}', '2023-07-29 10:09:54', '2023-07-29 10:09:54');
INSERT INTO `admin_operation_log` VALUES (88, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"_pjax\":\"#pjax-container\",\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"2\"}', '2023-07-29 10:09:54', '2023-07-29 10:09:54');
INSERT INTO `admin_operation_log` VALUES (89, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"_pjax\":\"#pjax-container\",\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"1\"}', '2023-07-29 10:09:57', '2023-07-29 10:09:57');
INSERT INTO `admin_operation_log` VALUES (90, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"1\"}', '2023-07-29 10:10:03', '2023-07-29 10:10:03');
INSERT INTO `admin_operation_log` VALUES (91, 1, 'admin/source_chapter/274508/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 10:10:10', '2023-07-29 10:10:10');
INSERT INTO `admin_operation_log` VALUES (92, 1, 'admin/source_comic', 'GET', '172.27.0.1', '{\"id\":\"13968\"}', '2023-07-29 10:10:12', '2023-07-29 10:10:12');
INSERT INTO `admin_operation_log` VALUES (93, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"comic_id\":\"13968\"}', '2023-07-29 10:10:15', '2023-07-29 10:10:15');
INSERT INTO `admin_operation_log` VALUES (94, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"comic_id\":\"13968\",\"page\":\"2\",\"_pjax\":\"#pjax-container\"}', '2023-07-29 10:10:18', '2023-07-29 10:10:18');
INSERT INTO `admin_operation_log` VALUES (95, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"comic_id\":\"13968\",\"_pjax\":\"#pjax-container\",\"page\":\"3\"}', '2023-07-29 10:10:20', '2023-07-29 10:10:20');
INSERT INTO `admin_operation_log` VALUES (96, 1, 'admin/source_chapter/274473/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 10:10:24', '2023-07-29 10:10:24');
INSERT INTO `admin_operation_log` VALUES (97, 1, 'admin/source_comic', 'GET', '172.27.0.1', '{\"id\":\"13970\"}', '2023-07-29 10:10:36', '2023-07-29 10:10:36');
INSERT INTO `admin_operation_log` VALUES (98, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"comic_id\":\"13970\"}', '2023-07-29 10:10:39', '2023-07-29 10:10:39');
INSERT INTO `admin_operation_log` VALUES (99, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"comic_id\":\"13970\",\"page\":\"2\",\"_pjax\":\"#pjax-container\"}', '2023-07-29 10:10:42', '2023-07-29 10:10:42');
INSERT INTO `admin_operation_log` VALUES (100, 1, 'admin/source_chapter/274561/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 10:10:45', '2023-07-29 10:10:45');
INSERT INTO `admin_operation_log` VALUES (101, 1, 'admin/source_chapter/274569/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 10:10:52', '2023-07-29 10:10:52');
INSERT INTO `admin_operation_log` VALUES (102, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"1\"}', '2023-07-29 10:11:00', '2023-07-29 10:11:00');
INSERT INTO `admin_operation_log` VALUES (103, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"1\"}', '2023-07-29 10:11:05', '2023-07-29 10:11:05');
INSERT INTO `admin_operation_log` VALUES (104, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"1\"}', '2023-07-29 10:11:10', '2023-07-29 10:11:10');
INSERT INTO `admin_operation_log` VALUES (105, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"5\",\"_pjax\":\"#pjax-container\"}', '2023-07-29 10:11:19', '2023-07-29 10:11:19');
INSERT INTO `admin_operation_log` VALUES (106, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"_pjax\":\"#pjax-container\",\"page\":\"8\"}', '2023-07-29 10:11:24', '2023-07-29 10:11:24');
INSERT INTO `admin_operation_log` VALUES (107, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"_pjax\":\"#pjax-container\",\"page\":\"1\"}', '2023-07-29 10:11:27', '2023-07-29 10:11:27');
INSERT INTO `admin_operation_log` VALUES (108, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"1\"}', '2023-07-29 10:14:34', '2023-07-29 10:14:34');
INSERT INTO `admin_operation_log` VALUES (109, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"1\"}', '2023-07-29 10:19:20', '2023-07-29 10:19:20');
INSERT INTO `admin_operation_log` VALUES (110, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"1\",\"_sort\":{\"column\":\"updated_at\",\"type\":\"desc\"},\"_pjax\":\"#pjax-container\"}', '2023-07-29 10:19:28', '2023-07-29 10:19:28');
INSERT INTO `admin_operation_log` VALUES (111, 1, 'admin/source_chapter/274615/edit', 'GET', '172.27.0.1', '[]', '2023-07-29 10:19:34', '2023-07-29 10:19:34');
INSERT INTO `admin_operation_log` VALUES (112, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"1\",\"_sort\":{\"column\":\"updated_at\",\"type\":\"desc\"}}', '2023-07-29 10:21:27', '2023-07-29 10:21:27');
INSERT INTO `admin_operation_log` VALUES (113, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"1\",\"_sort\":{\"column\":\"updated_at\",\"type\":\"desc\"}}', '2023-07-29 10:24:34', '2023-07-29 10:24:34');
INSERT INTO `admin_operation_log` VALUES (114, 1, 'admin/source_chapter', 'GET', '172.27.0.1', '{\"73f68bba75445ada15c3cab5bc5adc19\":\"1\",\"page\":\"1\",\"_sort\":{\"column\":\"updated_at\",\"type\":\"desc\"}}', '2023-07-29 10:24:38', '2023-07-29 10:24:38');

-- ----------------------------
-- Table structure for admin_permissions
-- ----------------------------
DROP TABLE IF EXISTS `admin_permissions`;
CREATE TABLE `admin_permissions`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `slug` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `http_method` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `http_path` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `admin_permissions_name_unique`(`name`) USING BTREE,
  UNIQUE INDEX `admin_permissions_slug_unique`(`slug`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_permissions
-- ----------------------------
INSERT INTO `admin_permissions` VALUES (1, 'All permission', '*', '', '*', NULL, NULL);
INSERT INTO `admin_permissions` VALUES (2, 'Dashboard', 'dashboard', 'GET', '/', NULL, NULL);
INSERT INTO `admin_permissions` VALUES (3, 'Login', 'auth.login', '', '/auth/login\r\n/auth/logout', NULL, NULL);
INSERT INTO `admin_permissions` VALUES (4, 'User setting', 'auth.setting', 'GET,PUT', '/auth/setting', NULL, NULL);
INSERT INTO `admin_permissions` VALUES (5, 'Auth management', 'auth.management', '', '/auth/roles\r\n/auth/permissions\r\n/auth/menu\r\n/auth/logs', NULL, NULL);

-- ----------------------------
-- Table structure for admin_role_menu
-- ----------------------------
DROP TABLE IF EXISTS `admin_role_menu`;
CREATE TABLE `admin_role_menu`  (
  `role_id` int(11) NOT NULL,
  `menu_id` int(11) NOT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  INDEX `admin_role_menu_role_id_menu_id_index`(`role_id`, `menu_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_role_menu
-- ----------------------------
INSERT INTO `admin_role_menu` VALUES (1, 2, NULL, NULL);

-- ----------------------------
-- Table structure for admin_role_permissions
-- ----------------------------
DROP TABLE IF EXISTS `admin_role_permissions`;
CREATE TABLE `admin_role_permissions`  (
  `role_id` int(11) NOT NULL,
  `permission_id` int(11) NOT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  INDEX `admin_role_permissions_role_id_permission_id_index`(`role_id`, `permission_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_role_permissions
-- ----------------------------
INSERT INTO `admin_role_permissions` VALUES (1, 1, NULL, NULL);

-- ----------------------------
-- Table structure for admin_role_users
-- ----------------------------
DROP TABLE IF EXISTS `admin_role_users`;
CREATE TABLE `admin_role_users`  (
  `role_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  INDEX `admin_role_users_role_id_user_id_index`(`role_id`, `user_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_role_users
-- ----------------------------
INSERT INTO `admin_role_users` VALUES (1, 1, NULL, NULL);

-- ----------------------------
-- Table structure for admin_roles
-- ----------------------------
DROP TABLE IF EXISTS `admin_roles`;
CREATE TABLE `admin_roles`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `slug` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `admin_roles_name_unique`(`name`) USING BTREE,
  UNIQUE INDEX `admin_roles_slug_unique`(`slug`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_roles
-- ----------------------------
INSERT INTO `admin_roles` VALUES (1, 'Administrator', 'administrator', '2023-07-29 09:20:40', '2023-07-29 09:20:40');

-- ----------------------------
-- Table structure for admin_user_permissions
-- ----------------------------
DROP TABLE IF EXISTS `admin_user_permissions`;
CREATE TABLE `admin_user_permissions`  (
  `user_id` int(11) NOT NULL,
  `permission_id` int(11) NOT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  INDEX `admin_user_permissions_user_id_permission_id_index`(`user_id`, `permission_id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for admin_users
-- ----------------------------
DROP TABLE IF EXISTS `admin_users`;
CREATE TABLE `admin_users`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `username` varchar(190) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `remember_token` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `admin_users_username_unique`(`username`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 2 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of admin_users
-- ----------------------------
INSERT INTO `admin_users` VALUES (1, 'admin', '$2y$10$KmDos9C4ebRrN80da2BGv.Pejp.IDQxJD7CyYZsOcK/Jdla5HoSlG', 'Administrator', NULL, 'ocFHS3VV6tF8nqHTkMDnVHjc7nE1gEKKEX35zr8uOWSxfNcefugu5O4hVcn8', '2023-07-29 09:20:40', '2023-07-29 09:20:40');

-- ----------------------------
-- Table structure for fail_info
-- ----------------------------
DROP TABLE IF EXISTS `fail_info`;
CREATE TABLE `fail_info`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `source` tinyint(1) NOT NULL DEFAULT 1 COMMENT '采集源 1:快看 2:腾讯',
  `type` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0漫画列表 1漫画 2章节 3图片',
  `err` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '错误关键词',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '地址',
  `info` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '失败信息记录',
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `type`(`type`) USING BTREE,
  INDEX `err`(`err`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for failed_jobs
-- ----------------------------
DROP TABLE IF EXISTS `failed_jobs`;
CREATE TABLE `failed_jobs`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `connection` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `queue` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `payload` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `exception` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `failed_at` timestamp(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for migrations
-- ----------------------------
DROP TABLE IF EXISTS `migrations`;
CREATE TABLE `migrations`  (
  `id` int(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `migration` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `batch` int(11) NOT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of migrations
-- ----------------------------
INSERT INTO `migrations` VALUES (1, '2014_10_12_000000_create_users_table', 1);
INSERT INTO `migrations` VALUES (2, '2014_10_12_100000_create_password_resets_table', 1);
INSERT INTO `migrations` VALUES (3, '2016_01_04_173148_create_admin_tables', 1);
INSERT INTO `migrations` VALUES (4, '2019_08_19_000000_create_failed_jobs_table', 1);

-- ----------------------------
-- Table structure for password_resets
-- ----------------------------
DROP TABLE IF EXISTS `password_resets`;
CREATE TABLE `password_resets`  (
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `token` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  INDEX `password_resets_email_index`(`email`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for source_chapter
-- ----------------------------
DROP TABLE IF EXISTS `source_chapter`;
CREATE TABLE `source_chapter`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sid` int(11) NOT NULL DEFAULT 0,
  `scid` int(11) NOT NULL DEFAULT 0,
  `comic_id` int(11) NOT NULL,
  `source` tinyint(1) NOT NULL DEFAULT 1 COMMENT '采集源 1:快看 2:腾讯',
  `source_chapter_id` int(11) NOT NULL COMMENT '源章节id',
  `source_url` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '源url',
  `cover` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '封面',
  `title` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标题',
  `sort` int(11) NOT NULL DEFAULT 0,
  `is_free` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0免费 1收费',
  `source_data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0未审核 1通过',
  `retry` tinyint(1) NOT NULL DEFAULT 0,
  `view_type` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0条漫 1页漫',
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `source`(`source`, `comic_id`, `source_url`) USING BTREE,
  INDEX `comic_id`(`comic_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 323197 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '采集-漫画章节' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for source_comic
-- ----------------------------
DROP TABLE IF EXISTS `source_comic`;
CREATE TABLE `source_comic`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `sid` int(11) NOT NULL DEFAULT 0,
  `source` tinyint(1) NOT NULL DEFAULT 1 COMMENT '采集源 1:快看 2:腾讯',
  `source_id` int(11) NOT NULL COMMENT '源漫画id',
  `source_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '源url',
  `cover` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '封面',
  `title` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '标题',
  `author` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '作者',
  `label` json NOT NULL COMMENT '标签',
  `category` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '分类',
  `region` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '地区',
  `chapter_count` int(11) NOT NULL DEFAULT 0 COMMENT '章节数量',
  `chapter_count_download` int(11) NOT NULL DEFAULT 0 COMMENT '章节数量(已下载)',
  `like` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '喜欢',
  `popularity` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '0' COMMENT '人气热度',
  `is_free` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0免费 1收费',
  `is_finish` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0连载 1完结',
  `description` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '描述',
  `source_data` text CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '源数据',
  `chapter_pick` int(11) NOT NULL DEFAULT 0 COMMENT '章节拨片',
  `retry` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0默认 1重抓',
  `status` tinyint(1) NOT NULL DEFAULT 0 COMMENT '0未审核 1通过',
  `last_chapter_update_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最新章节更新时间',
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP(0),
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `cover_h` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '横板封面',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `source_id`(`source`, `source_id`) USING BTREE,
  UNIQUE INDEX `source_uri`(`source_url`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 35795 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '采集-漫画' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for source_image
-- ----------------------------
DROP TABLE IF EXISTS `source_image`;
CREATE TABLE `source_image`  (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `source` tinyint(1) NULL DEFAULT 0 COMMENT '采集源 1:快看 2:腾讯',
  `comic_id` int(11) NOT NULL DEFAULT 0,
  `chapter_id` int(11) NOT NULL,
  `images` json NOT NULL,
  `source_data` json NOT NULL,
  `state` tinyint(1) NOT NULL DEFAULT 0 COMMENT '资源获取:0未开始 1已完成',
  `updated_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_at` datetime(0) NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `chapter_id`(`chapter_id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 361 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `email_verified_at` timestamp(0) NULL DEFAULT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `remember_token` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NULL DEFAULT NULL,
  `created_at` timestamp(0) NULL DEFAULT NULL,
  `updated_at` timestamp(0) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `users_email_unique`(`email`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_unicode_ci ROW_FORMAT = Dynamic;

SET FOREIGN_KEY_CHECKS = 1;
