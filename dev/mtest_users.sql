CREATE TABLE users
(
    id INT(11) PRIMARY KEY NOT NULL AUTO_INCREMENT,
    login VARCHAR(50) NOT NULL,
    pass VARCHAR(300) NOT NULL,
    name VARCHAR(300) NOT NULL,
    last_access INT(11)
);
#CREATE UNIQUE INDEX users_id_uindex ON users (id);
#CREATE UNIQUE INDEX users_login_uindex ON users (login);