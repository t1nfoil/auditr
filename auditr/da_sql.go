package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type DAActivity struct {
	ID                  int       `gorm:"primaryKey;column:id"`
	EventID             string    `gorm:"column:da_event_id"`
	ActivityDateTime    time.Time `gorm:"column:activity_date_time"`
	ActivityDisplayName string    `gorm:"column:activity_display_name"`
	Category            string    `gorm:"column:category"`
	CorrelationID       string    `gorm:"column:correlation_id"`
	LoggedByService     string    `gorm:"column:logged_by_service"`
	OperationType       string    `gorm:"column:operation_type"`
	Result              int       `gorm:"column:result"`
	ResultReason        string    `gorm:"column:result_reason"`
}

func (DAActivity) TableName() string {
	return "da_activities"
}

type DAActivityInitiatedBy struct {
	ID             int    `gorm:"primaryKey;column:id"`
	DAActivityID   int    `gorm:"column:da_activity_id;references:id"`
	AppID          string `gorm:"column:app_id"`
	AppName        string `gorm:"column:app_name"`
	UserName       string `gorm:"column:user_name"`
	UserID         string `gorm:"column:user_id"`
	IPAddress      string `gorm:"column:ip_address"`
	AdditionalData string `gorm:"column:additional_data"`
}

func (DAActivityInitiatedBy) TableName() string {
	return "da_activities_initiated_by"
}

type DAActivityAdditionalDetail struct {
	ID           int    `gorm:"primaryKey;column:id"`
	DAActivityID int    `gorm:"column:da_activity_id;references:id"`
	Key          string `gorm:"column:key"`
	Value        string `gorm:"column:value"`
}

func (DAActivityAdditionalDetail) TableName() string {
	return "da_activities_additional_details"
}

type DAActivityTargetResource struct {
	ID                 int    `gorm:"primaryKey;column:id"`
	DAActivityID       int    `gorm:"column:da_activity_id;references:id"`
	ResourceID         string `gorm:"column:resource_id"`
	DisplayName        string `gorm:"column:display_name"`
	GroupType          string `gorm:"column:group_type"`
	TargetResourceType string `gorm:"column:target_resource_type"`
	TypeEscaped        string `gorm:"column:type_escaped"`
	UserPrincipalName  string `gorm:"column:user_principal_name"`
	AdditionalData     string `gorm:"column:additional_data"`
}

func (DAActivityTargetResource) TableName() string {
	return "da_activities_target_resources"
}

type DAActivityTargetResourceModifiedProperty struct {
	ID                  int    `gorm:"primaryKey;column:id"`
	DAActivityID        int    `gorm:"column:da_activity_id;references:id"`
	DATargetResourcesID int    `gorm:"column:da_target_resources_id;references:da_target_resources_id"`
	DisplayName         string `gorm:"column:display_name"`
	NewValue            string `gorm:"column:new_value"`
	OldValue            string `gorm:"column:old_value"`
	AdditionalData      string `gorm:"column:additional_data"`
}

func (DAActivityTargetResourceModifiedProperty) TableName() string {
	return "da_activities_target_resources_modified_properties"
}

func processDAActivity(modelValue DirectoryAuditWrite) (DAActivity, error) {
	activity := DAActivity{}
	if modelValue.DirectoryAudit.GetId() != nil {
		activity.EventID = *modelValue.DirectoryAudit.GetId()
	}
	if modelValue.DirectoryAudit.GetActivityDateTime() != nil {
		activity.ActivityDateTime = *modelValue.DirectoryAudit.GetActivityDateTime()
	}
	if modelValue.DirectoryAudit.GetActivityDisplayName() != nil {
		activity.ActivityDisplayName = *modelValue.DirectoryAudit.GetActivityDisplayName()
	}
	if modelValue.DirectoryAudit.GetCategory() != nil {
		activity.Category = *modelValue.DirectoryAudit.GetCategory()
	}
	if modelValue.DirectoryAudit.GetCorrelationId() != nil {
		activity.CorrelationID = *modelValue.DirectoryAudit.GetCorrelationId()
	}
	if modelValue.DirectoryAudit.GetResult() != nil {
		activity.Result = int(*modelValue.DirectoryAudit.GetResult())
	}
	if modelValue.DirectoryAudit.GetResultReason() != nil {
		activity.ResultReason = *modelValue.DirectoryAudit.GetResultReason()
	}
	if modelValue.DirectoryAudit.GetOperationType() != nil {
		activity.OperationType = *modelValue.DirectoryAudit.GetOperationType()
	}
	if modelValue.DirectoryAudit.GetLoggedByService() != nil {
		activity.LoggedByService = *modelValue.DirectoryAudit.GetLoggedByService()
	}
	return activity, nil
}

func processDAActivityAdditionalDetails(modelValue DirectoryAuditWrite) ([]DAActivityAdditionalDetail, error) {
	activityAdditionalDetails := make([]DAActivityAdditionalDetail, 0)
	if additionalDetails := modelValue.DirectoryAudit.GetAdditionalDetails(); additionalDetails != nil {
		for _, v := range additionalDetails {
			if v.GetKey() != nil && v.GetValue() != nil {
				cleanedString := strings.ReplaceAll(*v.GetValue(), "\r\n", "")
				activityAdditionalDetails = append(activityAdditionalDetails, DAActivityAdditionalDetail{
					Key:   *v.GetKey(),
					Value: cleanedString,
				})
			}
		}
	}
	return activityAdditionalDetails, nil
}

func processDAActivityInitiatedBy(modelValue DirectoryAuditWrite) (DAActivityInitiatedBy, error) {
	var activityInitiatedBy DAActivityInitiatedBy
	initiatedBy := modelValue.DirectoryAudit.GetInitiatedBy()
	if initiatedBy == nil {
		return activityInitiatedBy, nil
	}
	user := initiatedBy.GetUser()
	if user != nil {
		if displayName := user.GetDisplayName(); displayName != nil {
			activityInitiatedBy.UserName = *displayName
		}
		if userID := user.GetId(); userID != nil {
			activityInitiatedBy.UserID = *userID
		}
		if ipAddress := user.GetIpAddress(); ipAddress != nil {
			activityInitiatedBy.IPAddress = *ipAddress
		}
	}
	app := initiatedBy.GetApp()
	if app != nil {
		if appID := app.GetAppId(); appID != nil {
			activityInitiatedBy.AppID = *appID
		}
		if appName := app.GetDisplayName(); appName != nil {
			activityInitiatedBy.AppName = *appName
		}
	}

	return activityInitiatedBy, nil
}

func processDAActivityTargetResources(modelValue DirectoryAuditWrite) ([]DAActivityTargetResource, error) {
	activityTargetResources := make([]DAActivityTargetResource, 0)
	if modelValue.DirectoryAudit.GetAdditionalData() != nil {
		for _, resource := range modelValue.DirectoryAudit.GetTargetResources() {
			var targetResource DAActivityTargetResource
			if resource.GetDisplayName() != nil {
				targetResource.DisplayName = *resource.GetDisplayName()
			}
			if resource.GetGroupType() != nil {
				targetResource.GroupType = resource.GetGroupType().String()
			}
			if resource.GetId() != nil {
				targetResource.ResourceID = *resource.GetId()
			}
			if resource.GetTypeEscaped() != nil {
				targetResource.TypeEscaped = *resource.GetTypeEscaped()
			}
			if resource.GetUserPrincipalName() != nil {
				targetResource.UserPrincipalName = *resource.GetUserPrincipalName()
			}
			if resource.GetGroupType() != nil {
				targetResource.GroupType = resource.GetGroupType().String()
			}
			if resource.GetOdataType() != nil {
				targetResource.TargetResourceType = *resource.GetOdataType()
			}
			activityTargetResources = append(activityTargetResources, targetResource)
		}
	}
	return activityTargetResources, nil
}

func processDAActivityTargetResourcesModifiedProperties(modelValue DirectoryAuditWrite) ([]DAActivityTargetResourceModifiedProperty, error) {
	targetResourceModifiedProperties := make([]DAActivityTargetResourceModifiedProperty, 0)
	if modelValue.DirectoryAudit.GetAdditionalData() == nil {
		return targetResourceModifiedProperties, nil
	}

	for _, resource := range modelValue.DirectoryAudit.GetTargetResources() {
		modifiedProperties := resource.GetModifiedProperties()
		if modifiedProperties == nil {
			continue
		}

		for _, mp := range modifiedProperties {
			modifiedProperty := DAActivityTargetResourceModifiedProperty{
				DisplayName: getStringValue(mp.GetDisplayName()),
				NewValue:    getStringValue(mp.GetNewValue()),
				OldValue:    getStringValue(mp.GetOldValue()),
			}

			if additionalData := mp.GetAdditionalData(); additionalData != nil {
				data, err := json.Marshal(additionalData)
				if err != nil {
					return nil, err
				}
				modifiedProperty.AdditionalData = string(data)
			}

			targetResourceModifiedProperties = append(targetResourceModifiedProperties, modifiedProperty)
		}
	}

	return targetResourceModifiedProperties, nil
}

func getStringValue(str *string) string {
	if str != nil {
		return *str
	}
	return ""
}

func WriteDirectoryAuditsToSQL() (int, error) {
	var count int = 0
	for modelKey, modelValue := range directoryAuditsMap {
		activity, err := processDAActivity(modelValue)
		if err != nil {
			return 0, err
		}
		activityAdditionalDetails, err := processDAActivityAdditionalDetails(modelValue)
		if err != nil {
			return 0, err
		}
		activityInitiatedBy, err := processDAActivityInitiatedBy(modelValue)
		if err != nil {
			return 0, err
		}
		activityTargetResources, err := processDAActivityTargetResources(modelValue)
		if err != nil {
			return 0, err
		}
		targetResourceModifiedProperties, err := processDAActivityTargetResourcesModifiedProperties(modelValue)
		if err != nil {
			return 0, err
		}
		if !modelValue.hasModelBeenWritten {
			tx := auditrDb.Begin()
			if err := tx.Error; err != nil {
				return 0, err
			}
			if err := tx.Create(&activity).Error; err != nil {
				tx.Rollback()
				return 0, err
			}
			for _, detail := range activityAdditionalDetails {
				detail.DAActivityID = activity.ID
				if err := tx.Create(&detail).Error; err != nil {
					tx.Rollback()
					return 0, err
				}
			}
			activityInitiatedBy.DAActivityID = activity.ID
			if err := tx.Create(&activityInitiatedBy).Error; err != nil {
				tx.Rollback()
				return 0, err
			}
			for _, targetResource := range activityTargetResources {
				targetResource.DAActivityID = activity.ID
				if err := tx.Create(&targetResource).Error; err != nil {
					tx.Rollback()
					return 0, err
				}
				for _, modifiedProperty := range targetResourceModifiedProperties {
					modifiedProperty.DAActivityID = activity.ID
					modifiedProperty.DATargetResourcesID = targetResource.ID
					if err := tx.Create(&modifiedProperty).Error; err != nil {
						tx.Rollback()
						return 0, err
					}
				}
			}
			if err := tx.Commit().Error; err != nil {
				fmt.Printf("error committing directory audit transactions to sql: %v  ... rolling back\n", err)
				errRollback := tx.Rollback()
				if errRollback != nil {
					return 0, errRollback.Error
				}
				return 0, err
			}
			modelValue.hasModelBeenWritten = true
			directoryAuditsMap[modelKey] = modelValue
			count++
			activity = DAActivity{}
		}
	}
	return count, nil
}
