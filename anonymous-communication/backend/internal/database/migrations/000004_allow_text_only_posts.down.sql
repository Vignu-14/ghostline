ALTER TABLE posts
    DROP CONSTRAINT IF EXISTS posts_content_required;

DELETE FROM posts
WHERE image_url IS NULL;

ALTER TABLE posts
    ALTER COLUMN image_url SET NOT NULL;
