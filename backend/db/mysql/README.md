## Common Usages

### Auto increment

```sql
SHOW VARIABLES LIKE 'auto_inc%';
SET @@auto_increment_increment=2;
```

### delete

```sql
DELETE FROM table WHERE usn=0;
```

## DDLs

```sql
create table user
(
	usn bigint auto_increment comment 'user serial number' primary key,
	uid bigint default '100000' not null comment 'user id',
	avatar varchar(255) null comment 'user avatar url',
	birthday bigint default '0' not null comment 'birthday EPOCH timestamp',
	country varchar(255) null,
	email varchar(255) not null,
	gender int default '0' null comment '0: unknown, 1: female, 2: male',
	last_login bigint default '0' not null,
	login_count bigint default '0' not null,
	nickname varchar(255) not null,
	salt varchar(255) null comment 'salt',
	secret varchar(255) null comment 'hash(salt + password)',
	since bigint default '0' not null,
	constraint user_usn_uindex
		unique (usn),
	constraint user_uid_uindex
		unique (uid)
)
;

```

## TODO

add columns to table user

* last_ip
* app_version
* os
* timezone
* app_language
* os_locale
* device_type
* mcc
* phone_number
