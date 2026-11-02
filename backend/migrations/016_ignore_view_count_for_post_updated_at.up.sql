CREATE OR REPLACE FUNCTION update_post_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    IF (NEW.title, NEW.slug, NEW.excerpt, NEW.content, NEW.cover_image, NEW.status, NEW.is_pinned, NEW.author_id, NEW.published_at)
       IS DISTINCT FROM
       (OLD.title, OLD.slug, OLD.excerpt, OLD.content, OLD.cover_image, OLD.status, OLD.is_pinned, OLD.author_id, OLD.published_at) THEN
        NEW.updated_at = NOW();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER trg_posts_updated_at ON posts;
CREATE TRIGGER trg_posts_updated_at
    BEFORE UPDATE ON posts
    FOR EACH ROW EXECUTE FUNCTION update_post_updated_at_column();
