create table skipscantest (
	account_id text,
	sequence_id bigint,
	data int
);

create index on skipscantest (
	account_id, sequence_id desc
);
