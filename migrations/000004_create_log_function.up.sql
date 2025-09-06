CREATE OR REPLACE FUNCTION log_item_change() RETURNS TRIGGER AS
$$
BEGIN
    IF (TG_OP = 'INSERT') THEN
        INSERT INTO item_history (item_id, action, changed_by, new_values)
        VALUES (NEW.id, 'create', current_setting('app.user', true), to_jsonb(NEW));
        RETURN NEW;
    ELSIF (TG_OP = 'UPDATE') THEN
        INSERT INTO item_history (item_id, action, changed_by, old_values, new_values)
        VALUES (NEW.id, 'update', current_setting('app.user', true), to_jsonb(OLD), to_jsonb(NEW));
        RETURN NEW;
    ELSIF (TG_OP = 'DELETE') THEN
        INSERT INTO item_history (item_id, action, changed_by, old_values)
        VALUES (OLD.id, 'delete', current_setting('app.user', true), to_jsonb(OLD));
        RETURN OLD;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;
