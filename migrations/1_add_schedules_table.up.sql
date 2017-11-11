create table schedules (
    id integer primary key autoincrement,
    path string not null unique,
    size integer,
    mod_time datetime,
    time_to_destroy datetime
);

create index file_path_index on schedules(path);
