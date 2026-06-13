INSERT INTO artists (id, name, information, image_url) VALUES
(1,'The Weeknd','Canadian singer-songwriter',NULL),
(2,'Dua Lipa','British-Albanian singer',NULL),
(3,'Coldplay','British rock band',NULL),
(4,'Taylor Swift','American singer-songwriter',NULL),
(5,'Ed Sheeran','English singer-songwriter',NULL);

INSERT INTO genres (id, name) VALUES
(1,'Pop'),
(2,'Synth-pop'),
(3,'Rock'),
(4,'Alternative Rock');

INSERT INTO albums
(id,title,type,genre_id,cover_image_url,release_date)
VALUES
(1,'After Hours','LP',2,NULL,'2020-03-20'),
(2,'Dawn FM','LP',2,NULL,'2022-01-07'),
(3,'Future Nostalgia','LP',1,NULL,'2020-03-27'),
(4,'A Head Full of Dreams','LP',3,NULL,'2015-12-04'),
(5,'Midnights','LP',1,NULL,'2022-10-21'),
(6,'Divide','LP',1,NULL,'2017-03-03'),
(7,'Equals','LP',1,NULL,'2021-10-29');

INSERT INTO artists_albums VALUES
(1,1),
(1,2),
(2,3),
(3,4),
(4,5),
(5,6),
(5,7);

INSERT INTO songs
(id,title,file_id,genre_id,duration,bpm,release_date)
VALUES
(1,'Blinding Lights','file_001',2,200,171,'2020-03-20'),
(2,'Save Your Tears','file_002',2,215,118,'2020-03-20'),
(3,'In Your Eyes','file_003',2,238,100,'2020-03-20'),

(4,'Gasoline','file_004',2,212,123,'2022-01-07'),
(5,'Sacrifice','file_005',2,188,122,'2022-01-07'),
(6,'Take My Breath','file_006',2,221,121,'2022-01-07'),

(7,'Dont Start Now','file_007',1,183,124,'2020-03-27'),
(8,'Levitating','file_008',1,203,103,'2020-03-27'),
(9,'Physical','file_009',1,193,147,'2020-03-27'),

(10,'Adventure of a Lifetime','file_010',3,251,112,'2015-12-04'),
(11,'Hymn for the Weekend','file_011',3,258,90,'2015-12-04'),
(12,'Up and Up','file_012',3,405,82,'2015-12-04'),

(13,'Anti-Hero','file_013',1,201,97,'2022-10-21'),
(14,'Lavender Haze','file_014',1,202,98,'2022-10-21'),
(15,'Karma','file_015',1,204,90,'2022-10-21'),

(16,'Shape of You','file_016',1,233,96,'2017-03-03'),
(17,'Perfect','file_017',1,263,63,'2017-03-03'),
(18,'Castle on the Hill','file_018',1,261,135,'2017-03-03'),

(19,'Bad Habits','file_019',1,231,126,'2021-10-29'),
(20,'Shivers','file_020',1,207,141,'2021-10-29');

INSERT INTO albums_songs VALUES
(1,1,1),
(2,1,2),
(3,1,3),

(1,2,4),
(2,2,5),
(3,2,6),

(1,3,7),
(2,3,8),
(3,3,9),

(1,4,10),
(2,4,11),
(3,4,12),

(1,5,13),
(2,5,14),
(3,5,15),

(1,6,16),
(2,6,17),
(3,6,18),

(1,7,19),
(2,7,20);

INSERT INTO artists_songs VALUES

-- After Hours
(1,1,'Main Artist'),
(1,2,'Main Artist'),
(4,2,'Featured Artist'),
(1,3,'Main Artist'),

-- Dawn FM
(1,4,'Main Artist'),
(1,5,'Main Artist'),
(1,6,'Main Artist'),

-- Future Nostalgia
(2,7,'Main Artist'),
(2,8,'Main Artist'),
(5,8,'Featured Artist'),
(2,9,'Main Artist'),

-- A Head Full of Dreams
(3,10,'Main Artist'),
(3,11,'Main Artist'),
(1,11,'Featured Artist'),
(3,12,'Main Artist'),

-- Midnights
(4,13,'Main Artist'),
(4,14,'Main Artist'),
(4,15,'Main Artist'),
(2,15,'Featured Artist'),

-- Divide
(5,16,'Main Artist'),
(5,17,'Main Artist'),
(5,18,'Main Artist'),

-- Equals
(5,19,'Main Artist'),
(5,20,'Main Artist');