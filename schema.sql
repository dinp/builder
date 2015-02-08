USE builder;

SET NAMES utf8;

DROP TABLE IF EXISTS build;
CREATE TABLE build (
  id        INT UNSIGNED  NOT NULL AUTO_INCREMENT PRIMARY KEY,
  app       VARCHAR(64)   NOT NULL DEFAULT '',
  version   VARCHAR(64)   NOT NULL DEFAULT '',
  resume    VARCHAR(255)  NOT NULL DEFAULT '',
  base      VARCHAR(1024) NOT NULL DEFAULT '',
  image     VARCHAR(1024) NOT NULL DEFAULT '',
  tarball   VARCHAR(1024) NOT NULL DEFAULT '',
  repo      VARCHAR(255)  NOT NULL DEFAULT '',
  branch    VARCHAR(64)   NOT NULL DEFAULT '' COMMENT 'branch or tag',
  status    VARCHAR(255)  NOT NULL DEFAULT 'saved in db',
  user_id   INT UNSIGNED  NOT NULL,
  user_name VARCHAR(64)   NOT NULL,
  create_at DATETIME      NOT NULL,
  KEY idx_build_app(app),
  KEY idx_build_user_id(user_id)
)
  ENGINE =innodb
  DEFAULT CHARSET =utf8
  COLLATE =utf8_general_ci;

  
