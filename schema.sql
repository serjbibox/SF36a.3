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
INSERT INTO news_hash(link, news_hash, pub_time)
VALUES (
    'test link 1',
    'testhash1',
    1
),
(
    'test link 2',
    'testhash2',
    2
);

INSERT INTO news(title, content, pub_time, link)
VALUES (
    'title1',
    'content1',
    1,
    'link1'
),
(
    'title2',
    'content2',
    2,
    'link2'
);