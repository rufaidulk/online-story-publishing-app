-- phpMyAdmin SQL Dump
-- version 5.1.1
-- https://www.phpmyadmin.net/
--
-- Host: subscription_service_db
-- Generation Time: Feb 13, 2022 at 09:21 AM
-- Server version: 8.0.27
-- PHP Version: 7.4.20

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `subscription_service`
--

-- --------------------------------------------------------

--
-- Table structure for table `payment_transactions`
--

DROP TABLE IF EXISTS `payment_transactions`;
CREATE TABLE `payment_transactions` (
  `id` bigint NOT NULL,
  `ref_no` varchar(255) NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `type` enum('CREDIT','DEBIT') NOT NULL,
  `txn_category` varchar(255) NOT NULL,
  `status` tinyint NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `sequences`
--

DROP TABLE IF EXISTS `sequences`;
CREATE TABLE `sequences` (
  `id` int NOT NULL,
  `type` varchar(255) NOT NULL,
  `seq_no` bigint NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `sequences`
--

INSERT INTO `sequences` (`id`, `type`, `seq_no`) VALUES
(1, 'subscriptionTxn', 0);

-- --------------------------------------------------------

--
-- Table structure for table `subscriptions`
--

DROP TABLE IF EXISTS `subscriptions`;
CREATE TABLE `subscriptions` (
  `id` bigint NOT NULL,
  `user_uuid` binary(16) NOT NULL,
  `user_uuid_text` varchar(40) GENERATED ALWAYS AS (bin_to_uuid(`user_uuid`)) VIRTUAL NOT NULL,
  `subscription_plan_id` int NOT NULL,
  `payment_transaction_id` bigint NOT NULL,
  `is_premium` tinyint NOT NULL,
  `expiry_date` date DEFAULT NULL,
  `status` tinyint NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- --------------------------------------------------------

--
-- Table structure for table `subscription_plans`
--

DROP TABLE IF EXISTS `subscription_plans`;
CREATE TABLE `subscription_plans` (
  `id` int NOT NULL,
  `period_type` tinyint NOT NULL,
  `name` varchar(255) NOT NULL,
  `plan_type` tinyint NOT NULL,
  `amount` decimal(10,2) NOT NULL,
  `period_in_days` int NOT NULL,
  `status` tinyint NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

--
-- Dumping data for table `subscription_plans`
--

INSERT INTO `subscription_plans` (`id`, `period_type`, `name`, `plan_type`, `amount`, `period_in_days`, `status`, `created_at`, `updated_at`) VALUES
(1, 10, 'BasicFree', 11, '0.00', 0, 10, '2022-02-13 08:52:10', '2022-02-13 08:52:10'),
(2, 10, 'Premium30', 10, '100.00', 30, 10, '2022-02-13 08:52:10', '2022-02-13 08:52:10'),
(3, 11, 'Premium365', 10, '1200.00', 365, 10, '2022-02-13 08:52:10', '2022-02-13 08:52:10');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `payment_transactions`
--
ALTER TABLE `payment_transactions`
  ADD PRIMARY KEY (`id`),
  ADD UNIQUE KEY `ref_no` (`ref_no`);

--
-- Indexes for table `sequences`
--
ALTER TABLE `sequences`
  ADD UNIQUE KEY `type` (`type`);

--
-- Indexes for table `subscriptions`
--
ALTER TABLE `subscriptions`
  ADD PRIMARY KEY (`id`),
  ADD KEY `user_uuid` (`user_uuid`),
  ADD KEY `subscription_plan_id` (`subscription_plan_id`),
  ADD KEY `payment_transaction_id` (`payment_transaction_id`);

--
-- Indexes for table `subscription_plans`
--
ALTER TABLE `subscription_plans`
  ADD PRIMARY KEY (`id`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `payment_transactions`
--
ALTER TABLE `payment_transactions`
  MODIFY `id` bigint NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `subscriptions`
--
ALTER TABLE `subscriptions`
  MODIFY `id` bigint NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `subscription_plans`
--
ALTER TABLE `subscription_plans`
  MODIFY `id` int NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=4;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;