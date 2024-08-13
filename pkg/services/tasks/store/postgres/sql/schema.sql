CREATE TABLE IF NOT EXISTS tasks (
  _id             BIGSERIAL PRIMARY KEY,
  url             text      NOT NULL,
  method          text      NOT NULL,
  namespace       text      NOT NULL,
  params          json,
  headers         json,
  body            json,
  start_unix      bigint   NOT NULL CHECK (start_unix > 0),
  end_unix        bigint   NOT NULL CHECK (end_unix >= 0),
  interval        text     NOT NULL
);