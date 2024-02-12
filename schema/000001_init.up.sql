create table users
(
    id  bigserial not null unique primary key ,
    login varchar unique not null ,
    password varchar not null,
    created_at timestamp not null default now()
);

create or replace function update_version()
    returns trigger as $$
begin
    new.update_at = now();
    new.version = old.version + 1;
    return new;
end;
$$ language plpgsql;

create table cards
(
    id  bigserial not null unique primary key ,
    user_id int references users (id) not null ,
    number  bytea not null unique,
    expired_at bytea not null,
    cvv bytea not null,
    version int default 1,
    created_at timestamp not null default now(),
    update_at timestamp not null default now()
);

create trigger update_cards_trigger
    before update on cards
    for each row
execute function update_version();

create table user_credentials
(
    id  bigserial not null unique primary key ,
    user_id int references users (id) not null ,
    login bytea unique not null ,
    password bytea not null,
    version int default 1,
    created_at timestamp not null default now(),
    update_at timestamp not null default now()
);

create trigger update_user_credentials_trigger
    before update on user_credentials
    for each row
execute function update_version();

create table text_data
(
    id  bigserial not null unique primary key ,
    user_id int references users (id) not null ,
    large_text bytea not null ,
    version int default 1,
    created_at timestamp not null default now(),
    update_at timestamp not null default now()
);

create trigger update_text_data_trigger
    before update on text_data
    for each row
execute function update_version();

create table binary_file
(
    id  bigserial not null unique primary key ,
    user_id int references users (id) not null ,
    path varchar not null ,
    created_at timestamp not null default now(),
    file_name varchar not null,
    size int not null
);