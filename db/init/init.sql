CREATE TABLE IF NOT EXISTS teams (
  team_name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users (
  user_id TEXT PRIMARY KEY,
  user_name TEXT NOT NULL,
  team_name TEXT REFERENCES teams(team_name),
  is_active BOOLEAN NOT NULL
);

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'pr_status') THEN
        CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS pull_requests (
  pr_id TEXT PRIMARY KEY,
  pr_name TEXT NOT NULL,
  author_id TEXT REFERENCES users(user_id),
  pr_status pr_status NOT NULL,
  created_at TIMESTAMP NOT NULL,
  merged_at TIMESTAMP
);

CREATE TABLE pr_reviewers (
  pr_id TEXT REFERENCES pull_requests(pr_id) ON DELETE CASCADE,
  user_id TEXT REFERENCES users(user_id),
  PRIMARY KEY (pr_id, user_id)
);
