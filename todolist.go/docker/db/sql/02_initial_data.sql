INSERT INTO `tasks` (`title`) VALUES ("買い物");
INSERT INTO `tasks` (`title`, `detail`) VALUES ("レポート", "やる");
INSERT INTO `tasks` (`title`) VALUES ("散歩");
INSERT INTO `tasks` (`title`, `detail`) VALUES ("勉強", "やるべし");

INSERT INTO `users` (`username`, `passward`) VALUES ("shibuya", "d74ff0ee8da3b9806b18c877dbf29bbde50b5bd8e4dad7a3a725000feb82e8f1");
INSERT INTO `users` (`username`, `passward`) VALUES ("tatsukichi", "d74ff0ee8da3b9806b18c877dbf29bbde50b5bd8e4dad7a3a725000feb82e8f1");

INSERT INTO `owners` (`username`, `taskid`) VALUES ("shibuya", 1);
INSERT INTO `owners` (`username`, `taskid`) VALUES ("shibuya", 2);
INSERT INTO `owners` (`username`, `taskid`) VALUES ("tatsukichi", 3);
INSERT INTO `owners` (`username`, `taskid`) VALUES ("tatsukichi", 4);
