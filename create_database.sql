USE dwtakehome;

-- For Firstname and Surname field, we can use VARCHAR(255) as we don't know the maximum length of the name.
CREATE TABLE IF NOT EXISTS users (
    user_id INT AUTO_INCREMENT PRIMARY KEY,
    username VARCHAR(32) UNIQUE NOT NULL,
    firstname VARCHAR(255) NOT NULL,
    surname VARCHAR(255) NOT NULL,
    email VARCHAR(64) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL
);

-- For FK on user_id, we let ON DELETE CASCADE to delete all the posts of the user when the user is deleted.
-- Same manner to other social media platforms.
CREATE TABLE IF NOT EXISTS posts (
    post_id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    message TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- We need to create a table for comments and user first before create CommentUser table
-- because many-to-many relationship is created by creating a new table that contains the PK of both tables.
CREATE TABLE IF NOT EXISTS comments(
    comment_id INT AUTO_INCREMENT PRIMARY KEY,
    comment_msg TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS CommentUser(
    comment_id INT NOT NULL,
    user_id INT NOT NULL,
    PRIMARY KEY (comment_id, user_id),
    FOREIGN KEY (comment_id) REFERENCES comments(comment_id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- On the other hands, when we delete table, we need to delete a table that link to 
-- many-to-many relationship table first.
