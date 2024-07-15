-- +goose Up
alter table posts 
alter column title type varchar(255);

alter table posts 
alter column description type varchar(500);

-- +goose Down
alter table posts 
alter column title type varchar(60);

alter table posts 
alter column description type varchar(300);

