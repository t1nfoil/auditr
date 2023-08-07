DROP DATABASE audit;
DROP TABLE IF EXISTS da_activities CASCADE;
DROP TABLE IF EXISTS da_activities_initiated_by CASCADE;
DROP TABLE IF EXISTS da_activities_additional_details CASCADE;
DROP TABLE IF EXISTS da_target_resources CASCADE;
DROP TABLE IF EXISTS da_target_resources_modified_properties CASCADE;
DROP INDEX IF EXISTS idx_activities_correlationid CASCADE;
DROP INDEX IF EXISTS idx_targetresources_resourceid CASCADE;


CREATE DATABASE audit;
\c audit

