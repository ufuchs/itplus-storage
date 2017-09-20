CREATE DATABASE `itp_home` DEFAULT CHARACTER SET utf8 COLLATE utf8_bin;

CREATE USER 'itplus'@'%' IDENTIFIED BY 'abc12345';

CREATE TABLE `fhem`.`history` (
    TIMESTAMP TIMESTAMP,
    DEVICE varchar(64),
    TYPE varchar(64),
    EVENT varchar(512),
    READING varchar(64),
    VALUE varchar(128),
    UNIT varchar(32)
);

CREATE TABLE `fhem`.`current` (
    TIMESTAMP TIMESTAMP,
    DEVICE varchar(64),
    TYPE varchar(64),
    EVENT varchar(512),
    READING varchar(64),
    VALUE varchar(128),
    UNIT varchar(32)
);

--
--
--
CREATE TABLE `itp_home`.`hub` (
    HubID       int NOT NULL AUTO_INCREMENT,
    Hostname    varchar(256),
    Alias       varchar(32),
    PRIMARY KEY (HubID)
);

--
--
--
CREATE TABLE `itp_home`.`gateway` (
    GatewayID   int NOT NULL AUTO_INCREMENT,
    HubID       int NOT NULL,
    GatewayType varchar(96),
    Hostname    varchar(32),
    Alias       varchar(32),
    PRIMARY KEY (GatewayID),
    CONSTRAINT FK_Gateway_Hub FOREIGN KEY (HubID) REFERENCES hub(HubID)
);

--
--
--
CREATE TABLE `itp-home`.`measurements` (
	MeasureID int NOT NULL AUTO_INCREMENT,
	GatewayID int NOT NULL,
    Num            int, 
    Alias          varchar(32),
    PhenomenonTime bigint,   
    Lon            float(18), 
    Lat            float(18),  
    Alt            float(18),  
    Temp           float(18),  
    Pressure       float(18),  
    Humidity       float(18),  
    LowBattery     boolean,    
    PRIMARY KEY (MeasureID),
    CONSTRAINT FK_Measurements_Gateway FOREIGN KEY (GatewayID) REFERENCES gateway(GatewayID)
);

drop table `itp-home`.`measurements`;


GRANT ALL ON `itp_home`.* TO 'itplus'@'%';
GRANT ALL ON `itp_home`.* TO 'itplus'@'192.168.178.10';

CREATE INDEX Search_Idx ON `fhem`.`history` (DEVICE, READING, TIMESTAMP);

