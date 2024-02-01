CREATE TABLE IF NOT EXISTS users (
    userid SERIAL PRIMARY KEY,
    name CHAR(256) NOT NULL,
    email CHAR(256) NOT NULL UNIQUE,
    hashed_password CHAR(60) NOT NULL,
    created TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS classes (
    classid INT PRIMARY KEY,
    name CHAR(256) NOT NULL,
    link CHAR(256) NOT NULL,
    professor CHAR(256) NOT NULL
);

CREATE TABLE IF NOT EXISTS notifications (
    notificationid SERIAL PRIMARY KEY,
    email CHAR(256) NOT NULL,
    classid INT NOT NULL,
    expires TIMESTAMP NOT NULL,
    CONSTRAINT unique_notification UNIQUE(email, classid),
    FOREIGN KEY (email) REFERENCES users(email) ON DELETE CASCADE,
    FOREIGN KEY (classid) REFERENCES classes(classid) ON DELETE CASCADE
);

INSERT INTO users (name, email, hashed_password, created) VALUES (
    'Hunter Tratar',
    'hunter@gmail.com',
    '$2a$12$NuTjWXm3KKntReFwyBVHyuf/to.HEwTy.eS206TNfkGfr6HzGJSWG',
    '2022-01-01 10:23:56'
);

INSERT INTO classes (classid, name, link, professor) VALUES (
    1,
    'Test Class',
    'testclass.com',
    'Professor Test'
);

INSERT INTO classes (classid, name, link, professor) VALUES (
    2,
    'Test Class 2',
    'testclass2.com',
    'Professor Test 2'
);