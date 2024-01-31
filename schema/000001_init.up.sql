create table users
(
    id  bigserial not null unique primary key ,
    login varchar unique not null ,
    password varchar not null,
    created_at timestamp not null default now()
);

create table cards
(
    id  bigserial not null unique primary key ,
    user_id int references users (id) not null ,
    number varchar not null,
    expired_at varchar not null,
    cvv varchar not null,
    version int default 1,
    created_at timestamp not null default now(),
    update_at timestamp not null default now()
);

create or replace function update_cards()
    returns trigger as $$
begin
    new.update_at = now();
    new.version = old.version + 1;
    return new;
end;
$$ language plpgsql;

create trigger update_cards_trigger
    before update on cards
    for each row
execute function update_cards();