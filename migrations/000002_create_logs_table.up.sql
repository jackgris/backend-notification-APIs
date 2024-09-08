CREATE TABLE logs (
    id SERIAL PRIMARY KEY,
    user_id INT,
    category VARCHAR(50),
    message TEXT,
    notification_type VARCHAR(50),
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
