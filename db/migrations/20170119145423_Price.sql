-- +goose Up
CREATE TABLE Price (
  product_id   INT       NOT NULL,
  product_name TEXT      NOT NULL,
  cost         FLOAT     NOT NULL,
  version      TEXT      NOT NULL,
  is_latest    BOOLEAN   NOT NULL
);
-- SQL in section 'Up' is executed when this migration is applied


-- +goose Down
DROP TABLE Price;
-- SQL section 'Down' is executed when this migration is rolled back

