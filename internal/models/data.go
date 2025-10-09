package models

import "time"

type ParentContactInfo struct {
	Id                      int    `json:"id"`
	Type                    string `json:"type"`
	LabelSalutation         string `json:"labelSalutation"`
	LetterSalutation        string `json:"letterSalutation"`
	RelationType            string `json:"relationType"`
	Deceased                int    `json:"deceased"`
	Title                   string `json:"title"`
	Forename                string `json:"forename"`
	Surname                 string `json:"surname"`
	Initials                string `json:"initials"`
	MiddleNames             string `json:"middleNames"`
	SecondaryDeceased       int    `json:"secondaryDeceased"`
	SecondaryTitle          string `json:"secondaryTitle"`
	SecondaryForename       string `json:"secondaryForename"`
	SecondarySurname        string `json:"secondarySurname"`
	SecondaryInitials       string `json:"secondaryInitials"`
	SecondaryMiddleNames    string `json:"secondaryMiddleNames"`
	FirstLineAddress        string `json:"firstLineAddress"`
	SecondLineAddress       string `json:"secondLineAddress"`
	ThirdLineAddress        string `json:"thirdLineAddress"`
	Town                    string `json:"town"`
	County                  string `json:"county"`
	Postcode                string `json:"postcode"`
	DaytimeTelephone        string `json:"daytimeTelephone"`
	EveningTelephone        string `json:"eveningTelephone"`
	WorkTelephone           string `json:"workTelephone"`
	EmailAddress1           string `json:"emailAddress1"`
	EmailAddress2           string `json:"emailAddress2"`
	Mobile1                 string `json:"mobile1"`
	Mobile2                 string `json:"mobile2"`
	JustContact             int    `json:"justContact"`
	MailMergeAll            int    `json:"mailMergeAll"`
	BillingMailMerge        int    `json:"billingMailMerge"`
	CorrespondenceMailMerge int    `json:"correspondenceMailMerge"`
	ReportsMailMerge        int    `json:"reportsMailMerge"`
	Sos                     string `json:"sos"`
	Ordinal                 int    `json:"ordinal"`
	DualContact             int    `json:"dualContact"`
}

type CurrentUserInfo struct {
	UserName         string      `json:"userName"`
	UserCode         string      `json:"userCode"`
	Title            string      `json:"title"`
	FirstName        string      `json:"firstName"`
	Surname          string      `json:"surname"`
	Fullname         string      `json:"fullname"`
	NickName         interface{} `json:"nickName"`
	Email            string      `json:"email"`
	Website          interface{} `json:"website"`
	Notes            string      `json:"notes"`
	PrimaryRole      string      `json:"primaryRole"`
	DirectLine       string      `json:"directLine"`
	Extension        string      `json:"extension"`
	Address          string      `json:"address"`
	PostCode         string      `json:"postCode"`
	HomeTelephone    string      `json:"homeTelephone"`
	MobileNumber     string      `json:"mobileNumber"`
	HomeEmail        string      `json:"homeEmail"`
	Homepage         string      `json:"homepage"`
	HidePersonal     int         `json:"hidePersonal"`
	OtherInformation string      `json:"otherInformation"`
	HideOtherInfo    int         `json:"hideOtherInfo"`
	UserType         string      `json:"userType"`
	Gender           string      `json:"gender"`
	BirthDate        interface{} `json:"birthDate"`
	UserTypeId       interface{} `json:"userTypeId"`
}

type TemplateData struct {
	ParentContactInfo ParentContactInfo
	CurrentUserInfo   CurrentUserInfo
	ChildrenInfo      []ChildInfo
	ConditionTypes    ConditionTypes
	Languages         Languages
}

type ChildFormInfo struct {
	Conditions   Conditions `json:"conditions"`
	PhotoConsent bool       `json:"photoConsent"`
	MedConsent   bool       `json:"medConsent"`
	TripsConsent bool       `json:"tripsConsent"`
	IsEal        bool       `json:"isEal"`
	Languages    []string   `json:"languages"`
	ToBeUpdated  bool       `json:"toBeUpdated"`
}

type RequestBodyCondition struct {
	Archived      bool      `json:"archived"`
	DateEntered   string    `json:"dateEntered"`
	FurtherInfo   string    `json:"furtherInfo"`
	GroupId       int       `json:"groupId"`
	Key           string    `json:"key"`
	LastUpdated   time.Time `json:"lastUpdated"`
	LastUpdatedBy string    `json:"lastUpdatedBy"`
	Sensitivity   int       `json:"sensitivity"`
	SeverityId    int       `json:"severityId"`
	Treatment     string    `json:"treatment"`
	Type          string    `json:"type"`
	ToBeUploaded  bool      `json:"toBeUploaded"`
}

type ChildFormInfoRequestBody struct {
	SchoolId     string                 `json:"schoolId"`
	Conditions   []RequestBodyCondition `json:"conditions"`
	ToDelete     []RequestBodyCondition `json:"toDelete"`
	PhotoConsent bool                   `json:"photoConsent"`
	MedConsent   bool                   `json:"medConsent"`
	TripsConsent bool                   `json:"tripsConsent"`
	IsEal        bool                   `json:"isEal"`
	Languages    []string               `json:"languages"`
}

type ChildInfo struct {
	SchoolId          string `json:"schoolId"`
	Forename          string `json:"forename"`
	Surname           string `json:"surname"`
	PreferredName     string `json:"preferredName"`
	Gender            string `json:"gender"`
	DateOfBirth       string `json:"dateOfBirth"`
	AcademicHouse     string `json:"academicHouse"`
	AcademicHouseCode string `json:"academicHouseCode"`
	BoardingHouse     string `json:"boardingHouse"`
	BoardingHouseCode string `json:"boardingHouseCode"`
	FormCode          string `json:"formCode"`
	Year              int    `json:"year"`
	YearCode          string `json:"yearCode"`
	YearName          string `json:"yearName"`
	SchoolDivision    int    `json:"schoolDivision"`
	Mobile            string `json:"mobile"`
	Email             string `json:"email"`
	Username          string `json:"username"`
	Status            int    `json:"status"`
	Contacts          []struct {
		Id                      int    `json:"id"`
		Type                    string `json:"type"`
		LabelSalutation         string `json:"labelSalutation"`
		LetterSalutation        string `json:"letterSalutation"`
		RelationType            string `json:"relationType"`
		Deceased                int    `json:"deceased"`
		Title                   string `json:"title"`
		Forename                string `json:"forename"`
		Surname                 string `json:"surname"`
		Initials                string `json:"initials"`
		MiddleNames             string `json:"middleNames"`
		SecondaryDeceased       int    `json:"secondaryDeceased"`
		SecondaryTitle          string `json:"secondaryTitle"`
		SecondaryForename       string `json:"secondaryForename"`
		SecondarySurname        string `json:"secondarySurname"`
		SecondaryInitials       string `json:"secondaryInitials"`
		SecondaryMiddleNames    string `json:"secondaryMiddleNames"`
		FirstLineAddress        string `json:"firstLineAddress"`
		SecondLineAddress       string `json:"secondLineAddress"`
		ThirdLineAddress        string `json:"thirdLineAddress"`
		Town                    string `json:"town"`
		County                  string `json:"county"`
		Postcode                string `json:"postcode"`
		DaytimeTelephone        string `json:"daytimeTelephone"`
		EveningTelephone        string `json:"eveningTelephone"`
		WorkTelephone           string `json:"workTelephone"`
		EmailAddress1           string `json:"emailAddress1"`
		EmailAddress2           string `json:"emailAddress2"`
		Mobile1                 string `json:"mobile1"`
		Mobile2                 string `json:"mobile2"`
		JustContact             int    `json:"justContact"`
		MailMergeAll            int    `json:"mailMergeAll"`
		BillingMailMerge        int    `json:"billingMailMerge"`
		CorrespondenceMailMerge int    `json:"correspondenceMailMerge"`
		ReportsMailMerge        int    `json:"reportsMailMerge"`
		Sos                     string `json:"sos"`
		Ordinal                 int    `json:"ordinal"`
		DualContact             int    `json:"dualContact"`
	} `json:"contacts"`
}

type Conditions []struct {
	Archived      *bool   `json:"archived,omitempty"`
	DateEntered   *string `json:"dateEntered,omitempty"`
	DateReviewed  *string `json:"dateReviewed,omitempty"`
	FurtherInfo   *string `json:"furtherInfo,omitempty"`
	GroupId       *int    `json:"groupId,omitempty"`
	Key           *string `json:"key,omitempty"`
	LastUpdated   *string `json:"lastUpdated,omitempty"`
	LastUpdatedBy *string `json:"lastUpdatedBy,omitempty"`
	Sensitivity   *int    `json:"sensitivity,omitempty"`
	SeverityId    *int    `json:"severityId,omitempty"`
	Treatment     *string `json:"treatment,omitempty"`
	Trigger       *string `json:"trigger,omitempty"`
	Type          *string `json:"type,omitempty"`
}
type CustomFields struct {
	Area                        *string `json:"area,omitempty"`
	DefaultValue                *string `json:"defaultValue,omitempty"`
	Id                          *int    `json:"id,omitempty"`
	InputType                   *string `json:"inputType,omitempty"`
	Name                        *string `json:"name,omitempty"`
	Section                     *string `json:"section,omitempty"`
	SupportsMultipleListItems   *bool   `json:"supportsMultipleListItems,omitempty"`
	SystemConfigurationListName *string `json:"systemConfigurationListName,omitempty"`
	Value                       *string `json:"value,omitempty"`
}

type ConditionTypes []struct {
	Allergy *bool   `json:"allergy,omitempty"`
	Deleted *bool   `json:"deleted,omitempty"`
	Id      *int    `json:"id,omitempty"`
	Name    *string `json:"name,omitempty"`
}

type Languages []struct {
	Common *bool `json:"common,omitempty"`
	Ctf    *struct {
		Code        *string `json:"code,omitempty"`
		Description *string `json:"description,omitempty"`
	} `json:"ctf,omitempty"`
	Description *string `json:"description,omitempty"`
	Dfe         *struct {
		Code        *string `json:"code,omitempty"`
		Description *string `json:"description,omitempty"`
	} `json:"dfe,omitempty"`
	Id  *int `json:"id,omitempty"`
	Isc *struct {
		Code        *string `json:"code,omitempty"`
		Description *string `json:"description,omitempty"`
	} `json:"isc,omitempty"`
	ListType *string `json:"listType,omitempty"`
	Name     *string `json:"name,omitempty"`
}

type GetStudent *struct {
	AcademicHouse   *string `json:"academicHouse,omitempty"`
	BirthCounty     *string `json:"birthCounty,omitempty"`
	Birthplace      *string `json:"birthplace,omitempty"`
	BoardingHouse   *string `json:"boardingHouse,omitempty"`
	BoardingStatus  *string `json:"boardingStatus,omitempty"`
	Dob             *string `json:"dob,omitempty"`
	EnrolmentDate   *string `json:"enrolmentDate,omitempty"`
	EnrolmentStatus *string `json:"enrolmentStatus,omitempty"`
	EnrolmentTerm   *string `json:"enrolmentTerm,omitempty"`
	EnrolmentYear   *int    `json:"enrolmentYear,omitempty"`
	Ethnicity       *string `json:"ethnicity,omitempty"`
	FamilyId        *int    `json:"familyId,omitempty"`
	Forename        *string `json:"forename,omitempty"`
	FormGroup       *string `json:"formGroup,omitempty"`
	FullName        *string `json:"fullName,omitempty"`
	FutureSchoolId  *int    `json:"futureSchoolId,omitempty"`
	Gender          *string `json:"gender,omitempty"`
	HomeAddresses   *[]struct {
		Address1 *string `json:"address1,omitempty"`
		Address2 *string `json:"address2,omitempty"`
		Address3 *string `json:"address3,omitempty"`
		Country  *string `json:"country,omitempty"`
		County   *string `json:"county,omitempty"`
		Id       *int    `json:"id,omitempty"`
		Postcode *string `json:"postcode,omitempty"`
		Private  *bool   `json:"private,omitempty"`
		Town     *string `json:"town,omitempty"`
	} `json:"homeAddresses,omitempty"`
	Id                  *int    `json:"id,omitempty"`
	Initials            *string `json:"initials,omitempty"`
	IsVisaRequired      *bool   `json:"isVisaRequired,omitempty"`
	LabelSalutation     *string `json:"labelSalutation,omitempty"`
	LanguageIsoMappings *[]struct {
		ISAMSName    *string `json:"iSAMSName,omitempty"`
		IsoCode2Char *string `json:"isoCode2Char,omitempty"`
		IsoCode3Char *string `json:"isoCode3Char,omitempty"`
		IsoName      *string `json:"isoName,omitempty"`
	} `json:"languageIsoMappings,omitempty"`
	Languages              *[]string `json:"languages,omitempty"`
	LastUpdated            *string   `json:"lastUpdated,omitempty"`
	LatestPhotoId          *int      `json:"latestPhotoId,omitempty"`
	LeavingDate            *string   `json:"leavingDate,omitempty"`
	LeavingReason          *string   `json:"leavingReason,omitempty"`
	LeavingYearGroup       *int      `json:"leavingYearGroup,omitempty"`
	LetterSalutation       *string   `json:"letterSalutation,omitempty"`
	Middlenames            *string   `json:"middlenames,omitempty"`
	MobileNumber           *string   `json:"mobileNumber,omitempty"`
	Nationalities          *[]string `json:"nationalities,omitempty"`
	NationalityIsoMappings *[]struct {
		ISAMSName    *string `json:"iSAMSName,omitempty"`
		IsoCode2Char *string `json:"isoCode2Char,omitempty"`
		IsoCode3Char *string `json:"isoCode3Char,omitempty"`
		IsoName      *string `json:"isoName,omitempty"`
	} `json:"nationalityIsoMappings,omitempty"`
	OfficialName         *string `json:"officialName,omitempty"`
	PersonGuid           *string `json:"personGuid,omitempty"`
	PersonId             *int    `json:"personId,omitempty"`
	PersonalEmailAddress *string `json:"personalEmailAddress,omitempty"`
	PreferredName        *string `json:"preferredName,omitempty"`
	PreviousName         *string `json:"previousName,omitempty"`
	Religion             *string `json:"religion,omitempty"`
	RemovalGrounds       *struct {
		Code        *string `json:"code,omitempty"`
		Description *string `json:"description,omitempty"`
	} `json:"removalGrounds,omitempty"`
	ResidentCountry    *string `json:"residentCountry,omitempty"`
	SchoolCode         *string `json:"schoolCode,omitempty"`
	SchoolEmailAddress *string `json:"schoolEmailAddress,omitempty"`
	SchoolId           *string `json:"schoolId,omitempty"`
	Surname            *string `json:"surname,omitempty"`
	SystemStatus       *string `json:"systemStatus,omitempty"`
	Title              *string `json:"title,omitempty"`
	TutorEmployeeId    *int    `json:"tutorEmployeeId,omitempty"`
	UniquePupilNumber  *string `json:"uniquePupilNumber,omitempty"`
	YearGroup          *int    `json:"yearGroup,omitempty"`
}
