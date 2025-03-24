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

Table email_verifications {
  id integer [primary key]
  email varchar(100) [ref: > users.email]
  token varchar(255)
  expired_at timestamp
  used boolean [default: false]
  created_at timestamp
}

Table password_resets {
  id integer [primary key]
  email varchar(100) [ref: > users.email]
  token varchar(255)
  expired_at timestamp
  used boolean [default: false]
  created_at timestamp
}