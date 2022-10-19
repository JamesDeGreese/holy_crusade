CREATE TABLE balance (
   id SERIAL PRIMARY KEY,
   city_id    INT,
   gold       INT,
   population INT,
   workers    INT,
   solders    INT,
   heroes     INT,
   CONSTRAINT fk_city FOREIGN KEY(city_id) REFERENCES cities(id) ON DELETE CASCADE
);