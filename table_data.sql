/*
 Navicat Premium Data Transfer

 Source Server         : 127.0.0.1_5432
 Source Server Type    : PostgreSQL
 Source Server Version : 110005
 Source Host           : 127.0.0.1:5432
 Source Catalog        : excel
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 110005
 File Encoding         : 65001

 Date: 19/11/2019 14:42:17
*/


-- ----------------------------
-- Table structure for table_data
-- ----------------------------
DROP TABLE IF EXISTS "public"."table_data";
CREATE TABLE "public"."table_data" (
  "table_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "describe" varchar(255) COLLATE "pg_catalog"."default",
  "version" int8 NOT NULL,
  "date" varchar(64) COLLATE "pg_catalog"."default",
  "data" varchar(65535) COLLATE "pg_catalog"."default" NOT NULL
)
;
ALTER TABLE "public"."table_data" OWNER TO "dbuser";

-- ----------------------------
-- Primary Key structure for table table_data
-- ----------------------------
ALTER TABLE "public"."table_data" ADD CONSTRAINT "table_data_pkey" PRIMARY KEY ("table_name");
