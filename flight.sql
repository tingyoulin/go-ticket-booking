CREATE DATABASE IF NOT EXISTS `ticket_booking` /*!40100 DEFAULT CHARACTER SET utf8 COLLATE utf8_unicode_ci */;
USE `ticket_booking`;

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `passenger`
--

DROP TABLE IF EXISTS `passengers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `passengers` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(200) COLLATE utf8_unicode_ci NOT NULL,
  `email` varchar(200) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `passengers`
--

LOCK TABLES `passengers` WRITE;
/*!40000 ALTER TABLE `passengers` DISABLE KEYS */;
INSERT INTO `passengers` VALUES 
(1,'Jaden Dawson','jaden.dawson@example.com','2023-11-01 10:00:00','2023-11-01 10:00:00'),
(2,'Katherine Meyer','katherine.meyer@example.com','2023-11-01 10:00:00','2023-11-01 10:00:00');
/*!40000 ALTER TABLE `passengers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `flight`
--

DROP TABLE IF EXISTS `flights`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `flights` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `departure` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `destination` varchar(100) COLLATE utf8_unicode_ci NOT NULL,
  `departure_time` datetime NOT NULL,
  `available_seats` int(11) NOT NULL,
  `price` float NOT NULL,
  `status` varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `lock_version` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE INDEX idx_flights_route_time_available_seats ON flights (departure_time, departure, destination, available_seats);

--
-- Dumping data for table `flights`
--

LOCK TABLES `flights` WRITE;
/*!40000 ALTER TABLE `flights` DISABLE KEYS */;
INSERT INTO `flights` VALUES 
(1, 'New York', 'Los Angeles', '2023-12-01 08:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(2, 'Chicago', 'Miami', '2023-12-02 09:30:00', 180, 199.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(3, 'San Francisco', 'Seattle', '2023-12-03 07:45:00', 150, 149.99, 'delayed', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(4, 'Dallas', 'Houston', '2023-12-04 11:15:00', 120, 89.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(5, 'Boston', 'Washington', '2023-12-05 14:00:00', 160, 179.99, 'canceled', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(6, 'Las Vegas', 'Phoenix', '2023-12-06 16:30:00', 200, 249.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(7, 'Atlanta', 'Orlando', '2023-12-07 12:00:00', 220, 159.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(8, 'Philadelphia', 'Newark', '2023-12-08 10:15:00', 140, 99.99, 'delayed', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(9, 'Denver', 'Salt Lake City', '2023-12-09 13:45:00', 160, 199.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(10, 'New York', 'Los Angeles', '2023-12-01 08:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(11, 'New York', 'Los Angeles', '2023-12-02 09:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(12, 'New York', 'Los Angeles', '2023-12-03 10:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(13, 'New York', 'Los Angeles', '2023-12-04 11:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(14, 'New York', 'Los Angeles', '2023-12-05 12:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(15, 'New York', 'Los Angeles', '2023-12-06 13:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(16, 'New York', 'Los Angeles', '2023-12-07 14:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(17, 'New York', 'Los Angeles', '2023-12-08 15:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(18, 'New York', 'Los Angeles', '2023-12-09 16:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(19, 'New York', 'Los Angeles', '2023-12-10 17:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(20, 'New York', 'Los Angeles', '2023-12-11 18:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(21, 'New York', 'Los Angeles', '2023-12-12 19:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(22, 'New York', 'Los Angeles', '2023-12-13 20:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(23, 'New York', 'Los Angeles', '2023-12-14 21:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(24, 'New York', 'Los Angeles', '2023-12-15 22:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(25, 'New York', 'Los Angeles', '2023-12-16 23:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(26, 'New York', 'Los Angeles', '2023-12-17 07:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(27, 'New York', 'Los Angeles', '2023-12-18 08:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(28, 'New York', 'Los Angeles', '2023-12-19 09:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0),
(29, 'New York', 'Los Angeles', '2023-12-20 10:00:00', 200, 299.99, 'on_time', '2023-11-01 10:00:00', '2023-11-01 10:00:00', 0);
/*!40000 ALTER TABLE `flights` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `booking`
--

DROP TABLE IF EXISTS `bookings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `bookings` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `passenger_id` int(11) unsigned NOT NULL,
  `flight_id` int(11) unsigned NOT NULL,
  `seats` int(11) NOT NULL,
  `status` varchar(20) COLLATE utf8_unicode_ci NOT NULL,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  FOREIGN KEY (`passenger_id`) REFERENCES `passengers`(`id`) ON DELETE CASCADE,
  FOREIGN KEY (`flight_id`) REFERENCES `flights`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE INDEX idx_bookings_flight_status_passenger_id ON bookings (flight_id, passenger_id, status);

/*!40101 SET character_set_client = @saved_cs_client */;
