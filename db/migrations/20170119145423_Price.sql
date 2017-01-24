-- +goose Up
CREATE TABLE Price (
  product_id   INT       NOT NULL,
  product_name TEXT      NOT NULL,
  cost         FLOAT     NOT NULL,
  created_at   TIMESTAMP NOT NULL DEFAULT now(),
  updated_at   TIMESTAMP NOT NULL DEFAULT now()
);
-- SQL in section 'Up' is executed when this migration is applied


-- +goose Down
DROP TABLE Price;
-- SQL section 'Down' is executed when this migration is rolled back

