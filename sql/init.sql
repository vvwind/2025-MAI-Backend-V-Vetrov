-- init.sql
-- This script will run when the PostgreSQL container starts

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_name VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    user_role VARCHAR(50) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS users_creds (
    user_id SERIAL PRIMARY KEY,
    email VARCHAR(100) REFERENCES users(email) ON DELETE CASCADE,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Create products table
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    seller_name VARCHAR(100) REFERENCES users(user_name) ON DELETE CASCADE,
    seller_id INT REFERENCES users(id) ON DELETE CASCADE,
    product_image text NOT NULL,
    product_description TEXT NOT NULL,
    price INTEGER NOT NULL,
    amount INTEGER NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- Optionally insert some test data
INSERT INTO users (user_name, email, user_role, created_at, updated_at)
VALUES 
    ('main_admin', 'admin@example.com', 'admin', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('elon_musk', 'elon2024@example.com', 'seller', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('ivan_ivanov', 'russian_bro@example.com', 'customer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('ronald_mcdonald', 'american_bro@example.com', 'seller', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    ('anna_petrova', 'russian_girl@example.com', 'customer', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
ON CONFLICT (user_name) DO NOTHING;

-- Insert into users_creds table (matching user_ids)
INSERT INTO users_creds (user_id, email, password_hash, created_at, updated_at)
VALUES 
    -- 3mRH8@tCO9_0
    (1, 'admin@example.com', '$2a$10$G8HPv2tPsVLp3Wntm.Yw5uU0sPlyuc7IoqJg6.6nYnSv.RJMYmdTm', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    -- 123mRH8@tCO9_0
    (2, 'elon2024@example.com', '$2a$10$Fv7RlNslndxJYGDgYeTJ/ubyWgqxPvovZhGGl8yi24eufWm8WkxkO', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    -- 12345mRH8@tCO9_0
    (3, 'russian_bro@example.com', '$2a$10$dY8hPqExR9071yf51DwpkO92iFJuyJXRjCUcv.Cqxp3.d08esmPtO', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    -- 123mRH8@tCO9_0
    (4, 'american_bro@example.com', '$2a$10$Fv7RlNslndxJYGDgYeTJ/ubyWgqxPvovZhGGl8yi24eufWm8WkxkO', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    -- 123mRH8@tCO9_0
    (5, 'russian_girl@example.com', '$2a$10$Fv7RlNslndxJYGDgYeTJ/ubyWgqxPvovZhGGl8yi24eufWm8WkxkO', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);

INSERT INTO products (
    title, 
    seller_name,
    seller_id,
    product_image, 
    product_description, 
    price,
    amount,
    created_at, 
    updated_at
)
VALUES
    ('Smartphone Ultra Super X', 
    'elon_musk',
    2,
    'http://localhost/static/phone.jpg', 
    'Latest smartphone with 128GB storage and OLED display', 
    79900,
    50,
    CURRENT_TIMESTAMP, 
    CURRENT_TIMESTAMP),
    
    ('Wireless Headphones Grusha', 
    'elon_musk',
    2,
    'http://localhost/static/headphones.jpg', 
    'Noise cancelling with 30-hour battery life', 
    19900,
    100,
    CURRENT_TIMESTAMP, 
    CURRENT_TIMESTAMP),
    
    ('Smart Watch Pro', 
    'ronald_mcdonald',
    4,
    'http://localhost/static/watch.jpg', 
    'Fitness tracking with heart rate and SpO2 monitoring', 
    24900,
    30,
    CURRENT_TIMESTAMP, 
    CURRENT_TIMESTAMP),
    
    ('Bluetooth Speaker JRIBL', 
    'ronald_mcdonald',
    4,
    'http://localhost/static/speaker.jpg', 
    'Waterproof portable speaker with 20W output', 
    8900,
    75,
    CURRENT_TIMESTAMP, 
    CURRENT_TIMESTAMP),
    
    ('4K Ultra HD TV', 
    'elon_musk',
    2,
    'http://localhost/static/TV.jpg', 
    '55-inch 4K Smart TV with HDR10+', 
    59900,
    15,
    CURRENT_TIMESTAMP, 
    CURRENT_TIMESTAMP),
    
    ('Gaming Laptop XVidia', 
    'ronald_mcdonald',
    4,
    'http://localhost/static/laptop.jpg', 
    'RTX 3070, 16GB RAM, 1TB SSD', 
    149900,
    10,
    CURRENT_TIMESTAMP, 
    CURRENT_TIMESTAMP)
ON CONFLICT (id) DO NOTHING;