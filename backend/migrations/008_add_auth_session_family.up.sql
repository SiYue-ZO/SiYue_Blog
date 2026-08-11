ALTER TABLE auth_sessions
    ADD COLUMN family_id VARCHAR(64);

-- Existing refresh sessions predate token families. Walk each replacement
-- chain from its root so reuse of an old token also revokes its active child.
WITH RECURSIVE session_families AS (
    SELECT session.token_hash, session.token_hash AS family_id
    FROM auth_sessions AS session
    WHERE NOT EXISTS (
        SELECT 1
        FROM auth_sessions AS parent
        WHERE parent.replaced_by = session.token_hash
    )

    UNION ALL

    SELECT child.token_hash, family.family_id
    FROM session_families AS family
    JOIN auth_sessions AS parent ON parent.token_hash = family.token_hash
    JOIN auth_sessions AS child ON child.token_hash = parent.replaced_by
)
UPDATE auth_sessions AS session
SET family_id = family.family_id
FROM session_families AS family
WHERE session.token_hash = family.token_hash;

-- Defensive fallback for malformed legacy cycles or orphaned records.
UPDATE auth_sessions
SET family_id = token_hash
WHERE family_id IS NULL;

ALTER TABLE auth_sessions
    ALTER COLUMN family_id SET NOT NULL;

CREATE INDEX idx_auth_sessions_family ON auth_sessions (family_id);
