ALTER TABLE users ADD COLUMN variant VARCHAR(1) NOT NULL DEFAULT 'Z';

-- In doing A/B testing, we want to assign a variant to each user. The variant is a single character, either 'A' or 'B'. We'll default to 'Z' for now.