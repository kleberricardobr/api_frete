CREATE TABLE IF NOT EXISTS freight_carriers (
    id SERIAL PRIMARY KEY,    
    carrier_name VARCHAR(255) NOT NULL,
    service VARCHAR(100) NOT NULL,
    deadline INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP    
);