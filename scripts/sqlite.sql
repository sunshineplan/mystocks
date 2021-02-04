-- Initialize the database.

CREATE TABLE user (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  username TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL DEFAULT '123456',
  uid TEXT UNIQUE NOT NULL DEFAULT ''
);

CREATE TABLE stock (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL,
  idx TEXT NOT NULL,
  code TEXT NOT NULL
);

CREATE TABLE seq (
  user_id INTEGER NOT NULL,
  stock_id INTEGER NOT NULL,
  seq INTEGER NOT NULL
);

CREATE TRIGGER add_user AFTER INSERT ON user
BEGIN
    INSERT INTO stock (user_id, idx, code)
    VALUES
      (new.id, 'SSE', '000001'),
      (new.id, 'SZSE', '399001'),
      (new.id, 'SZSE', '399106'),
      (new.id, 'SZSE', '399005'),
      (new.id, 'SZSE', '399006');
END;

CREATE TRIGGER add_seq AFTER INSERT ON stock
BEGIN
    INSERT INTO seq
      (user_id, stock_id, seq)
    VALUES
      (new.user_id, new.id, (SELECT IFNULL(MAX(seq)+1, 1) FROM seq WHERE user_id = new.user_id));
END;

CREATE TRIGGER reorder AFTER DELETE ON stock
BEGIN
    DELETE FROM seq
    WHERE user_id = old.user_id AND seq = (SELECT seq FROM seq WHERE user_id = old.user_id AND stock_id = old.id);
    UPDATE seq SET seq = seq-1
    WHERE user_id = old.user_id AND seq > (SELECT seq FROM seq WHERE user_id = old.user_id AND stock_id = old.id);
END;

INSERT INTO user (id, username, password)
VALUES (0, 'guest', '');
