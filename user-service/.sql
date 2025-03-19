Table users {
  id integer [primary key]
  name varchar(50) 
  email varchar(100) [unique]
  password varchar(255)
  balance decimal(10,2) [default: 0]
  is_verified boolean [default: false]
  created_at timestamp
  updated_at timestamp
  deleted_at timestamp
}