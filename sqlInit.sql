-- ----------------------------
-- Table structure for tag_data
-- ----------------------------
DROP TABLE IF EXISTS "public"."tag_data";
CREATE TABLE "public"."tag_data" (
  "tag_name" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "date" varchar(64) COLLATE "pg_catalog"."default",
  "describe" varchar(128) COLLATE "pg_catalog"."default"
)
;
ALTER TABLE "public"."tag_data" OWNER TO "dbuser";

-- ----------------------------
-- Primary Key structure for table tag_data
-- ----------------------------
ALTER TABLE "public"."tag_data" ADD CONSTRAINT "tag_data_pkey" PRIMARY KEY ("tag_name");

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


-- ----------------------------
-- Table structure for user
-- ----------------------------
DROP TABLE IF EXISTS "public"."user";
CREATE TABLE "public"."user" (
  "user_name" varchar(64) COLLATE "pg_catalog"."default" NOT NULL,
  "password" varchar(255) COLLATE "pg_catalog"."default" NOT NULL,
  "permission" varchar(255) COLLATE "pg_catalog"."default" NOT NULL
)
;
ALTER TABLE "public"."user" OWNER TO "dbuser";

-- ----------------------------
-- Primary Key structure for table user
-- ----------------------------
ALTER TABLE "public"."user" ADD CONSTRAINT "user_pkey" PRIMARY KEY ("user_name");