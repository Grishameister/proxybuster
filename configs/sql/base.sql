create table if not exists requests(
	id serial primary key,
	url text not null,
	req text not null
);
