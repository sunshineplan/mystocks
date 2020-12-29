-- Initialize the database.

CREATE TABLE user (
  id INT PRIMARY KEY AUTO_INCREMENT,
  username VARCHAR(20) UNIQUE NOT NULL,
  password VARCHAR(120) NOT NULL DEFAULT '123456'
);

CREATE TABLE stock (
  id INT PRIMARY KEY AUTO_INCREMENT,
  user_id INT NOT NULL,
  idx VARCHAR(10) NOT NULL,
  code VARCHAR(20) NOT NULL
);

CREATE TABLE seq (
  user_id INT NOT NULL,
  stock_id INT NOT NULL,
  seq INT NOT NULL
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
    SET @seq := (SELECT IFNULL(MAX(seq)+1, 1) FROM seq WHERE user_id = new.user_id);
    INSERT INTO seq
      (user_id, stock_id, seq)
    VALUES
      (new.user_id, new.id, @seq);
END;;

CREATE TRIGGER reorder AFTER DELETE ON stock
FOR EACH ROW BEGIN
    SET @seq := (SELECT seq FROM seq WHERE user_id = old.user_id AND stock_id = old.id);
    DELETE FROM seq
    WHERE user_id = old.user_id AND seq = @seq;
    UPDATE seq SET seq = seq-1
    WHERE user_id = old.user_id AND seq > @seq;
END;;
DELIMITER ;

SET SESSION sql_mode = 'NO_AUTO_VALUE_ON_ZERO';
INSERT INTO user (id, username, password) VALUES (0, 'guest', '');
