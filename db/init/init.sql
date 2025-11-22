DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'pr_status') THEN
        CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS team (
  team_name TEXT PRIMARY KEY,
  chat_id BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS user (
  user_id UUID PRIMARY KEY,
  user_name TEXT NOT NULL,
  team_name TEXT REFERENCES team(team_name),
  is_active BOOLEAN NOT NULL
);

CREATE TABLE IF NOT EXISTS pull_request (
  pr_id UUID PRIMARY KEY,
  pr_name TEXT NOT NULL,
  author_id UUID REFERENCES user(user_id),
  status pr_status NOT NULL,
  assigned_reviewers UUID[],
  created_at TIMESTAMP NOT NULL
  merged_at TIMESTAMP
);
