create table if not exists missions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
        constraint missions_id_unique UNIQUE,
    assigned_cat_id UUId unique,
    completion_status varchar DEFAULT 'incomplete',

    FOREIGN KEY (assigned_cat_id) REFERENCES cats(id)
);

-- add 1-1 to cats 
ALTER TABLE cats
ADD COLUMN mission_id UUID DEFAULT NULL,
ADD CONSTRAINT fk_cats_mission FOREIGN KEY (mission_id) REFERENCES missions(id),
ADD CONSTRAINT unique_mission_assignment UNIQUE (mission_id);

-- trigger for preventing deletion if mission has a cat assigned
create or REPLACE FUNCTION prevent_mission_with_assigned_cat_deletion() 
RETURNS TRIGGER as $$
BEGIN 
    if OLD.assigned_cat_id is not null then
        raise EXCEPTION 'cannot delete a mission with assigned cat';
    end if;
    RETURN old;
end;
$$ LANGUAGE PLPGSQL;

create trigger trg_prevent_mission_with_assigned_cat_deletion
before DELETE on missions 
for EACH row EXECUTE FUNCTION prevent_mission_with_assigned_cat_deletion();

-- trigger to prevent cat assignment to completed missions
-- wasn't in task but sounds reasonable (needs to be communicated)
create or REPLACE FUNCTION prevent_cat_assignment_to_completed_missions()
RETURNS TRIGGER as $$
BEGIN
    if OLD.assigned_cat_id != NEW.assigned_cat_id AND old.completion_status = 'complete' then
        raise EXCEPTION 'can not reasign cat to completed task';
    end if;
    return OLD;
END;
$$ LANGUAGE PLPGSQL;

create trigger trg_prevent_cat_assignment_to_completed_missions
before DELETE on missions 
for EACH row EXECUTE FUNCTION prevent_cat_assignment_to_completed_missions();

-- Targets --
create table if not exists targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid()
        constraint targets_id_unique UNIQUE,
    name varchar default '',
    country varchar default '',
    completion_status varchar default 'incomplete',
    notes text default '',
    mission_id uuid not null,

    FOREIGN KEY (mission_id) REFERENCES missions(id) 
        on DELETE CASCADE
);

create unique index if not exists targets_mission_id_index
    on targets (mission_id);

--trigger to not allow deletion if targets count <= 1 
CREATE OR REPLACE FUNCTION check_min_targets() RETURNS TRIGGER AS $$
DECLARE
    tcount INTEGER;
BEGIN
    SELECT COUNT(*) INTO tcount
    FROM targets 
    WHERE mission_id = OLD.mission_id;
    
    IF tcount <= 1 THEN
        RAISE EXCEPTION 'Cannot delete all mission targets. The minimal allowed number of targets is 1';
    END IF;
    
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

create trigger trg_check_min_targets
before DELETE on targets
for each row EXECUTE FUNCTION check_min_targets();

-- trigger to not allow insertion  if targets count >=3
create or replace function check_max_targets()
returns trigger as $$
DECLARE
    tcount INTEGER;
BEGIN
    SELECT count(*) as tcount
    FROM targets 
    WHERE mission_id = NEW.mission_id;
    if (tcount >= 3) THEN
        raise EXCEPTION 'can not add new mission targets. maximum allowed number of targets is 3';
    end if;
    RETURN NEW;
END;
$$ LANGUAGE PLPGSQL;


create trigger trg_check_max_targets
before INSERT on targets
for each row EXECUTE FUNCTION check_max_targets();

-- trigger to prevent deletion of accomplished target
create or replace function check_target_complition()
returns trigger as $$
BEGIN
    if OLD.completion_status = 'complete' then
        raise EXCEPTION 'can not delete accomplished target';
    end if;
    return old;
END;
$$ LANGUAGE PLPGSQL;

create trigger trg_check_target_complition
before DELETE on targets
for each row EXECUTE FUNCTION check_target_complition();

-- trigger to not allow modifying notes on completed targets & missions 
create or replace function prevent_modifying_notes_on_completed_targets_or_missions() 
RETURNS trigger as $$
DECLARE 
    mission_status varchar;
BEGIN
    if old.notes != new.notes then
        if OLD.completion_status = 'complete' then 
            raise exception 'can not modify notes on completed target';
        end if;
        select completion_status into mission_status from missions where id = OLD.mission_id;
        if mission_status = 'complete' then 
            raise exception 'can not modify target''s notes on completed mission';
        end if;
    end if;
    RETURN new;
END;
$$ LANGUAGE PLPGSQL;


create trigger trg_prevent_modifying_notes_on_completed_targets_or_missions
before update on targets
for each row EXECUTE FUNCTION prevent_modifying_notes_on_completed_targets_or_missions();
