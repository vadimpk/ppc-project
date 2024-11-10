-- +goose Up
-- +goose StatementBegin

-- Businesses
CREATE TABLE businesses
(
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    logo_url     TEXT,
    color_scheme JSONB,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Users and Authentication
CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    business_id   INTEGER REFERENCES businesses (id),
    email         VARCHAR(255),
    phone         VARCHAR(20),
    full_name     VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255),
    role          VARCHAR(20)  NOT NULL CHECK (role IN ('admin', 'employee', 'client')),
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (business_id, email),
    UNIQUE (business_id, phone)
);

-- Services
CREATE TABLE services
(
    id          SERIAL PRIMARY KEY,
    business_id INTEGER REFERENCES businesses (id),
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    duration    INTEGER      NOT NULL, -- in minutes
    price       INTEGER      NOT NULL, -- in cents
    is_active   BOOLEAN                  DEFAULT true,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Employees
CREATE TABLE employees
(
    id             SERIAL PRIMARY KEY,
    business_id    INTEGER REFERENCES businesses (id),
    user_id        INTEGER REFERENCES users (id),
    specialization TEXT,
    is_active      BOOLEAN                  DEFAULT true,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Employee Services (Many-to-Many relationship)
CREATE TABLE employee_services
(
    employee_id INTEGER REFERENCES employees (id),
    service_id  INTEGER REFERENCES services (id),
    PRIMARY KEY (employee_id, service_id)
);

-- Regular Schedule (weekly recurring)
CREATE TABLE schedule_templates
(
    id          SERIAL PRIMARY KEY,
    employee_id INTEGER REFERENCES employees (id),
    day_of_week INTEGER CHECK (day_of_week BETWEEN 0 AND 6),
    start_time  TIME NOT NULL,
    end_time    TIME NOT NULL,
    is_break    BOOLEAN                  DEFAULT false,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Schedule Overrides (specific dates)
CREATE TABLE schedule_overrides
(
    id             SERIAL PRIMARY KEY,
    employee_id    INTEGER REFERENCES employees (id),
    override_date  DATE NOT NULL,
    start_time     TIME,
    end_time       TIME,
    is_working_day BOOLEAN                  DEFAULT true,
    is_break       BOOLEAN                  DEFAULT false,
    created_at     TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- Appointments
CREATE TABLE appointments
(
    id            SERIAL PRIMARY KEY,
    business_id   INTEGER REFERENCES businesses (id),
    client_id     INTEGER REFERENCES users (id),
    employee_id   INTEGER REFERENCES employees (id),
    service_id    INTEGER REFERENCES services (id),
    start_time    TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time      TIMESTAMP WITH TIME ZONE NOT NULL,
    status        VARCHAR(20) CHECK (status IN ('scheduled', 'completed', 'cancelled', 'no_show')),
    reminder_time INTEGER, -- minutes before appointment
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indices for better query performance
CREATE INDEX idx_users_business ON users (business_id);
CREATE INDEX idx_services_business ON services (business_id);
CREATE INDEX idx_employees_business ON employees (business_id);
CREATE INDEX idx_appointments_business ON appointments (business_id);
CREATE INDEX idx_appointments_client_id ON appointments (client_id);
CREATE INDEX idx_appointments_employee_id ON appointments (employee_id);
CREATE INDEX idx_appointments_start_time ON appointments (start_time);
CREATE INDEX idx_schedule_templates_employee_id ON schedule_templates (employee_id);
CREATE INDEX idx_schedule_overrides_employee_date ON schedule_overrides (employee_id, override_date);
CREATE INDEX idx_employee_services_employee_id ON employee_services (employee_id);
CREATE INDEX idx_employee_services_service_id ON employee_services (service_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_employee_services_service_id;
DROP INDEX IF EXISTS idx_employee_services_employee_id;
DROP INDEX IF EXISTS idx_schedule_overrides_employee_date;
DROP INDEX IF EXISTS idx_schedule_templates_employee_id;
DROP INDEX IF EXISTS idx_appointments_start_time;
DROP INDEX IF EXISTS idx_appointments_employee_id;
DROP INDEX IF EXISTS idx_appointments_client_id;
DROP INDEX IF EXISTS idx_appointments_business;
DROP INDEX IF EXISTS idx_employees_business;
DROP INDEX IF EXISTS idx_services_business;
DROP INDEX IF EXISTS idx_users_business;

DROP TABLE IF EXISTS appointments;
DROP TABLE IF EXISTS schedule_override_slots;
DROP TABLE IF EXISTS schedule_overrides;
DROP TABLE IF EXISTS schedule_templates;
DROP TABLE IF EXISTS employee_services;
DROP TABLE IF EXISTS employees;
DROP TABLE IF EXISTS services;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS businesses;
-- +goose StatementEnd
