USE master;

CREATE TABLE IF NOT EXISTS geo_location (
    ip_address VARCHAR(20),
    country_code VARCHAR(5),
    country VARCHAR(50),
    city VARCHAR(30),
    latitude DECIMAL(30, 15),
    longitude	DECIMAL(30, 15),
    mystery_value DECIMAL(30, 15)
);