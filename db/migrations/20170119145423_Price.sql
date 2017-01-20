
-- +goose Up
CREATE TABLE Price (
  product_id 				  int not null,
  product_name	  		text not null,
  cost				        float not null
);
-- SQL in section 'Up' is executed when this migration is applied


-- +goose Down
DROP TABLE Price;
-- SQL section 'Down' is executed when this migration is rolled back

