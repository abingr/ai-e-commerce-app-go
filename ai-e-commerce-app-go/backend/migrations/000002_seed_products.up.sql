INSERT INTO products (name, description, brand, category, price_cents, stock_quantity, image_url, is_active)
VALUES
    (
        'AeroBook Pro 14',
        'Lightweight laptop with a high-refresh display, 16GB memory, and fast SSD storage for productivity.',
        'AeroTech',
        'Laptops',
        129900,
        12,
        'https://images.unsplash.com/photo-1496181133206-80ce9b88a853?auto=format&fit=crop&w=900&q=80',
        true
    ),
    (
        'NovaPhone X',
        '5G smartphone with OLED display, dual camera system, and all-day battery life.',
        'NovaMobile',
        'Phones',
        89900,
        30,
        'https://images.unsplash.com/photo-1511707171634-5f897ff02aa9?auto=format&fit=crop&w=900&q=80',
        true
    ),
    (
        'KeyForge Mechanical Keyboard',
        'Hot-swappable mechanical keyboard with tactile switches and programmable lighting.',
        'KeyForge',
        'Accessories',
        12900,
        40,
        'https://images.unsplash.com/photo-1587829741301-dc798b83add3?auto=format&fit=crop&w=900&q=80',
        true
    ),
    (
        'ViewMax 27 Monitor',
        '27-inch QHD monitor with accurate colors, slim bezels, and USB-C connectivity.',
        'ViewMax',
        'Monitors',
        34900,
        18,
        'https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?auto=format&fit=crop&w=900&q=80',
        true
    ),
    (
        'PulseBeat Wireless Headphones',
        'Noise-cancelling wireless headphones with comfortable ear cups and long battery life.',
        'PulseBeat',
        'Audio',
        19900,
        25,
        'https://images.unsplash.com/photo-1505740420928-5e560c06d30e?auto=format&fit=crop&w=900&q=80',
        true
    ),
    (
        'FitTime Smart Watch',
        'Fitness-focused smart watch with heart-rate tracking, GPS, and water resistance.',
        'FitTime',
        'Wearables',
        24900,
        22,
        'https://images.unsplash.com/photo-1523275335684-37898b6baf30?auto=format&fit=crop&w=900&q=80',
        true
    )
ON CONFLICT DO NOTHING;
