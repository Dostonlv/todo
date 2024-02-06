
create table todo (
  id serial primary key,
  title text not null,
  completed boolean not null default false
  created_at timestamp not null default now()
);


-- create a user table
create table users (
  id serial primary key,
  login varchar not null unique,
  password varchar default '',
  created_at timestamp with time zone default now()
);


insert into users (login, password) values ('admin', 'admin');
