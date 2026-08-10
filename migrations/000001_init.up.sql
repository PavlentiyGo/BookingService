CREATE TYPE roles AS ENUM('user','admin');
CREATE TYPE statuses AS ENUM('active','cancelled');


CREATE TABLE users(
    id UUID PRIMARY KEY,
    role roles NOT NULL,
    email VARCHAR CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    password_hash VARCHAR,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE rooms(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR UNIQUE NOT NULL,
    description VARCHAR,
    capacity INTEGER DEFAULT 0,
    createdAt TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE schedules(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID UNIQUE NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    days_of_week SMALLINT[] CHECK ( array_length(days_of_week,1) <= 7 ),
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    CHECK ( start_time < end_time )
);

CREATE TABLE slots(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE ,
    start_time TIMESTAMPTZ NOT NULL,
    end_time TIMESTAMPTZ NOT NULL,
    CHECK ( start_time < end_time ),
    UNIQUE(room_id,start_time,end_time)
);

CREATE TABLE bookings(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slot_id UUID NOT NULL REFERENCES slots(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status statuses DEFAULT 'active',
    conference_link VARCHAR,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE UNIQUE INDEX idx_bookings_one_active_per_slot
    ON bookings(slot_id)
    WHERE status = 'active'