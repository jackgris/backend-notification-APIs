-- Insert 5 mock users into the users table

INSERT INTO users (name, email, phone, subscribed_categories, notification_channels)
VALUES
-- User 1: Subscribed to Sports and Finance, prefers SMS and Email notifications
('Alice Johnson', 'alice.johnson@example.com', '+12345678901',
  ARRAY['Sports', 'Finance'], ARRAY['SMS', 'Email']),

-- User 2: Subscribed to Films, prefers Push notifications
('Bob Smith', 'bob.smith@example.com', '+12345678902',
  ARRAY['Films'], ARRAY['PushNotification']),

-- User 3: Subscribed to Sports and Films, prefers SMS and Push notifications
('Charlie Lee', 'charlie.lee@example.com', '+12345678903',
  ARRAY['Sports', 'Films'], ARRAY['SMS', 'PushNotification']),

-- User 4: Subscribed to Finance, prefers Email notifications
('Dana White', 'dana.white@example.com', '+12345678904',
  ARRAY['Finance'], ARRAY['Email']),

-- User 5: Subscribed to all categories, prefers all notification types
('Eve Black', 'eve.black@example.com', '+12345678905',
  ARRAY['Sports', 'Finance', 'Films'], ARRAY['SMS', 'Email', 'PushNotification']);
