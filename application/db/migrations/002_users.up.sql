-- Add columns for user activity and lock status
ALTER TABLE app_users
ADD COLUMN is_active BOOLEAN DEFAULT TRUE,
ADD COLUMN is_locked BOOLEAN DEFAULT FALSE;

CREATE TABLE login_histories (
    login_id INT GENERATED ALWAYS AS IDENTITY,
    user_id INT REFERENCES app_users(user_id),
    username TEXT,
    login_time TIMESTAMPTZ DEFAULT NOW(),
    login_status BOOLEAN
);

-- Create function to check and lock user if 3 consecutive failed login attempts within 2 hours
CREATE OR REPLACE FUNCTION check_and_lock_user() RETURNS TRIGGER AS $$
DECLARE
    login_count INT := 0;
    last_login_status BOOLEAN := NULL;
    login_cursor CURSOR FOR
        SELECT login_status
        FROM login_histories
        WHERE user_id = NEW.user_id
        AND login_time >= NOW() - INTERVAL '2 hours'
        ORDER BY login_time DESC
		LIMIT 3;

BEGIN
    OPEN login_cursor;
    LOOP
        FETCH login_cursor INTO last_login_status;
        IF last_login_status = FALSE THEN
            login_count := login_count + 1;
        ELSE
            login_count := 0; -- Reset login count if login_status is true
        END IF;

        IF login_count = 3 THEN
            UPDATE app_users
            SET is_locked = TRUE
            WHERE user_id = NEW.user_id;
            EXIT;
        END IF;

        EXIT WHEN NOT FOUND;
    END LOOP;
    CLOSE login_cursor;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger to call the check_and_lock_user function
CREATE TRIGGER check_login_attempts_trigger
AFTER INSERT ON login_histories
FOR EACH ROW
EXECUTE FUNCTION check_and_lock_user();

-- Create trigger to update is_locked to false if new row login_status is true
CREATE OR REPLACE FUNCTION update_is_locked() RETURNS TRIGGER AS $$
BEGIN
    IF NEW.login_status = TRUE THEN
        UPDATE app_users
        SET is_locked = FALSE
        WHERE user_id = NEW.user_id;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_is_locked_trigger
AFTER INSERT ON login_histories
FOR EACH ROW
EXECUTE FUNCTION update_is_locked();