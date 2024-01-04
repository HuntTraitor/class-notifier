CREATE TABLE IF NOT EXISTS notifications (
    notificationid SERIAL PRIMARY KEY,
    email CHAR(256) NOT NULL,
    classid INT NOT NULL,
    expires TIMESTAMP NOT NULL,
    CONSTRAINT unique_notification UNIQUE(email, classid),
    FOREIGN KEY (email) REFERENCES users(email),
    FOREIGN KEY (classid) REFERENCES classes(classid)
);