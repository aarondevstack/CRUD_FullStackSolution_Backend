-- Create "users" table
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `username` varchar(255) NOT NULL,
  `password` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `role` enum('user','admin') NOT NULL DEFAULT "user",
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `email` (`email`),
  UNIQUE INDEX `username` (`username`)
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "addresses" table
CREATE TABLE `addresses` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `street` varchar(255) NOT NULL,
  `city` varchar(255) NOT NULL,
  `state` varchar(255) NOT NULL,
  `zip` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  `user_addresses` bigint NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `addresses_users_addresses` (`user_addresses`),
  CONSTRAINT `addresses_users_addresses` FOREIGN KEY (`user_addresses`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "blogs" table
CREATE TABLE `blogs` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL,
  `content` longtext NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  `user_blogs` bigint NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `blogs_users_blogs` (`user_blogs`),
  CONSTRAINT `blogs_users_blogs` FOREIGN KEY (`user_blogs`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
-- Create "comments" table
CREATE TABLE `comments` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `content` longtext NOT NULL,
  `created_at` timestamp NOT NULL,
  `updated_at` timestamp NOT NULL,
  `blog_comments` bigint NOT NULL,
  `user_comments` bigint NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `comments_blogs_comments` (`blog_comments`),
  INDEX `comments_users_comments` (`user_comments`),
  CONSTRAINT `comments_blogs_comments` FOREIGN KEY (`blog_comments`) REFERENCES `blogs` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `comments_users_comments` FOREIGN KEY (`user_comments`) REFERENCES `users` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_bin;
