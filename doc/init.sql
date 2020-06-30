CREATE DATABASE IF NOT EXISTS capacity;

DROP TABLE IF EXISTS `r_work`;
CREATE TABLE `r_work` (
  `guid` VARCHAR(255) NOT NULL,
  `name` VARCHAR(255) NOT NULL DEFAULT '',
  `workspace` VARCHAR(50) NOT NULL DEFAULT '',
  `endpoint_a` VARCHAR(255) NOT NULL DEFAULT '',
  `endpoint_b` VARCHAR(255) NOT NULL DEFAULT '',
  `metric_a` VARCHAR(255) NOT NULL DEFAULT '',
  `metric_b` VARCHAR(255) NOT NULL DEFAULT '',
  `time_select` VARCHAR(50) NOT NULL DEFAULT '',
  `legend_x` VARCHAR(255) NOT NULL DEFAULT '',
  `legend_y` VARCHAR(255) NOT NULL DEFAULT '',
  `output` VARCHAR(2000) NOT NULL DEFAULT '',
  `expr` VARCHAR(50) NOT NULL DEFAULT '',
  `func_a` VARCHAR(50) NOT NULL DEFAULT '',
  `func_b` VARCHAR(50) NOT NULL DEFAULT '',
  `level` INT NOT NULL DEFAULT 0,
  `update_at` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`guid`)
) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;