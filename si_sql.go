package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type SIActivity struct {
	ID                               int    `gorm:"column:id"`
	EventID                          string `gorm:"column:si_event_id"`
	AppDisplayName                   string `gorm:"column:app_display_name"`
	AppID                            string `gorm:"column:app_id"`
	AppliedConditionalAccessPolicies string `gorm:"column:applied_conditional_access_policies"`
	ClientAppUsed                    string `gorm:"column:client_app_used"`
	ConditionalAccessStatus          string `gorm:"column:conditional_access_status"`
	CorrelationID                    string `gorm:"column:correlation_id"`
	CreatedDateTime                  string `gorm:"column:created_date_time"`
	DeviceDetail                     string `gorm:"column:device_detail"`
	IPAddress                        string `gorm:"column:ip_address"`
	IsInteractive                    string `gorm:"column:is_interactive"`
	Location                         string `gorm:"column:location"`
	ResourceDisplayName              string `gorm:"column:resource_display_name"`
	ResourceID                       string `gorm:"column:resource_id"`
	RiskDetail                       string `gorm:"column:risk_detail"`
	RiskEventTypes                   string `gorm:"column:risk_event_types"`
	RiskEventTypesV2                 string `gorm:"column:risk_event_types_v2"`
	RiskLevelAggregated              string `gorm:"column:risk_level_aggregated"`
	RiskLevelDuringSignIn            string `gorm:"column:risk_level_during_sign_in"`
	RiskState                        string `gorm:"column:risk_state"`
	Status                           string `gorm:"column:status"`
	UserDisplayName                  string `gorm:"column:user_display_name"`
	UserID                           string `gorm:"column:user_id"`
	UserPrincipalName                string `gorm:"column:user_principal_name"`
}

func (SIActivity) TableName() string {
	return "si_activities"
}

type SIActivityDeviceDetails struct {
	ID              int    `gorm:"primaryKey;column:id"`
	SIActivityID    int    `gorm:"column:si_activity_id;references:si_activity_id"`
	AdditionalData  string `gorm:"column:additional_data"`
	Browser         string `gorm:"column:browser"`
	DeviceID        string `gorm:"column:device_id"`
	DisplayName     string `gorm:"column:display_name"`
	IsCompliant     bool   `gorm:"column:is_compliant"`
	IsManaged       bool   `gorm:"column:is_managed"`
	ODataType       string `gorm:"column:o_data_type"`
	OperatingSystem string `gorm:"column:operating_system"`
	TrustType       string `gorm:"column:trust_type"`
}

func (SIActivityDeviceDetails) TableName() string {
	return "si_activities_device_details"
}

type SIActivityLocation struct {
	ID              int    `gorm:"primaryKey;column:id"`
	SIActivityID    int    `gorm:"column:si_activity_id;references:si_activity_id"`
	City            string `gorm:"column:city"`
	CountryOrRegion string `gorm:"column:country_or_region"`
	State           string `gorm:"column:state"`
}

func (SIActivityLocation) TableName() string {
	return "si_activities_location"
}

type SIActivityGeoCoordinates struct {
	ID             int    `gorm:"primaryKey;column:id"`
	SIActivityID   int    `gorm:"column:si_activity_id;references:si_activity_id"`
	SILocationID   int    `gorm:"column:si_location_id;references:si_location_id"`
	AdditionalData string `gorm:"column:additional_data"`
	Latitude       string `gorm:"column:latitude"`
	Longitude      string `gorm:"column:longitude"`
}

func (SIActivityGeoCoordinates) TableName() string {
	return "si_activities_geocoordinates"
}

type SIActivityStatus struct {
	ID                int    `gorm:"primaryKey;id"`
	SIActivityID      int    `gorm:"column:si_activity_id;references:si_activity_id"`
	AdditionalData    string `gorm:"column:additional_data"`
	AdditionalDetails string `gorm:"column:additional_detail"`
	ErrorCode         string `gorm:"column:error_code"`
	FailureReason     string `gorm:"column:failure_reason"`
}

func (SIActivityStatus) TableName() string {
	return "si_activities_status"
}

func processSIActivity(modelValue SignInWrite) (SIActivity, error) {
	var activity SIActivity
	signIn := modelValue.SignIn

	if signIn.GetId() != nil {
		activity.EventID = *signIn.GetId()
	}
	if signIn.GetAppDisplayName() != nil {
		activity.AppDisplayName = *signIn.GetAppDisplayName()
	}
	if signIn.GetAppId() != nil {
		activity.AppID = *signIn.GetAppId()
	}
	if policies := signIn.GetAppliedConditionalAccessPolicies(); policies != nil {
		p, err := json.Marshal(policies)
		if err != nil {
			return SIActivity{}, err
		}
		activity.AppliedConditionalAccessPolicies = string(p)
	}
	if signIn.GetClientAppUsed() != nil {
		activity.ClientAppUsed = *signIn.GetClientAppUsed()
	}
	if signIn.GetConditionalAccessStatus() != nil {
		activity.ConditionalAccessStatus = signIn.GetConditionalAccessStatus().String()
	}
	if signIn.GetCorrelationId() != nil {
		activity.CorrelationID = *signIn.GetCorrelationId()
	}
	if signIn.GetCreatedDateTime() != nil {
		activity.CreatedDateTime = signIn.GetCreatedDateTime().Format("2006-01-02 15:04:05.999999 -07:00")
	}
	if signIn.GetIpAddress() != nil {
		activity.IPAddress = *signIn.GetIpAddress()
	}
	if signIn.GetIsInteractive() != nil {
		activity.IsInteractive = strconv.FormatBool(*signIn.GetIsInteractive())
	}
	if signIn.GetResourceDisplayName() != nil {
		activity.ResourceDisplayName = *signIn.GetResourceDisplayName()
	}
	if signIn.GetResourceId() != nil {
		activity.ResourceID = *signIn.GetResourceId()
	}
	if modelValue.SignIn.GetRiskDetail() != nil {
		activity.RiskDetail = modelValue.SignIn.GetRiskDetail().String()
	}
	if riskEventTypes := modelValue.SignIn.GetRiskEventTypes(); riskEventTypes != nil {
		var err error
		r, err := json.Marshal(riskEventTypes)
		if err != nil {
			return SIActivity{}, err
		}
		activity.RiskEventTypes = string(r)
	}
	if riskEventTypesV2 := modelValue.SignIn.GetRiskEventTypesV2(); riskEventTypesV2 != nil {
		var err error
		r, err := json.Marshal(modelValue.SignIn.GetRiskEventTypesV2())
		if err != nil {
			return SIActivity{}, err
		}
		activity.RiskEventTypesV2 = string(r)
	}
	if modelValue.SignIn.GetRiskLevelAggregated() != nil {
		activity.RiskLevelAggregated = modelValue.SignIn.GetRiskLevelAggregated().String()
	}
	if modelValue.SignIn.GetRiskLevelDuringSignIn() != nil {
		activity.RiskLevelDuringSignIn = modelValue.SignIn.GetRiskLevelDuringSignIn().String()
	}
	if modelValue.SignIn.GetRiskState() != nil {
		activity.RiskState = modelValue.SignIn.GetRiskState().String()
	}
	if modelValue.SignIn.GetUserDisplayName() != nil {
		activity.UserDisplayName = *modelValue.SignIn.GetUserDisplayName()
	}
	if modelValue.SignIn.GetUserId() != nil {
		activity.UserID = *modelValue.SignIn.GetUserId()
	}
	if modelValue.SignIn.GetUserPrincipalName() != nil {
		activity.UserPrincipalName = *modelValue.SignIn.GetUserPrincipalName()
	}

	return activity, nil
}

func processSIActivityDeviceDetails(modelValue SignInWrite) (SIActivityDeviceDetails, error) {
	var device SIActivityDeviceDetails
	signInDevice := modelValue.SignIn.GetDeviceDetail()

	if signInDevice == nil {
		return device, nil
	}
	if additionalData := signInDevice.GetAdditionalData(); additionalData != nil {
		d, err := json.Marshal(additionalData)
		if err != nil {
			return SIActivityDeviceDetails{}, err
		}
		device.AdditionalData = string(d)
	}
	if browser := signInDevice.GetBrowser(); browser != nil {
		b, err := json.Marshal(browser)
		if err != nil {
			return SIActivityDeviceDetails{}, err
		}
		device.Browser = string(b)
	}
	if signInDevice.GetDeviceId() != nil {
		device.DeviceID = *signInDevice.GetDeviceId()
	}
	if signInDevice.GetDisplayName() != nil {
		device.DisplayName = *signInDevice.GetDisplayName()
	}
	if signInDevice.GetIsCompliant() != nil {
		device.IsCompliant = *signInDevice.GetIsCompliant()
	}
	if signInDevice.GetIsManaged() != nil {
		device.IsManaged = *signInDevice.GetIsManaged()
	}
	if signInDevice.GetOdataType() != nil {
		device.ODataType = *signInDevice.GetOdataType()
	}
	if signInDevice.GetOperatingSystem() != nil {
		device.OperatingSystem = *signInDevice.GetOperatingSystem()
	}
	if signInDevice.GetTrustType() != nil {
		device.TrustType = *signInDevice.GetTrustType()
	}

	return device, nil
}

func processSIActivityLocation(modelValue SignInWrite) (SIActivityLocation, error) {
	var location SIActivityLocation
	signInLocation := modelValue.SignIn.GetLocation()

	if signInLocation == nil {
		return SIActivityLocation{}, nil
	}
	if city := signInLocation.GetCity(); city != nil {
		location.City = *city
	}
	if countryOrRegion := signInLocation.GetCountryOrRegion(); countryOrRegion != nil {
		location.CountryOrRegion = *countryOrRegion
	}
	if state := signInLocation.GetState(); state != nil {
		location.State = *state
	}

	return location, nil
}

func processSIActivityGeoCoordinates(modelValue SignInWrite) (SIActivityGeoCoordinates, error) {
	var coordinates SIActivityGeoCoordinates
	signInLocation := modelValue.SignIn.GetLocation()
	signInGeoCoordinates := signInLocation.GetGeoCoordinates()

	if signInGeoCoordinates == nil {
		return coordinates, nil
	}
	if additionalData := signInGeoCoordinates.GetAdditionalData(); additionalData != nil {
		g, err := json.Marshal(additionalData)
		if err != nil {
			return SIActivityGeoCoordinates{}, err
		}
		coordinates.AdditionalData = string(g)
	}
	if latitude := signInGeoCoordinates.GetLatitude(); latitude != nil {
		coordinates.Latitude = strconv.FormatFloat(*latitude, 'f', -1, 64)
	}
	if longitude := signInGeoCoordinates.GetLongitude(); longitude != nil {
		coordinates.Longitude = strconv.FormatFloat(*longitude, 'f', -1, 64)
	}

	return coordinates, nil
}

func processSIActivityStatus(modelValue SignInWrite) (SIActivityStatus, error) {
	var status SIActivityStatus
	signInStatus := modelValue.SignIn.GetStatus()

	if signInStatus == nil {
		return SIActivityStatus{}, nil
	}
	if additionalDetails := signInStatus.GetAdditionalDetails(); additionalDetails != nil {
		status.AdditionalDetails = *additionalDetails
	}
	if errorCode := signInStatus.GetErrorCode(); errorCode != nil {
		status.ErrorCode = strconv.Itoa(int(*errorCode))
	}
	if failureReason := signInStatus.GetFailureReason(); failureReason != nil {
		status.FailureReason = *failureReason
	}

	return status, nil
}

func WriteSignInActivityToSQL() (int, error) {
	var count = 0
	for modelKey, modelValue := range signInActivityMap {
		activity, err := processSIActivity(modelValue)
		if err != nil {
			return 0, err
		}
		device, err := processSIActivityDeviceDetails(modelValue)
		if err != nil {
			return 0, err
		}
		location, err := processSIActivityLocation(modelValue)
		if err != nil {
			return 0, err
		}
		coordinates, err := processSIActivityGeoCoordinates(modelValue)
		if err != nil {
			return 0, err
		}
		status, err := processSIActivityStatus(modelValue)
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

			device.SIActivityID = activity.ID
			if err := tx.Create(&device).Error; err != nil {
				tx.Rollback()
				return 0, err
			}
			location.SIActivityID = activity.ID
			if err := tx.Create(&location).Error; err != nil {
				tx.Rollback()
				return 0, err
			}
			coordinates.SIActivityID = activity.ID
			coordinates.SILocationID = location.ID
			if err := tx.Create(&coordinates).Error; err != nil {
				tx.Rollback()
				return 0, err
			}
			status.SIActivityID = activity.ID
			if err := tx.Create(&status).Error; err != nil {
				tx.Rollback()
				return 0, err
			}
			if err := tx.Commit().Error; err != nil {
				fmt.Printf("error committing sign-in transactions to sql: %v  ... rolling back\n", err)
				errRollback := tx.Rollback()
				if errRollback != nil {
					return 0, errRollback.Error
				}
				return 0, err
			}
			modelValue.hasModelBeenWritten = true
			signInActivityMap[modelKey] = modelValue
			count++
			activity = SIActivity{}
		}
	}
	return count, nil
}
