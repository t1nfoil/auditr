\c audit;

CREATE TABLE da_activities (
    id SERIAL PRIMARY KEY,
    da_event_id VARCHAR(255),
    activity_date_time TIMESTAMPTZ NOT NULL,
    activity_display_name VARCHAR(255),
    category VARCHAR(255),
    correlation_id VARCHAR(255),
    logged_by_service VARCHAR(255),
    operation_type VARCHAR(255),
    result INT,
    result_reason TEXT
);

CREATE TABLE da_activities_initiated_by (
    id SERIAL PRIMARY KEY,
    da_activity_id INTEGER REFERENCES da_activities(id),
    app_id VARCHAR(255),
    app_name VARCHAR(255),
    user_name VARCHAR(255),
    user_id VARCHAR(255),
    ip_address VARCHAR(255),
    additional_data TEXT
);

CREATE TABLE da_activities_additional_details (
    id SERIAL PRIMARY KEY,
    da_activity_id INTEGER REFERENCES da_activities(id),
    key TEXT,
    value TEXT
);

CREATE TABLE da_activities_target_resources (
    id SERIAL PRIMARY KEY,
    da_activity_id INT REFERENCES da_activities(id),
    resource_id VARCHAR(255),
    display_name VARCHAR(255),
    group_type VARCHAR(255),
    target_resource_type VARCHAR(255),
    type_escaped VARCHAR(255),
    user_principal_name VARCHAR(255),
    additional_data TEXT
);

CREATE TABLE da_activities_target_resources_modified_properties (
    id SERIAL PRIMARY KEY,
    da_activity_id INT REFERENCES da_activities(id),
    da_target_resources_id INT REFERENCES da_activities_target_resources(id),
    display_name VARCHAR(255),
    new_value TEXT,
    old_value TEXT,
    additional_data TEXT
);


CREATE INDEX idx_initiated_by_da_activity_id ON da_activities_initiated_by(da_activity_id);
CREATE INDEX idx_additional_details_da_activity_id ON da_activities_additional_details(da_activity_id);
CREATE INDEX idx_target_resources_da_activity_id ON da_activities_target_resources(da_activity_id);
CREATE INDEX idx_target_resources_modified_properties_da_activity_id ON da_activities_target_resources_modified_properties(da_activity_id);
CREATE INDEX idx_activities_correlation_id ON da_activities(correlation_id);
CREATE INDEX idx_targetResources_resource_id ON da_activities_target_resources(resource_id);
CREATE INDEX idx_activities_da_event_id ON da_activities(da_event_id);
CREATE INDEX idx_activities_activity_display_name ON da_activities(activity_display_name);
