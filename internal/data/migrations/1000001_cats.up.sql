create table if not exists cats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
        constraint cats_id_unique UNIQUE,
    name varchar DEFAULT '', 
    years_of_experience integer DEFAULT 0,
    breed varchar not null,
    salary DECIMAL(12,2) DEFAULT 0.00
)