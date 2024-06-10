CREATE TABLE IF NOT EXISTS users
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    pass_hash BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    status VARCHAR(255) NOT NULL DEFAULT 'active'
    is_admin BOOLEAN NOT NULL DEFAULT FALSE
    );
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

CREATE TABLE IF NOT EXISTS modules(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    api_key VARCHAR(255) NOT NULL,
    plant_family_id INT NOT NULL,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    status VARCHAR(255) NOT NULL DEFAULT 'active'
    FOREIGN KEY (plant_family_id) REFERENCES plantFamily (id) ON DELETE CASCADE
)

CREATE TABLE IF NOT EXISTS user_modules(
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    module_id INT NOT NULL,
    status VARCHAR(255) NOT NULL DEFAULT 'active',
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (module_id) REFERENCES modules (id) ON DELETE CASCADE
)

CREATE TABLE IF NOT EXISTS recommendations(
    id INT AUTO_INCREMENT PRIMARY KEY,
    temperature_min FLOAT NOT NULL,
    temperature_max FLOAT NOT NULL,
    humidity_in_min FLOAT NOT NULL,
    humidity_in_max FLOAT NOT NULL,
    humidity_out_min FLOAT NOT NULL,
    humidity_out_max FLOAT NOT NULL,
    illuminance_min FLOAT NOT NULL,
    illuminance_max FLOAT NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(255) NOT NULL DEFAULT 'active'
);

CREATE TABLE IF NOT EXISTS settings(
    module_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    temperature_min FLOAT NOT NULL,
    temperature_max FLOAT NOT NULL,
    humidity_in_min FLOAT NOT NULL,
    humidity_in_max FLOAT NOT NULL,
    humidity_out_min FLOAT NOT NULL,
    humidity_out_max FLOAT NOT NULL,
    illuminance_min FLOAT NOT NULL,
    illuminance_max FLOAT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (module_id) REFERENCES modules (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS plantFamily(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    recommendation_id INT,
    description TEXT,
    picture_url VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (recommendation_id) REFERENCES recommendations (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS moduleData(
    module_id INT NOT NULL,
    humidity_in FLOAT,
    humidity_out FLOAT,
    temperature FLOAT,
    illuminance FLOAT,
    measurement_time TIMESTAMP,
    FOREIGN KEY (module_id) REFERENCES modules (id) ON DELETE CASCADE
);