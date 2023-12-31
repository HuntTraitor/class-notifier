CREATE TABLE IF NOT EXISTS notifications (
    email CHAR(256) NOT NULL,
    classid INT NOT NULL,
    expires TIMESTAMP NOT NULL,
    PRIMARY KEY (email, classid),
    FOREIGN KEY (email) REFERENCES users(email),
    FOREIGN KEY (classid) REFERENCES classes(classid)
);