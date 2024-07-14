-- +goose Up
alter table feed_user
drop constraint feed_user_feed_id_fkey,
drop constraint feed_user_user_id_fkey;

alter table feed_user
add constraint fk_feed_user_feed_id
    foreign key (feed_id) references feeds(id) on delete cascade,
add constraint fk_feed_user_user_id
    foreign key (user_id) references users(id) on delete cascade;

-- +goose Down
alter table feed_user
drop constraint fk_feed_user_feed_id,
drop constraint fk_feed_user_user_id;

alter table feed_user
add constraint fk_feed_user_feed_id
    foreign key (feed_id) references feeds(id),
add constraint fk_feed_user_user_id
    foreign key (user_id) references users(id);
