CREATE TABLE IF NOT EXISTS tasks (
  _id         BIGSERIAL PRIMARY KEY,
  url        text      NOT NULL,
  method     text      NOT NULL,
  namespace  text      NOT NULL,
  params     json,
  headers    json,
  body       json
);