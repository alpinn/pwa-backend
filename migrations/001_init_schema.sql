CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(36) PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(100) NOT NULL,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS products (
    id VARCHAR(36) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    stock INTEGER NOT NULL DEFAULT 0,
    image_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL REFERENCES users(id),
    total_amount DECIMAL(10, 2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transaction_items (
    id VARCHAR(36) PRIMARY KEY,
    transaction_id VARCHAR(36) NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    product_id VARCHAR(36) NOT NULL REFERENCES products(id),
    product_name VARCHAR(100) NOT NULL,
    quantity INTEGER NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    subtotal DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_transactions_user_id ON transactions(user_id);
CREATE INDEX idx_transactions_status ON transactions(status);
CREATE INDEX idx_transaction_items_transaction_id ON transaction_items(transaction_id);
CREATE INDEX idx_products_name ON products(name);

INSERT INTO users (id, username, password, name, role, created_at, updated_at) VALUES
('user-001', 'admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'Admin User', 'admin', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

INSERT INTO products (id, name, description, price, stock, image_url, created_at, updated_at) VALUES
('prod-001', 'Nasi Goreng', 'Nasi goreng spesial dengan telur', 25000, 100, 'https://via.placeholder.com/150', NOW(), NOW()),
('prod-002', 'Mie Goreng', 'Mie goreng pedas', 20000, 80, 'https://via.placeholder.com/150', NOW(), NOW()),
('prod-003', 'Es Teh Manis', 'Es teh manis segar', 5000, 200, 'https://via.placeholder.com/150', NOW(), NOW()),
('prod-004', 'Es Jeruk', 'Es jeruk peras segar', 7000, 150, 'https://via.placeholder.com/150', NOW(), NOW()),
('prod-005', 'Ayam Bakar', 'Ayam bakar bumbu kecap', 35000, 50, 'https://via.placeholder.com/150', NOW(), NOW()),
('prod-006', 'Sate Ayam', 'Sate ayam 10 tusuk', 30000, 60, 'https://via.placeholder.com/150', NOW(), NOW()),
('prod-007', 'Gado-gado', 'Gado-gado sayur lengkap', 15000, 40, 'https://via.placeholder.com/150', NOW(), NOW()),
('prod-008', 'Soto Ayam', 'Soto ayam kuah kuning', 22000, 70, 'https://via.placeholder.com/150', NOW(), NOW())
ON CONFLICT (id) DO NOTHING;