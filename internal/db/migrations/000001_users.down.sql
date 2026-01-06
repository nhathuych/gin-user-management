-- Drop trigger
DROP TRIGGER IF EXISTS users_set_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_user_updated_at();

-- Drop table
DROP TABLE IF EXISTS users;
