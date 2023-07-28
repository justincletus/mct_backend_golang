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

  alter table `users` add column `status` enum('active', 'inactive') default 'inactive';

alter table `users` MODIFY column role VARCHAR(20) default "user";

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

alter table `orders` add column `report_id` int;
alter table `orders` add constraint `order_rept_fk` foreign key (`report_id`) references `reports` (`id`) on update cascade on delete cascade;

alter table `orders` MODIFY COLUMN `requisition_no` varchar(30);
alter table `orders` MODIFY column `purchase_order_no` VARCHAR(30);
alter table `orders` MODIFY COLUMN `delivery_note_no` VARCHAR(30);
alter table `orders` drop foreign key `ord_key`;

alter table `orders` drop column `user_id`;

alter table `orders` add column `job_id` int not null after `updated_at`;

alter table `orders` add constraint `ord_fk` foreign key (`job_id`) references `jobs`(`id`) on update CASCADE on delete CASCADE;

create table if not exists jobs(
  `id` int not null AUTO_INCREMENT primary KEY,
  `name` VARCHAR(50) null,
  `created_at` datetime not null default current_timestamp,
  `updated_at` datetime not null default current_timestamp on update current_timestamp,
  `user_id` int not null,
  key `job_fk`(`user_id`),
  constraint `job_fk` foreign key(`user_id`) references `users`(`id`) on delete cascade 
) engine=innodb auto_increment=1000 default charset=latin1;

alter table `jobs` add column `job_id` varchar(50) unique;

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

alter table `team_mems` CHANGE `member2` client_email varchar(30);
alter table `team_mems` change member3 sub_contractor varchar(30);

create table if not exists `members`(
  `id` int not null auto_increment,
  `email` varchar(100) not null,
  `team_id` int,
  `created_at` datetime not null default current_timestamp,
  `updated_at` datetime not null default current_timestamp on update current_timestamp,
  primary key(`id`),
  key `members_fk` (`team_id`),
  constraint `members_fk` foreign key (`team_id`) references `team_mems` (`id`) on update CASCADE on delete cascade
) engine=innodb AUTO_INCREMENT=1000 default charset=latin1;
  

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


  CREATE TABLE IF NOT EXISTS `reports` (
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
    `job_id` INT not null,
    `user_id` int not null,
    `report_type`
    primary key(`id`),
    key `mri_reports_relation_1` (`job_id`),
    constraint `mri_reports_relation_1` foreign key(`job_id`) references `jobs`(`id`) on delete cascade
  ) engine=innodb auto_increment=1000 default charset=latin1;

alter table `reports` drop foreign key `reports_relation_1`;
alter table `reports` drop column `job_id`;

alter table `reports` add column `order_id` int not null after `updated_at`;
alter table `reports` add constraint `mri_fk` foreign key (`order_id`) REFERENCES `orders`(`id`) on update cascade on delete cascade;

alter table `reports` add column `status` ENUM('pending', 'approved', 'rejected', 'info') default 'pending';

alter table `reports` add column `remark` VARCHAR(200);

alter table `reports` add column `user_id` int;

alter table `reports` add constraint `mki_ufk` foreign key (`user_id`) REFERENCES `users`(`id`) on update cascade;

alter table `mri_reports` rename to `reports`;

alter table `reports` add column `report_type` varchar(20);

alter table `reports` add column `file1` varchar(60);
alter table `reports` add column `file2` varchar(60);
alter table `reports` add column `file3` VARCHAR(60);
alter table `reports` add column `file4` varchar(60);

alter table `reports` add column `insp_eng_sign` varchar(50);


CREATE table IF NOT EXISTS `comments`(
  `id` int not null auto_increment,
  `approve_comment` varchar(50),
  `reject_comment` varchar(50),
  `report_id` int,
  `created_at` datetime not null default current_timestamp,
  `updated_at` datetime not null default current_timestamp on update current_timestamp,
  primary key(`id`),
  key `comment_fk`(`mri_report_id`),
  constraint `comment_fk` foreign key (`mri_report_id`) REFERENCES `mri_reports`(`id`) on update cascade
) ENGINE=Innodb AUTO_INCREMENT=1000 default charset=latin1;

alter table `comments` drop foreign key `comment_fk`;
alter table `comments` CHANGE `mri_report_id` `report_id` int;

alter table `comments` add constraint `comment_fk` foreign key (`report_id`) REFERENCES `reports` (`id`) on update CASCADE;

create table IF NOT EXISTS `client_reports`(
  `id` int not null auto_increment,
  `is_specification` TINYINT(1) not null DEFAULT '0',
  `comment` varchar(200),
  `signature` varchar(100),
  `signing_date` datetime not null default current_timestamp,
  `created_at` datetime not null default current_timestamp,
  `updated_at` datetime not null default current_timestamp on update current_timestamp,
  `report_id` int,
  primary key(`id`),
  key `client_fk`(`report_id`),
  constraint `client_fk` foreign key (`report_id`) references `reports`(`id`) on update CASCADE
) engine=innodb AUTO_INCREMENT=1000 default charset=latin1;

alter table `client_reports` add COLUMN name varchar(20);
alter table `client_reports` drop key `client_fk`;
alter table `client_reports` drop foreign key `client_fk`;
alter table `client_reports` add constraint `client_fk` foreign key `report_id` REFERENCES `reports` (`id`) on update cascade on delete cascade;

alter table `client_reports` add column `client_eng_sign` varchar(100);
alter table `client_reports` drop column `client_insp_sign`;

alter table `client_reports` add column `client_name` varchar(40);

alter table `client_reports` add column `client_sign_date` datetime not null DEFAULT CURRENT_TIMESTAMP; 


