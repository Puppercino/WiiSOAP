/*

< Database Template File >
This file automatically adds the default database and tables.
WiiSOAP uses MySQL.

This SQL File does not guarantee functionality as WiiSOAP is still in early development statements.
It is suggested that you should hold off from using WiiSOAP unless you are confident that you know what you are doing.
Follow and practice proper security practices before handling user data.

*/

-- Generation Time: Jan 23, 2019 at 12:40 PM
-- Server version: 10.1.37-MariaDB
-- PHP Version: 7.3.0

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
SET AUTOCOMMIT = 0;
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `wiisoap`
--

-- --------------------------------------------------------

--
-- Table structure for table `userbase`
--

CREATE TABLE `userbase` (
    `DeviceId` varchar(10) NOT NULL,
    `DeviceToken` varchar(64) NOT NULL COMMENT 'This token should be considered a secret, so after generation only the sha256sum of the md5 the Wii sends is inserted.',
    `AccountId` varchar(9) NOT NULL,
    `Region` varchar(2) NOT NULL,
    `Country` varchar(2) NOT NULL,
    `Language` varchar(2) NOT NULL,
    `SerialNo` varchar(11) NOT NULL,
    `DeviceCode` varchar(16) NOT NULL,
    PRIMARY KEY (`AccountId`),
    UNIQUE KEY `AccountId` (`AccountId`),
    UNIQUE KEY `userbase_DeviceId_uindex` (`DeviceId`),
    UNIQUE KEY `userbase_DeviceToken_uindex` (`DeviceToken`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;