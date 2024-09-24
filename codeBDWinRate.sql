create database WinRate;

use WinRate;

create table users (
	userID int auto_increment primary key,
    email varchar(255) not null unique, 
    password varchar(255) not null
);

create table decks (
	deckID int auto_increment primary key, 
    user_id int,
    deck_name varchar(255) not null,
    foreign key (user_id) references user(userID)
);

create table matches (
	id int auto_increment primary key,
    user_deck_id int,
    opponent_deck_id int,
    victories int default 0,
    defeats  int default 0,
    created_at timestamp default current_timestamp,
    foreign key (user_deck_id) references decks(deckID),
    foreign key (opponent_deck_id) references decks(deckID)
);

use WinRate;
insert into user (email, password)
values ('lukesky@gmail.com', 'force');


show tables;
