#!/bin/bash

docker run --name postgres-test -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d -p 5432:5432 postgres:latest

echo "Postgresql container is starting..."
sleep 3

docker exec -it postgres-test psql -U postgres -d postgres -c "CREATE DATABASE messages;"
sleep 2

echo "Database messages created successfully"

docker exec -it postgres-test psql -U postgres -d messages -c "
CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY,
    content VARCHAR(255) NOT NULL,
    recipient_phone_number VARCHAR(20) NOT NULL,
    sent BOOLEAN NOT NULL
);
"

sleep 3
echo "Table messages created successfully"

docker exec -it postgres-test psql -U postgres -d messages -c "
INSERT INTO messages (id, content, recipient_phone_number, sent) VALUES ('24228df7-d5e9-44f7-a44a-37b10a473e3f', 'Merhaba, nasılsın?', '+905551111111', true);
INSERT INTO messages (id, content, recipient_phone_number, sent) VALUES ('6596d0b0-0a9c-45b0-a48b-1193ca0d0a98', 'Toplantı saat 14:00''da.', '+905551111111', false);
INSERT INTO messages (id, content, recipient_phone_number, sent) VALUES ('7e4c157b-a819-45b2-a2e7-0f96e5b132f1', 'Yeni bir ürünümüz çıktı, incelediniz mi?', '+905551111111', true);
INSERT INTO messages (id, content, recipient_phone_number, sent) VALUES ('18389670-5502-43c1-9010-26a42482e8a7', 'Doğum günün kutlu olsun!', '+905551111111', true);
INSERT INTO messages (id, content, recipient_phone_number, sent) VALUES ('6db17e4b-b043-4769-a3f3-b120a7b99174', 'Yarın hava nasıl olacak?', '+905551111111', false);
"

echo "Data inserted successfully"
