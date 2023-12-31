CREATE TABLE IF NOT EXISTS classes (
    classid INT PRIMARY KEY,
    name CHAR(256) NOT NULL,
    link CHAR(256) NOT NULL,
    professor CHAR(256) NOT NULL
);