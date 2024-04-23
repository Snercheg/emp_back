CREATE TABLE IF NOT EXISTS users
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    pass_hash BYTEA NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    status VARCHAR(255) NOT NULL DEFAULT 'active'
    role_id INT NOT NULL DEFAULT 1,
    FOREIGN KEY (role_id) REFERENCES roles (id) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);

CREATE TABLE IF NOT EXISTS apps (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    secret_key VARCHAR(255) NOT NULL UNIQUE,
);

CREATE TABLE IF NOT EXISTS roles (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
);


CREATE TABLE IF NOT EXISTS modules(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    api_key VARCHAR(255) NOT NULL,
    plant_family_id INT NOT NULL,
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    status VARCHAR(255) NOT NULL DEFAULT 'active'
    FOREIGN KEY (plant_family_id) REFERENCES plantFamily (id) ON DELETE CASCADE
)

CREATE TABLE IF NOT EXISTS modules_users(
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    module_id INT NOT NULL,
    linked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (module_id) REFERENCES modules (id) ON DELETE CASCADE
)

CREATE TABLE IF NOT EXISTS recommendations(
    id INT AUTO_INCREMENT PRIMARY KEY,
    temperature_min FLOAT NOT NULL,
    temperature_max FLOAT NOT NULL,
    humidity_min FLOAT NOT NULL,
    humidity_max FLOAT NOT NULL,
    illuminance_min FLOAT NOT NULL,
    illuminance_max FLOAT NOT NULL,
    description_care TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(255) NOT NULL DEFAULT 'active'
);

CREATE TABLE IF NOT EXISTS settings(
    module_id INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    temperature_min FLOAT NOT NULL,
    temperature_max FLOAT NOT NULL,
    humidity_min FLOAT NOT NULL,
    humidity_max FLOAT NOT NULL,
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
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (recommendation_id) REFERENCES recommendations (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS deviceData(
    module_id INT NOT NULL,
    humidity FLOAT,
    temperature FLOAT,
    illuminance FLOAT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (module_id) REFERENCES modules (id) ON DELETE CASCADE
);