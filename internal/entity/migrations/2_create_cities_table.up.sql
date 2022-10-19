CREATE TABLE cities (
   id SERIAL PRIMARY KEY,
   user_id INT,
   token    VARCHAR(255),
   name   VARCHAR(255),
   rating INT,
   CONSTRAINT fk_user FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);