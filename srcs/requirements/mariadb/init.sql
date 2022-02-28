CREATE TABLE avito.`user` (
	id INT UNSIGNED AUTO_INCREMENT NOT NULL,
	PRIMARY KEY(id)
);

CREATE TABLE avito.`transaction` (
	id INT UNSIGNED AUTO_INCREMENT NOT NULL,
	user_id INT UNSIGNED NOT NULL,
	amount NUMERIC(15,2) NOT NULL,
	status INT NOT NULL,
	`from` varchar(100) NOT NULL,
	from_id INT UNSIGNED NOT NULL,
	comment varchar(1000) NOT NULL,
	PRIMARY KEY(id),
	FOREIGN KEY (user_id) REFERENCES `user`(id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;

CREATE TABLE avito.`balance` (
	id INT UNSIGNED AUTO_INCREMENT NOT NULL,
	user_id INT UNSIGNED NOT NULL,
	balance NUMERIC(15,2) UNSIGNED NOT NULL,
	PRIMARY KEY(id),
	FOREIGN KEY (user_id) REFERENCES `user`(id)
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;

INSERT INTO avito.`user` () VALUES
	 (),(),(),(),(),();

INSERT INTO avito.balance (user_id,balance) VALUES
	 (4,207.66),
	 (5,1000000000.46);

INSERT INTO avito.`transaction` (user_id,amount,status,`from`,from_id,comment) VALUES
	 (4,1.00,0,'balance',1,''),
	 (5,1000000000.23,1,'create',1,''),
	 (5,0.23,1,'balance',4,''),
	 (5,0.23,1,'balance',4,''),
	 (4,100.00,1,'bonus',0,''),
	 (4,100.00,1,'bonus',0,'спасибо что вы с нами');
