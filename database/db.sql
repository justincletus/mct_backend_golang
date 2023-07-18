CREATE DATABASE IF NOT EXISTS cms;
USE cms;

create table if not exists `users`(
  `id` int(10) not null auto_increment,
  `email` varchar(50) not null unique,
  `username` varchar(50) not null unique,
  `fullname` varchar(50) not null,
  `password` varchar(250) not null,
  `mobile` varchar(50),
  `code` varchar(50),
  `email_verified` tinyint(1) default '0',
  `role` varchar(30),
  `created_at` datetime not null default current_timestamp,
  `updated_at` datetime not null default current_timestamp on update current_timestamp,
  `deleted_at` datetime,
  primary key(`id`)
  ) engine=innodb auto_increment=1000 default charset=latin1;

CREATE TABLE IF NOT EXISTS `projects` (
    `id` INT(10) NOT NULL AUTO_INCREMENT,
    `project_id` VARCHAR(10) UNIQUE,
    `p_name` VARCHAR(10),
    `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1000 DEFAULT CHARSET=latin1;

CREATE TABLE IF NOT EXISTS `orders`(
  `id` int not null auto_increment,
  `project` VARCHAR(30) not null,
  `requisition_no` int not null,
  `purchase_order_no` int not null,
  `delivery_note_no` int,
  `date_of_delivery` datetime not null default current_timestamp,
  `description` longtext,
  `created_at` datetime not null default current_timestamp,
  `updated_at` datetime not null default current_timestamp on update current_timestamp,
  `user_id` int not null,
  primary key(`id`),
  key `ord_key` (`user_id`),
  constraint `ord_key` foreign key(`user_id`) references `users`(`id`) on delete cascade
) engine=innodb auto_increment=1000 default charset=latin1;

create table if not exists jobs(
  `id` int not null AUTO_INCREMENT primary KEY,
  `name` VARCHAR(50) null null,
  `created_at` datetime not null default current_timestamp,
  `updated_at` datetime not null default current_timestamp on update current_timestamp,
  `user_id` int not null,
  key `job_fk`(`user_id`),
  constraint `job_fk` foreign key(`user_id`) references `users`(`id`) on delete cascade 
) engine=innodb auto_increment=1000 default charset=latin1;

create table if not exists team_mems(
  `id` int not null auto_increment primary key,
  `title` varchar(50) not null,  
  `member1` varchar(30),
  `member2` varchar(30),
  `member3` varchar(30),
  `member4` varchar(30),
  `member5` varchar(30),
  `members` longtext,
  `created_at` datetime not null default current_timestamp,
  `updated_at` datetime not null default current_timestamp on update current_timestamp,
  `user_id` int not null,
  key `team_mems_fk`(`user_id`),
  constraint `team_mems_fk` foreign key(`user_id`) references `users`(`id`) on delete cascade
) engine=innodb auto_increment=1000 default charset=latin1;
  

CREATE TABLE IF NOT EXISTS 
  `managers` (
    `id` bigint unsigned NOT NULL AUTO_INCREMENT,
    `created_at` datetime(3) DEFAULT NULL,
    `updated_at` datetime(3) DEFAULT NULL,
    `deleted_at` datetime(3) DEFAULT NULL,
    `uid` bigint unsigned DEFAULT NULL,
    `user_id` bigint unsigned DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_managers_deleted_at` (`deleted_at`),
    KEY `idx_managers_uid` (`uid`),
    KEY `fk_users_manager` (`user_id`),
    CONSTRAINT `fk_managers_user` FOREIGN KEY (`uid`) REFERENCES `users` (`id`),
    CONSTRAINT `fk_users_manager` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
  ) ENGINE = InnoDB AUTO_INCREMENT = 3 DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci



use cms;
  CREATE TABLE IF NOT EXISTS `mri_reports` (
    `id` INT not null auto_increment,    
    `purchase_requisition` tinyint(1) not null default '0',
    `is_quality` tinyint(1) not null default '0',
    `is_quantity` tinyint(1) not null default '0',
    `is_damaged` tinyint(1) not null default '0',
    `is_sample_same` tinyint(1) not null default '0',
    `is_any_certification` tinyint(1) not null default '0',
    `is_document` tinyint(1) not null default '0',
    `is_material_certification` tinyint(1) not null default '0',
    `is_mill_certification` tinyint(1) not null default '0',
    `is_applied_finish` tinyint(1) not null default '0',
    `is_test_result` tinyint(1) not null default '0',
    `is_data_sheet` tinyint(1) not null default '0',
    `is_other` tinyint(1) not null default '0',
    `is_spare_delivery` tinyint(1) not null default '0',
    `is_material_comply` tinyint(1) not null default '0',
    `comment` varchar(255),
    `name` varchar(50) not null,
    `signature` varchar(50),
    `created_at` datetime not null default current_timestamp,
    `updated_at` datetime not null default current_timestamp on update current_timestamp,
    `project_id` INT not null
    primary key(`id`),
    key `mri_fkey` (`project_id`),
    constraint `mri_fk` foreign key(`project_id`) references `projects`(`id`) on delete cascade
  ) engine=innodb auto_increment=1000 default charset=latin1;


  

