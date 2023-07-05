CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS cube;
CREATE EXTENSION IF NOT EXISTS earthdistance;

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_by UUID REFERENCES users(id)
);

CREATE TABLE user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    user_id UUID REFERENCES users(id) NOT NULL,
    role_user TEXT NOT NULL
);

CREATE TABLE address (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name TEXT NOT NULL,
    latitude DECIMAL NOT NULL,
    longitude DECIMAL NOT NULL,
    user_id UUID REFERENCES users(id)
);

CREATE TABLE restaurants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name TEXT NOT NULL,
    latitude DECIMAL NOT NULL,
    longitude DECIMAL NOT NULL,
    created_by UUID REFERENCES users(id)
);

CREATE TABLE dishes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
    name TEXT NOT NULL,
    price INT NOT NULL,
    restaurant_id UUID REFERENCES restaurants(id),
    created_by UUID REFERENCES users(id)
);

INSERT INTO users(name, email, password) VALUES ('Siddhant', 'abc@gmail.com', '$2a$12$iwTamkdHO27Rhc43SfV2.u0AYT096KDAZOiSyX9FrOaU1h5NSGR8i');

WITH siddhant_user AS (
    SELECT id
    FROM users
    WHERE name = 'Siddhant'
    LIMIT 1
    )
INSERT INTO user_roles (user_id, role_user)
SELECT id, 'admin'
FROM siddhant_user;