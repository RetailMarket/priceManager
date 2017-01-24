
-- +goose Up
CREATE TABLE PriceUpdateRequest (
  product_id 				  int not null,
  product_name			  text not null,
  cost				        float not null,
  status			        text not null,
  created_at   TIMESTAMP NOT NULL DEFAULT now(),
  updated_at   TIMESTAMP NOT NULL DEFAULT now()
);
-- SQL in section 'Up' is executed when this migration is applied


-- +goose Down
DROP TABLE PriceUpdateRequest;
-- SQL section 'Down' is executed when this migration is rolled back

