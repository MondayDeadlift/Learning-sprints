CREATE TABLE words (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    text TEXT NOT NULL,
    color TEXT NOT NULL, -- hex-код цвета (#FF0000)
    user_ip TEXT, -- или user_id если будет авторизация
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    word_id INTEGER,
    voter_ip TEXT,
    FOREIGN KEY (word_id) REFERENCES words(id)
);