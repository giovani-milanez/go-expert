CREATE TABLE orders (
  id varchar(36) not null PRIMARY KEY, 
  price decimal(10,2) not null, 
  tax decimal(10,2) not null, 
  final_price decimal(10,2) not null
);