-- Drop trigger for checking and locking user on failed attempts
DROP TRIGGER IF EXISTS check_login_attempts_trigger ON login_histories;

-- Drop function for checking and locking user on failed attempts
DROP FUNCTION IF EXISTS check_and_lock_user;

-- Drop login_histories table
DROP TABLE IF EXISTS login_histories;

-- Remove columns for user activity and lock status
ALTER TABLE app_users
DROP COLUMN is_active,
DROP COLUMN is_locked;