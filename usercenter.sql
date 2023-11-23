use usercenter;
drop table verify_code;
CREATE TABLE IF NOT EXISTS verify_code (
    id bigint(20) not null primary key,
    verify_id varchar(1024) not null,
    answer varchar(512) not null,
    valid  smallint(4) not null default 1,
    stime int(11) unsigned,
    etime int(11) unsigned,
    key(verify_id)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;
