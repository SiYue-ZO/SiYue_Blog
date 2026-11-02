DROP TRIGGER trg_posts_updated_at ON posts;
DROP FUNCTION update_post_updated_at_column();
CREATE TRIGGER trg_posts_updated_at
    BEFORE UPDATE ON posts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
