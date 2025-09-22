package models

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
}

type ChildInfo struct {
	SchoolId          string      `json:"schoolId"`
	Forename          string      `json:"forename"`
	Surname           string      `json:"surname"`
	PreferredName     string      `json:"preferredName"`
	Gender            string      `json:"gender"`
	DateOfBirth       string      `json:"dateOfBirth"`
	AcademicHouse     *string     `json:"academicHouse"`
	AcademicHouseCode *string     `json:"academicHouseCode"`
	BoardingHouse     interface{} `json:"boardingHouse"`
	BoardingHouseCode interface{} `json:"boardingHouseCode"`
	FormCode          *string     `json:"formCode"`
	Year              *int        `json:"year"`
	YearCode          *string     `json:"yearCode"`
	YearName          *string     `json:"yearName"`
	SchoolDivision    *int        `json:"schoolDivision"`
	Mobile            interface{} `json:"mobile"`
	Email             interface{} `json:"email"`
	Username          interface{} `json:"username"`
	Status            int         `json:"status"`
	Contacts          interface{} `json:"contacts"`
}
