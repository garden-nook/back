alter table crops
    add column description text,
    add column soil_type_id int not null references soil_types(id),
    alter column sun_needs set default 1,
    alter column sun_needs set not null,
    alter column sun_needs drop default;
