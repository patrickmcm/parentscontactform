package util

import (
	"crypto/rand"
	"encoding/base64"
	"parentscontactform/internal/client"
	"parentscontactform/internal/models"
)

func GenerateRandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func PtrEq(x *string, y string) bool {
	if x == nil {
		return false
	}

	if *x == y {
		return true
	}

	return false
}

func InitStudent(currentStudent models.GetStudent) client.PutApiStudentsSchoolIdJSONRequestBody {
	newStudent := client.PutApiStudentsSchoolIdJSONRequestBody{
		AcademicHouse:        currentStudent.AcademicHouse,
		BirthCounty:          currentStudent.BirthCounty,
		Birthplace:           currentStudent.Birthplace,
		BoardingHouse:        currentStudent.BoardingHouse,
		BoardingStatus:       currentStudent.BoardingStatus,
		Dob:                  currentStudent.Dob,
		EnrolmentDate:        currentStudent.EnrolmentDate,
		EnrolmentTerm:        currentStudent.EnrolmentTerm,
		EnrolmentYear:        currentStudent.EnrolmentYear,
		Ethnicity:            currentStudent.Ethnicity,
		Forename:             currentStudent.Forename,
		FormGroup:            currentStudent.FormGroup,
		FullName:             currentStudent.FullName,
		FutureSchoolId:       currentStudent.FutureSchoolId,
		Gender:               currentStudent.Gender,
		Initials:             currentStudent.Initials,
		IsVisaRequired:       currentStudent.IsVisaRequired,
		LabelSalutation:      currentStudent.LabelSalutation,
		Languages:            currentStudent.Languages,
		LeavingDate:          currentStudent.LeavingDate,
		LeavingReason:        currentStudent.LeavingReason,
		LeavingYearGroup:     currentStudent.LeavingYearGroup,
		LetterSalutation:     currentStudent.LetterSalutation,
		MiddleNames:          currentStudent.Middlenames,
		Middlenames:          currentStudent.Middlenames,
		MobileNumber:         currentStudent.MobileNumber,
		Nationalities:        currentStudent.Nationalities,
		OfficialName:         currentStudent.OfficialName,
		PersonalEmailAddress: currentStudent.PersonalEmailAddress,
		PreferredName:        currentStudent.PreferredName,
		PreviousName:         currentStudent.PreviousName,
		Religion:             currentStudent.Religion,
		ResidentCountry:      currentStudent.ResidentCountry,
		SchoolCode:           currentStudent.SchoolCode,
		SchoolEmailAddress:   currentStudent.SchoolEmailAddress,
		Surname:              currentStudent.Surname,
		Title:                currentStudent.Title,
		TutorEmployeeId:      currentStudent.TutorEmployeeId,
		UniquePupilNumber:    currentStudent.UniquePupilNumber,
		YearGroup:            currentStudent.YearGroup,
	}

	return newStudent
}
