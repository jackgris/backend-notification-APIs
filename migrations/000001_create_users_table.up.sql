CREATE TABLE users (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100),
  email VARCHAR(100),
  phone VARCHAR(15),
  subscribed_categories TEXT[], -- e.g., '{Sports,Finance}'
  notification_channels TEXT[] -- e.g., '{SMS,Email,PushNotification}'
);
