CREATE TABLE IF NOT EXISTS jobruns (
  id INT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  created TIMESTAMP NOT NULL,
  job VARCHAR(255) NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_created ON jobruns (created);
CREATE INDEX IF NOT EXISTS idx_job ON jobruns (job);

CREATE TABLE IF NOT EXISTS alerts (
  id INT NOT NULL GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  jobrun_id INT NOT NULL,
  line INT NOT NULL,
  rule VARCHAR(1023) NOT NULL,
  description TEXT NOT NULL,
  CONSTRAINT fk_jobrun FOREIGN KEY(jobrun_id) REFERENCES jobruns(id)
);

CREATE INDEX IF NOT EXISTS idx_jobrun_id ON alerts (jobrun_id);
CREATE INDEX IF NOT EXISTS idx_order ON alerts (order);
