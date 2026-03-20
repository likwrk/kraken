CREATE TABLE sensor_readings (
    device_id VARCHAR(100) NOT NULL,
    sensor_id VARCHAR(100) NOT NULL,
    sensor_value DOUBLE NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    PRIMARY KEY (device_id, sensor_id, timestamp)
);
CREATE INDEX idx_reading_set_created ON sensor_readings(device_id, timestamp);


CREATE TABLE sensor_sets (
    device_id VARCHAR(100) NOT NULL,
    sensor_id VARCHAR(100) NOT NULL,
    sensor_value DOUBLE NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    PRIMARY KEY (device_id, sensor_id, timestamp)
);
CREATE INDEX idx_set_created ON sensor_sets(device_id, timestamp);


delete from sensor_sets where timestamp < 1773668004;
delete from sensor_readings where timestamp < 1773668004;
update sensor_readings set device_id = 'b2ee2e5c-dd93-41b6-be86-cfb1fda6c480';
update sensor_sets set device_id = 'b2ee2e5c-dd93-41b6-be86-cfb1fda6c480';
