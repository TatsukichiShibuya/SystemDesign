INSERT INTO `tasks` (`id`, `title`)           VALUES (1, "買い物");
INSERT INTO `tasks` (`id`, `title`, `detail`) VALUES (2, "レポート", "やる");
INSERT INTO `tasks` (`id`, `title`)           VALUES (3, "散歩");
INSERT INTO `tasks` (`id`, `title`, `detail`) VALUES (4, "勉強", "やるべし");

INSERT INTO `users` (`id`, `username`, `passward`) VALUES (1, "shibuya", "d74ff0ee8da3b9806b18c877dbf29bbde50b5bd8e4dad7a3a725000feb82e8f1");
INSERT INTO `users` (`id`, `username`, `passward`) VALUES (2, "tatsukichi", "d74ff0ee8da3b9806b18c877dbf29bbde50b5bd8e4dad7a3a725000feb82e8f1");

INSERT INTO `owners` (`userid`, `taskid`) VALUES (1, 1);
INSERT INTO `owners` (`userid`, `taskid`) VALUES (1, 2);
INSERT INTO `owners` (`userid`, `taskid`) VALUES (2, 3);
INSERT INTO `owners` (`userid`, `taskid`) VALUES (2, 4);
