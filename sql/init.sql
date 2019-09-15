DROP TABLE IF EXISTS quote;

CREATE TABLE quote (
	   id INT AUTO_INCREMENT UNIQUE,
	   quote_name VARCHAR(64),
	   quote VARCHAR(255),
	   owner VARCHAR(255),
	   date VARCHAR(16),
	   category VARCHAR(64),
	   PRIMARY KEY(quote_name)
);
