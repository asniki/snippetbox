# snippetbox

A web application, which lets paste and share snippets of text

![snippetbox web app](/img/homepage.png "snippetbox web app")

###  MySQL migration

    -- Create a new UTF-8 `snippetbox` database
    CREATE DATABASE snippetbox CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

    -- Create a new web user with SELECT and INSERT privileges only
    CREATE USER 'web'@'localhost';
    GRANT SELECT, INSERT, UPDATE, DELETE ON snippetbox.* TO 'web'@'localhost';

    -- Important: Make sure to swap 'pass' with a password of your own choosing
    ALTER USER 'web'@'localhost' IDENTIFIED BY 'pass';

    -- Switch to using the `snippetbox` database
    USE snippetbox;

    -- Create a `snippets` table
    CREATE TABLE snippets (
        id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
        title VARCHAR(100) NOT NULL,
        content TEXT NOT NULL,
        created DATETIME NOT NULL,
        expires DATETIME NOT NULL
    );

    -- Add an index on the created column
    CREATE INDEX idx_snippets_created ON snippets(created);

    -- Add some dummy records
    INSERT INTO snippets (title, content, created, expires) VALUES (
        'An old silent pond',
        'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō',
        UTC_TIMESTAMP(),
        DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
    );

    -- Create a `sessions` table for session manager
    CREATE TABLE sessions (
        token CHAR(43) PRIMARY KEY,
        data BLOB NOT NULL,
        expiry TIMESTAMP(6) NOT NULL
    );

    CREATE INDEX sessions_expiry_idx ON sessions (expiry);

    -- Create a `users` table
    CREATE TABLE users (
        id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
        name VARCHAR(255) NOT NULL,
        email VARCHAR(255) NOT NULL,
        hashed_password CHAR(60) NOT NULL,
        created DATETIME NOT NULL
    );

    ALTER TABLE users ADD CONSTRAINT users_uc_email UNIQUE (email);


### Create a self-signed certificate for localhost (for macOS)

    cd ./tls
    go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost


### Run web application

    go run ./cmd/web -addr=":4000" -dsn="user:pass@/snippetbox?parseTime=true"


### Build an executable binary

    go build -o /tmp/web ./cmd/web/
    cp -r ./tls /tmp/
    cd /tmp/
    ./web

### Open web application

[https://localhost:4000/](https://localhost:4000/)
