DROP TABLE IF EXISTS news;
CREATE TABLE news (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    pub_time INTEGER DEFAULT 0,
    link TEXT NOT NULL UNIQUE
);
DROP TABLE IF EXISTS news_hash;
CREATE TABLE news_hash (
    id SERIAL PRIMARY KEY,
    news_hash TEXT NOT NULL,
    pub_time INTEGER DEFAULT 0,
    link TEXT NOT NULL UNIQUE
);