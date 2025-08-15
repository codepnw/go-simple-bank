CREATE TYPE users_role AS ENUM ('USER', 'STAFF', 'ADMIN');

ALTER TABLE users ADD COLUMN role users_role;

UPDATE users SET role = 'USER' WHERE role IS NULL;

ALTER TABLE users ALTER COLUMN role SET NOT NULL;