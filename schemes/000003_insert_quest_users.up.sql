INSERT INTO users 
    (id, email, password, confirm_status, activation, name, surname, token_secret_key, confirmed_at, created_at)
    VALUES (1, 'butago_quest@quset.com', 'butago_quest', 'quest', true, 'Quest', 'Quest', '000000', NOW()::timestamp, NOW()::timestamp) ON CONFLICT (id) DO UPDATE SET name = 'Quest', surname = 'Quest';