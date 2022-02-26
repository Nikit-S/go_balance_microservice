CREATE TABLE avito.`transaction` (
	id INT UNSIGNED NOT NULL,
	user_id INT UNSIGNED NOT NULL,
	amount decimal(15,2) NOT NULL,
	status INT NOT NULL,
	`from` varchar(100) NOT NULL,
	from_id INT UNSIGNED NOT NULL,
	comment varchar(1000) NOT NULL
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;
CREATE INDEX transaction_id_IDX USING BTREE ON avito.`transaction` (id);
ALTER TABLE avito.`transaction` MODIFY COLUMN id int(10) unsigned auto_increment NOT NULL;

CREATE TABLE avito.`balance` (
	id INT UNSIGNED NOT NULL,
	user_id INT UNSIGNED NOT NULL,
	balance decimal(15,2) UNSIGNED NOT NULL
)
ENGINE=InnoDB
DEFAULT CHARSET=utf8
COLLATE=utf8_general_ci;
CREATE INDEX balance_id_IDX USING BTREE ON avito.balance (id);
ALTER TABLE avito.balance MODIFY COLUMN id int(10) unsigned auto_increment NOT NULL;
ALTER TABLE avito.balance ADD CONSTRAINT balance_UN UNIQUE KEY (user_id);

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
