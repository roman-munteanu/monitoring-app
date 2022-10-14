-- monitoringdb
CREATE DATABASE monitoringdb;

USE monitoringdb;

CREATE TABLE heroes (
    id INT UNSIGNED NOT NULL AUTO_INCREMENT,
    name VARCHAR(100) NOT NULL,
    category VARCHAR(50),
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8; 

INSERT INTO heroes (name, category) VALUES ('Gelu', 'Archer'), ('Crag Hack', 'Barbarian'), ('Solmyr', 'Wizard');

SELECT * FROM heroes;
