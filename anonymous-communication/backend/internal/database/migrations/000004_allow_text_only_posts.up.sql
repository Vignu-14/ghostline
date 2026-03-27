ALTER TABLE posts
    ALTER COLUMN image_url DROP NOT NULL;

ALTER TABLE posts
    ADD CONSTRAINT posts_content_required
    CHECK (
        image_url IS NOT NULL
        OR NULLIF(BTRIM(COALESCE(caption, '')), '') IS NOT NULL
    );
