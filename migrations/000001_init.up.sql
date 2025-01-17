CREATE TABLE clients (
    id UUID PRIMARY KEY,
    email VARCHAR(128) NOT NULL
);

CREATE TABLE orders (
    id UUID PRIMARY KEY,
    client_id UUID REFERENCES clients(id),
    created_at TIMESTAMP NOT NULL,
    confirmed_at TIMESTAMP,
    sended_at TIMESTAMP,
    delivered_at TIMESTAMP
);