CREATE TRIGGER items_insert_trigger
    AFTER INSERT ON items
    FOR EACH ROW
EXECUTE FUNCTION log_item_change();

CREATE TRIGGER items_update_trigger
    AFTER UPDATE ON items
    FOR EACH ROW
EXECUTE FUNCTION log_item_change();

CREATE TRIGGER items_delete_trigger
    BEFORE DELETE ON items
    FOR EACH ROW
EXECUTE FUNCTION log_item_change();