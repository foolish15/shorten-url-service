CREATE TABLE IF NOT EXISTS `blocks` (
  `type` enum('scheme','domain','regex') NOT NULL,
  `value` varchar(255) NOT NULL,
  PRIMARY KEY (`type`,`value`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


REPLACE INTO `blocks`
SET `type` = 'scheme',
    `value` = 'ftp';
    
REPLACE INTO `blocks`
SET `type` = 'scheme',
    `value` = 'line';

REPLACE INTO `blocks`
SET `type` = 'domain',
    `value` = 'blockdomain';

REPLACE INTO `blocks`
SET `type` = 'regex',
    `value` = '(?m)block-regex';