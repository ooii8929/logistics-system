CREATE TABLE IF NOT EXISTS locations (
    location_id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    title TEXT NOT NULL,
    city TEXT NOT NULL,
    address TEXT NOT NULL,
    created_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS recipients (
    id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    phone TEXT NOT NULL,
    created_dt DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS trackings (
    sno INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    tracking_status TEXT NOT NULL,
    estimated_delivery DATETIME NOT NULL,
    recipient_id INT UNSIGNED,
    current_location INT UNSIGNED,
    created_dt DATETIME NOT NULL,
    updated_dt DATETIME NOT NULL
);

CREATE TABLE IF NOT EXISTS tracking_histories (
    record_id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    tracking_status TEXT NOT NULL,
    sno INT UNSIGNED,
    location_id INT UNSIGNED,
    created_dt DATETIME NOT NULL
);

INSERT INTO locations (location_id, title, city, address, created_dt)
VALUES
(7, '台北物流中⼼', '台北市', '台北市中正區忠孝東路100號', NOW()),
(13, '新⽵物流中⼼', '新⽵市', '新⽵市東區光復路⼀段101號', NOW()),
(24, '台中物流中⼼', '台中市', '台中市⻄區⺠⽣路200號', NOW()),
(3, '桃園物流中⼼', '桃園市', '桃園市中壢區中央⻄路三段150號', NOW()),
(18, '⾼雄物流中⼼', '⾼雄市', '⾼雄市前⾦區成功⼀路82號', NOW()),
(9, '彰化物流中⼼', '彰化市', '彰化市中⼭路⼆段250號', NOW()),
(15, '嘉義物流中⼼', '嘉義市', '嘉義市東區⺠族路380號', NOW()),
(6, '宜蘭物流中⼼', '宜蘭市', '宜蘭市中⼭路⼆段58號', NOW()),
(21, '屏東物流中⼼', '屏東市', '屏東市⺠⽣路300號', NOW()),
(1, '花蓮物流中⼼', '花蓮市', '花蓮市國聯⼀路100號', NOW()),
(4, '台南物流中⼼', '台南市', '台南市安平區建平路18號', NOW()),
(11, '南投物流中⼼', '南投市', '南投市⾃由路67號', NOW()),
(23, '雲林物流中⼼', '雲林市', '雲林市中正路五段120號', NOW()),
(14, '基隆物流中⼼', '基隆市', '基隆市信⼀路50號', NOW()),
(8, '澎湖物流中⼼', '澎湖縣', '澎湖縣⾺公市中正路200號', NOW()),
(19, '⾦⾨物流中⼼', '⾦⾨縣', '⾦⾨縣⾦城鎮⺠⽣路90號', NOW());

INSERT INTO recipients (id, name, address, phone, created_dt)
VALUES
(1234, '賴⼩賴', '台北市中正區仁愛路⼆段99號', '091234567', NOW()),
(1235, '陳⼤明', '新北市板橋區⽂化路⼀段100號', '092345678', NOW()),
(1236, '林⼩芳', '台中市⻄區⺠⽣路200號', '093456789', NOW()),
(1237, '張美玲', '⾼雄市前⾦區成功⼀路82號', '094567890', NOW()),
(1238, '王⼩明', '台南市安平區建平路18號', '095678901', NOW()),
(1239, '劉⼤華', '新⽵市東區光復路⼀段101號', '096789012', NOW()),
(1240, '⿈⼩琳', '彰化市中⼭路⼆段250號', '097890123', NOW()),
(1241, '吳美美', '花蓮市國聯⼀路100號', '098901234', NOW()),
(1242, '蔡⼩虎', '屏東市⺠⽣路300號', '099012345', NOW()),
(1243, '鄭⼤勇', '基隆市信⼀路50號', '091123456', NOW()),
(1244, '謝⼩珍', '嘉義市東區⺠族路380號', '092234567', NOW()),
(1245, '潘⼤為', '宜蘭市中⼭路⼆段58號', '093345678', NOW()),
(1246, '趙⼩梅', '南投市⾃由路67號', '094456789', NOW()),
(1247, '周⼩⿓', '雲林市中正路五段120號', '095567890', NOW()),
(1248, '李⼤同', '澎湖縣⾺公市中正路200號', '096678901', NOW()),
(1249, '陳⼩凡', '⾦⾨縣⾦城鎮⺠⽣路90號', '097789012', NOW()),
(1250, '楊⼤明', '台北市信義區松仁路50號', '098890123', NOW()),
(1251, '吳⼩雯', '新北市中和區景平路100號', '099901234', NOW());
