-- Initialize the database.

CREATE TABLE user (
  id INT PRIMARY KEY AUTOINCREMENT,
  username VARCHAR(20) UNIQUE NOT NULL,
  password VARCHAR(120) NOT NULL DEFAULT '123456'
);

CREATE TABLE stock (
  user_id INT NOT NULL,
  idx VARCHAR(10) NOT NULL,
  code VARCHAR(20) NOT NULL,
  seq INT DEFAULT 0,
);

DELIMITER ;;
CREATE TRIGGER add_user AFTER INSERT ON user
FOR EACH ROW BEGIN
    INSERT INTO stock (user_id, idx, code)
    VALUES
      (new.id, 'SSE', '000001'),
      (new.id, 'SZSE', '399001'),
      (new.id, 'SZSE', '399106'),
      (new.id, 'SZSE', '399005'),
      (new.id, 'SZSE', '399006');
END;;

CREATE TRIGGER add_seq AFTER INSERT ON stock
FOR EACH ROW BEGIN
    UPDATE stock SET seq = (SELECT MAX(seq) + 1 FROM stock WHERE user_id = new.user_id)
    WHERE user_id = new.user_id AND idx = new.idx AND code = new.code;
END;;

CREATE TRIGGER reorder AFTER DELETE ON stock
FOR EACH ROW BEGIN
    UPDATE stock SET seq = seq - 1
    WHERE user_id = old.user_id AND seq > old.seq;
END;;
DELIMITER ;

INSERT INTO user (id, username, password)
VALUES (0, 'guest', '');
