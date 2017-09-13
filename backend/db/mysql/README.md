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
	usn bigint auto_increment primary key,
	uid bigint default '100000' not null,
	app_language varchar(16) null,
	app_version varchar(32) null,
	avatar varchar(255) null,
	birthday bigint default '0' not null,
	country varchar(16) null,
	device_type int default '0' null,
	email varchar(255) not null,
	email_verified tinyint(1) default '0' null,
	gender int default '0' null,
	is_official tinyint(1) default '0' null,
	last_ip varchar(128) null,
	last_login bigint default '0' not null,
	login_count bigint default '0' not null,
	mcc int default '0' null,
	nickname varchar(128) not null,
	os varchar(255) null,
	os_locale varchar(16) null,
	premium_level int default '0' null,
	salt varchar(255) null,
	secret varchar(255) null,
	since bigint default '0' not null,
	social_id varchar(32) null,
	status int default '0' null,
	timezone int default '0' null,
	token varchar(255) null,
	constraint user_usn_uindex
		unique (usn),
	constraint user_uid_uindex
		unique (uid)
)
;


```
