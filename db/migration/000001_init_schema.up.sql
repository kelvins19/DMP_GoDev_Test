CREATE TABLE users(id SERIAL PRIMARY KEY, username VARCHAR(50) UNIQUE, password_hash VARCHAR(100), created_at TIMESTAMP, updated_at TIMESTAMP);

INSERT INTO users(username, password_hash) VALUES('kelvin', '$2a$12$X5A51EQlbOx/5LtJi3uyN...4ijjlkHJ6zcJXXY2QkxVgb.iK6M0G');