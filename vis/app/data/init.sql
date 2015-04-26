
DROP TABLE IF EXISTS users;
CREATE TABLE users (
	id INTEGER PRIMARY KEY NOT NULL,
	username VARCHAR(255) UNIQUE NOT NULL,
	password VARCHAR(255) NOT NULL,
	name VARCHAR(255),
	email VARCHAR(255),
	role VARCHAR(255)
	);

INSERT INTO users (id, username, password, name, email, role) VALUES (1, 'tony', 'stark', 'Tony Stark', 'tony@stark.com', 'admin');

DROP TABLE IF EXISTS metrics;
CREATE TABLE metrics (
	id INTEGER PRIMARY KEY NOT NULL,
	name VARCHAR(255) UNIQUE NOT NULL,
	type VARCHAR(255) NOT NULL,
	detector TEXT NOT NULL,
	params TEXT,
	md5 VARCHAR(255)
	);

INSERT INTO metrics (id, name, type, detector, params, md5) VALUES (1, "Load", "call", "Load", "", "");
INSERT INTO metrics (id, name, type, detector, params, md5) VALUES (2, "CPURate", "call", "CPURate", "", "");
INSERT INTO metrics (id, name, type, detector, params, md5) VALUES (3, "MemRate", "call", "MemRate", "", "");
INSERT INTO metrics (id, name, type, detector, params, md5) VALUES (4, "DiskRate", "call", "DiskRate", "", "");
INSERT INTO metrics (id, name, type, detector, params, md5) VALUES (5, "DiskRead", "call", "DiskRead", "", "");
INSERT INTO metrics (id, name, type, detector, params, md5) VALUES (6, "DiskWrite", "call", "DiskWrite", "", "");
INSERT INTO metrics (id, name, type, detector, params, md5) VALUES (7, "NetRead", "call", "NetRead", "", "");
INSERT INTO metrics (id, name, type, detector, params, md5) VALUES (8, "NetWrite", "call", "NetWrite", "", "");

DROP TABLE IF EXISTS default_metrics;
CREATE TABLE default_metrics (
	id INTEGER PRIMARY KEY NOT NULL,
	interval INTEGER NOT NULL
	);

INSERT INTO default_metrics (id, interval) VALUES (1, 600);
INSERT INTO default_metrics (id, interval) VALUES (2, 600);
INSERT INTO default_metrics (id, interval) VALUES (3, 600);
INSERT INTO default_metrics (id, interval) VALUES (4, 600);
INSERT INTO default_metrics (id, interval) VALUES (5, 600);
INSERT INTO default_metrics (id, interval) VALUES (6, 600);
INSERT INTO default_metrics (id, interval) VALUES (7, 600);
INSERT INTO default_metrics (id, interval) VALUES (8, 600);

DROP TABLE IF EXISTS nodes;
CREATE TABLE nodes (
	id INTEGER PRIMARY KEY NOT NULL,
	addr VARCHAR(255) UNIQUE NOT NULL,
	type VARCHAR(255),
	name VARCHAR(255) UNIQUE,
	os  VARCHAR(255),
	cpu VARCHAR(255),
	core VARCHAR(255),
	mem VARCHAR(255),
	disk VARCHAR(255),
	uptime VARCHAR(255),
	ctime DATETIME,
	atime DATETIME
	);

DROP TABLE IF EXISTS metric_bindings;
CREATE TABLE metric_bindings (
	id INTEGER PRIMARY KEY NOT NULL,
	node INTEGER NOT NULL,
	metric INTEGER NOT NULL,
	interval INTEGER NOT NULL,
	ctime DATETIME,
	atime DATETIME,
	UNIQUE(node, metric)
	);

DROP TABLE IF EXISTS metric_records;
CREATE TABLE metric_records (
	id INTEGER PRIMARY KEY NOT NULL,
	node INTEGER NOT NULL,
	metric INTEGER NOT NULL,
	value TEXT NOT NULL,
	ctime DATETIME,
	UNIQUE(node, metric)
	);

DROP TABLE IF EXISTS config;
CREATE TABLE config (
	id INTEGER PRIMARY KEY NOT NULL,
	name VARCHAR(255) UNIQUE NOT NULL,
	value TEXT
	);
