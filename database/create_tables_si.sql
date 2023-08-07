\c audit;

CREATE TABLE si_activities (
    id SERIAL PRIMARY KEY,
    si_event_id VARCHAR(255),
    app_display_name VARCHAR(255),
    app_id VARCHAR(255),
    applied_conditional_access_policies TEXT,
    client_app_used VARCHAR(255),
    conditional_access_status VARCHAR(255),
    correlation_id VARCHAR(255),
    created_date_time TIMESTAMP,
    device_detail TEXT,
    ip_address VARCHAR(255),
    is_interactive BOOLEAN,
    location TEXT,
    resource_display_name VARCHAR(255),
    resource_id VARCHAR(255),
    risk_detail VARCHAR(255),
    risk_event_types TEXT,
    risk_event_types_v2 TEXT,
    risk_level_aggregated VARCHAR(255),
    risk_level_during_sign_in VARCHAR(255),
    risk_state VARCHAR(255),
    status VARCHAR(255),
    user_display_name VARCHAR(255),
    user_id VARCHAR(255),
    user_principal_name VARCHAR(255)
);

CREATE TABLE si_activities_device_details (
   id SERIAL PRIMARY KEY,
   si_activity_id INTEGER REFERENCES si_activities(id),
   additional_data TEXT,
   browser VARCHAR(255),
   device_id VARCHAR(255),
   display_name VARCHAR(255),
   is_compliant BOOLEAN,
   is_managed BOOLEAN,
   o_data_type VARCHAR(255),
   operating_system VARCHAR(255),
   trust_type VARCHAR(255)
);

CREATE TABLE si_activities_location (
   id SERIAL PRIMARY KEY,
   si_activity_id INTEGER REFERENCES si_activities(id),
   city VARCHAR(255),
   country_or_region VARCHAR(255),
   state VARCHAR(255)
);

CREATE TABLE si_activities_geocoordinates (
   id SERIAL PRIMARY KEY,
   si_activity_id INTEGER REFERENCES si_activities(id),
   si_location_id INTEGER REFERENCES si_activities_location(id),
   additional_data TEXT,
   latitude VARCHAR(64),
   longitude VARCHAR(64)
);

CREATE TABLE si_activities_status (
   id SERIAL PRIMARY KEY,
   si_activity_id INTEGER REFERENCES si_activities(id),
   additional_data TEXT,
   additional_detail TEXT,
   error_code VARCHAR(64),
   failure_reason TEXT
);

CREATE INDEX idx_siactivities_correlation_id ON si_activities(correlation_id);
CREATE INDEX idx_siactivities_app_id ON si_activities(app_id);
CREATE INDEX idx_device_details_si_activity_id ON si_activities_device_details(si_activity_id);
CREATE INDEX idx_location_si_activity_id ON si_activities_location(si_activity_id);
CREATE INDEX idx_geocoordinates_si_activity_id ON si_activities_geocoordinates(si_activity_id);
CREATE INDEX idx_status_si_activity_id ON si_activities_status(si_activity_id);
CREATE INDEX idx_siactivities_user_id ON si_activities(user_id);
CREATE INDEX idx_siactivities_created_date_time ON si_activities(created_date_time);
CREATE INDEX idx_siactivities_is_interactive ON si_activities(is_interactive);
