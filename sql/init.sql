DROP TABLE IF EXISTS quote;

CREATE TABLE quote (
	   id INT AUTO_INCREMENT,
	   quote VARCHAR(255),
	   owner VARCHAR(255),
	   date DATE,
	   PRIMARY KEY(id)
);
