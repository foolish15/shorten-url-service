CREATE TABLE IF NOT EXISTS `access_transactions` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `link_id` bigint(20) unsigned DEFAULT NULL,
  `link_url` varchar(255) DEFAULT NULL,
  `browser` varchar(255) DEFAULT NULL,
  `browser_version` longtext,
  `os` varchar(255) DEFAULT NULL,
  `os_version` longtext,
  `device_type` varchar(255) DEFAULT NULL,
  `user_agent` longtext,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_access_transactions_device_type` (`device_type`),
  KEY `idx_access_transactions_link_id` (`link_id`),
  KEY `idx_access_transactions_browser` (`browser`),
  KEY `idx_access_transactions_os` (`os`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `blocks` (
  `type` enum('scheme','domain','regex') NOT NULL,
  `value` varchar(255) NOT NULL,
  PRIMARY KEY (`type`,`value`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `links` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `link` text,
  `expire` bigint(20) DEFAULT NULL,
  `code` varchar(255) DEFAULT NULL,
  `hit` bigint(20) unsigned DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_links_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `tokens` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `sub` varchar(255) DEFAULT NULL,
  `exp` bigint(20) DEFAULT NULL,
  `u` json DEFAULT NULL,
  `code` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_tokens_code` (`code`),
  KEY `idx_tokens_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `users` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `email` varchar(255) DEFAULT NULL,
  `phone` varchar(255) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_users_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `user_auths` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint(20) unsigned DEFAULT NULL,
  `channel` enum('password','line','facebook','google') DEFAULT NULL,
  `channel_id` varchar(255) DEFAULT NULL,
  `channel_secret` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uidx_user_id_channel` (`user_id`,`channel`),
  CONSTRAINT `fk_users_user_auths` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;




REPLACE INTO `blocks`
SET `type` = 'scheme',
    `value` = 'ftp';
    
REPLACE INTO `blocks`
SET `type` = 'scheme',
    `value` = 'rtmp';

REPLACE INTO `blocks`
SET `type` = 'domain',
    `value` = 'blockdomain';

REPLACE INTO `blocks`
SET `type` = 'regex',
    `value` = '(?m)block-regex';