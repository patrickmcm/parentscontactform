package handlers

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"parentscontactform/internal/auth"
	"parentscontactform/internal/client"
	"parentscontactform/internal/middleware"
	"parentscontactform/internal/models"
	"parentscontactform/internal/session"
	"parentscontactform/internal/util"
	"strings"
	"time"
)

type Handler struct {
	staticFS   embed.FS
	templateFS embed.FS
	RestClient *client.ClientWithResponses
}

func NewHandler(staticFS embed.FS, templateFS embed.FS, restClient *client.ClientWithResponses) *Handler {
	return &Handler{staticFS, templateFS, restClient}
}

func (h *Handler) HandleFormGet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// 1. Get the session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// No cookie, not logged in
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// 2. Look up the session from memory
	sessionData, exists := session.Get(cookie.Value)

	if !exists {
		// Session not found (e.g., server was restarted)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	parentInfo, err := client.GetContactInfo(sessionData.ISAMSSessionId)
	if err != nil {
		middleware.LogAndError(r, w, "Error retriving your contact info", err.Error(), http.StatusInternalServerError)
		return
	}

	currentUserInfo, err := client.GetUserAccountDetails(sessionData)
	if err != nil {
		middleware.LogAndError(r, w, "Error retriving your account info", err.Error(), http.StatusInternalServerError)
		return
	}

	parentChildren, err := client.GetUserChildren(sessionData)
	if err != nil {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
	}

	/*
		NEXT STEPS:
			1. Get each childs med info (conditions, general consent)
			2. Fetch custom fields for photos consent
			3. Fetch if EAL
	*/

	medicalConditionsList, err := h.RestClient.GetApiMedicalConditiontypesWithResponse(ctx, nil)
	if err != nil {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
		return
	}

	languagesList, err := h.RestClient.GetApiSystemconfigurationListLanguagesWithResponse(ctx, nil)
	if err != nil {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
		return
	}

	if languagesList == nil {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := models.TemplateData{
		CurrentUserInfo: currentUserInfo,
		ChildrenInfo:    parentChildren,
		ConditionTypes:  *medicalConditionsList.JSON200.ConditionTypes,
		Languages:       *languagesList.JSON200.Items,
	}

	funcMap := template.FuncMap{
		"ptrEq": util.PtrEq,
	}

	t, err := template.New("form.gohtml").Funcs(funcMap).ParseFS(h.templateFS, "templates/form.gohtml")
	if err != nil {
		middleware.LogAndError(r, w, "failed to parse gohtml", err.Error(), http.StatusInternalServerError)
		return
	}

	if len(parentInfo) > 0 {
		templateData.ParentContactInfo = parentInfo[0]
	}

	err = t.Execute(w, templateData)
	if err != nil {
		middleware.LogAndError(r, w, "failed to execute template", err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) HandleChildFormGet(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	// 1. Get the session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// No cookie, not logged in
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	// 2. Look up the session from memory
	sessionData, exists := session.Get(cookie.Value)

	if !exists {
		// Session not found (e.g., server was restarted)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	parentChildren, err := client.GetUserChildren(sessionData)
	if err != nil {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
	}

	childrenInfo := make(map[string]models.ChildFormInfo)

	for _, child := range parentChildren {
		var childFormInfo models.ChildFormInfo

		if child.Status < 1 {
			continue
		}

		childrenConditionsInfo, err := h.RestClient.GetApiMedicalStudentsSchoolIdConditionsWithResponse(ctx, child.SchoolId, nil)
		if err != nil {
			middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
			return
		}

		if childrenConditionsInfo.JSON200.Conditions != nil {
			childFormInfo.Conditions = *childrenConditionsInfo.JSON200.Conditions
		}

		photosConsent, err := h.RestClient.GetApiStudentsSchoolIdCustomFieldsCustomFieldIdWithResponse(ctx, child.SchoolId, 27, nil)
		if err != nil {
			middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
			return
		}

		if photosConsent.JSON200 != nil {
			customField := *(*photosConsent.JSON200.CustomFields)[0].Value
			if customField == "true" {
				childFormInfo.PhotoConsent = true
			} else {
				childFormInfo.PhotoConsent = false
			}
		}

		medConsentResp, err := h.RestClient.GetApiMedicalStudentsSchoolIdConsentsWithResponse(ctx, child.SchoolId, nil)
		if err != nil {
			middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
			return
		}

		if medConsentResp.JSON200 != nil {
			consents := *medConsentResp.JSON200.Consents
			for _, consent := range consents {
				if *consent.TypeId == 18 {
					childFormInfo.MedConsent = *consent.ConsentGiven
				}
			}
		}

		isEalResp, err := h.RestClient.GetApiStudentsSchoolIdCustomFieldsCustomFieldIdWithResponse(ctx, child.SchoolId, 35, nil)
		if err != nil {
			middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
			return
		}

		if isEalResp.JSON200 != nil {
			customField := *(*isEalResp.JSON200.CustomFields)[0].Value
			if customField == "true" {
				childFormInfo.IsEal = true
			} else {
				childFormInfo.IsEal = false
			}
		}

		tripsConsentResp, err := h.RestClient.GetApiStudentsSchoolIdCustomFieldsCustomFieldIdWithResponse(ctx, child.SchoolId, 29, nil)
		if err != nil {
			middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
			return
		}

		if tripsConsentResp.JSON200 != nil {
			customField := *(*tripsConsentResp.JSON200.CustomFields)[0].Value
			if customField == "true" {
				childFormInfo.TripsConsent = true
			} else {
				childFormInfo.TripsConsent = false
			}
		}

		studentResp, err := h.RestClient.GetApiStudentsSchoolIdWithResponse(ctx, child.SchoolId, nil)
		if err != nil {
			middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
			return
		}

		if studentResp.JSON200 != nil {
			childFormInfo.Languages = *studentResp.JSON200.Languages
		}

		childrenInfo[child.SchoolId] = childFormInfo
	}

	encoded, err := json.Marshal(childrenInfo)
	if err != nil {
		middleware.LogAndError(r, w, "Unable to access your children's info", err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(encoded)
}

func (h *Handler) HandleChildFormPost(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	cookie, err := r.Cookie("session_id")
	if err != nil {
		// No cookie, not logged in
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	_, exists := session.Get(cookie.Value)

	if !exists {
		// Session not found (e.g., server was restarted)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		middleware.LogAndError(r, w, "malformed body", err.Error(), http.StatusInternalServerError)
		return
	}

	var childInfo models.ChildFormInfoRequestBody
	err = json.Unmarshal(body, &childInfo)
	if err != nil {
		middleware.LogAndError(r, w, "malformed body", err.Error(), http.StatusInternalServerError)
		return
	}

	for _, condition := range childInfo.ToDelete {
		delConditionResp, err := h.RestClient.DeleteApiMedicalStudentsSchoolIdConditionsConditionIdWithResponse(ctx, childInfo.SchoolId, condition.Key, nil)
		if err != nil || delConditionResp == nil {
			middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
			return
		}

		if delConditionResp.StatusCode() != http.StatusOK {
			middleware.LogAndError(r, w, "Unable to delete condition", string(delConditionResp.Body), http.StatusInternalServerError)
			return
		}
	}

	for _, condition := range childInfo.Conditions {
		if condition.ToBeUploaded && condition.Key != "" {
			op := "replace"

			groupIdPath := "/GroupId"
			typePath := "/Type"
			furtherInfoPath := "/FurtherInfo"
			treatmentPath := "/Treatment"

			groupIdStr := string(rune(condition.GroupId))

			updateBody := client.PatchApiMedicalStudentsSchoolIdConditionsConditionIdJSONBody{
				{
					Op:    &op,
					Path:  &groupIdPath,
					Value: &groupIdStr,
				},
				{
					Op:    &op,
					Path:  &typePath,
					Value: &condition.Type,
				},
				{
					Op:    &op,
					Path:  &furtherInfoPath,
					Value: &condition.FurtherInfo,
				},
				{
					Op:    &op,
					Path:  &treatmentPath,
					Value: &condition.Treatment,
				},
			}
			updatedCondResp, err := h.RestClient.PatchApiMedicalStudentsSchoolIdConditionsConditionIdWithResponse(ctx, childInfo.SchoolId, condition.Key, nil, updateBody)
			if err != nil {
				middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
				return
			}

			if updatedCondResp.StatusCode() != 200 {
				middleware.LogAndError(r, w, "Unable to update condition", string(updatedCondResp.Body), http.StatusInternalServerError)
				return
			}
		} else if condition.ToBeUploaded && condition.Key == "" {
			now := time.Now().Format(time.RFC3339)
			newConditionBody := client.PostApiMedicalStudentsSchoolIdConditionsJSONRequestBody{
				DateEntered: &now,
				FurtherInfo: &condition.FurtherInfo,
				GroupId:     &condition.GroupId,
				Treatment:   &condition.Treatment,
				Type:        &condition.Type,
			}
			newConditionResp, err := h.RestClient.PostApiMedicalStudentsSchoolIdConditionsWithResponse(ctx, childInfo.SchoolId, nil, newConditionBody)
			if err != nil || newConditionResp.StatusCode() != 201 {
				middleware.LogAndError(r, w, "Unable to connect to iSAMS", string(newConditionResp.Body), http.StatusInternalServerError)
				return
			}
		}
	}

	photoConsent := "false"
	if childInfo.PhotoConsent {
		photoConsent = "true"
	}

	updatedPhotoConsent := client.PatchApiStudentsSchoolIdCustomFieldsCustomFieldIdJSONRequestBody{
		Value: &photoConsent,
	}

	photoConsentResp, err := h.RestClient.PatchApiStudentsSchoolIdCustomFieldsCustomFieldIdWithResponse(ctx, childInfo.SchoolId, 27, nil, updatedPhotoConsent)
	if err != nil {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
		return
	}

	if photoConsentResp.StatusCode() != 200 {
		middleware.LogAndError(r, w, "Unable to update condition", string(photoConsentResp.Body), http.StatusInternalServerError)
		return
	}

	// med consent

	medConsentsResp, err := h.RestClient.GetApiMedicalStudentsSchoolIdConsentsWithResponse(ctx, childInfo.SchoolId, nil)
	if err != nil {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
		return
	}

	if medConsentsResp.JSON200 == nil {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
		return
	}

	medConsent := "false"
	if childInfo.MedConsent {
		medConsent = "true"
	}

	foundConsentRecord := false

	for _, consent := range *medConsentsResp.JSON200.Consents {
		if *consent.TypeId == 18 {
			foundConsentRecord = true
			op := "replace"
			consentGivenPath := "/ConsentGiven"
			updatedMedConsent := client.PatchApiMedicalStudentsSchoolIdConsentsConsentIdJSONBody{
				{
					Op:    &op,
					Path:  &consentGivenPath,
					Value: &medConsent,
				},
			}

			updatedConsentResp, err := h.RestClient.PatchApiMedicalStudentsSchoolIdConsentsConsentIdWithResponse(ctx, childInfo.SchoolId, *consent.Key, nil, updatedMedConsent)
			if err != nil || updatedConsentResp.StatusCode() != 200 {
				middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	if !foundConsentRecord {
		now := time.Now().Format(time.DateOnly)
		typeId := 18
		newConsentBody := client.PostApiMedicalStudentsSchoolIdConsentsJSONRequestBody{
			ConsentGiven: &childInfo.MedConsent,
			ReceivedDate: &now,
			TypeId:       &typeId,
		}
		newConsentResp, err := h.RestClient.PostApiMedicalStudentsSchoolIdConsentsWithResponse(ctx, childInfo.SchoolId, nil, newConsentBody)
		if err != nil || newConsentResp.StatusCode() != 200 {
			middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// eal
	isEal := "false"
	if childInfo.IsEal {
		isEal = "true"
	}

	isEalBody := client.PatchApiStudentsSchoolIdCustomFieldsCustomFieldIdJSONRequestBody{
		Value: &isEal,
	}

	isEalFieldResp, err := h.RestClient.PatchApiStudentsSchoolIdCustomFieldsCustomFieldIdWithResponse(ctx, childInfo.SchoolId, 35, nil, isEalBody)
	if err != nil {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
		return
	}

	if isEalFieldResp.StatusCode() != 200 {
		middleware.LogAndError(r, w, "Unable to update condition", string(isEalFieldResp.Body), http.StatusInternalServerError)
		return
	}

	currentStudentInfoResp, err := h.RestClient.GetApiStudentsSchoolIdWithResponse(ctx, childInfo.SchoolId, nil)
	if err != nil {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
		return
	}

	if currentStudentInfoResp.JSON200 == nil {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", string(currentStudentInfoResp.Body), http.StatusInternalServerError)
		return
	}

	updatedStudentBody := util.InitStudent(currentStudentInfoResp.JSON200)

	for i, v := range childInfo.Languages {
		if v == "English" {
			break
		} else if i == len(childInfo.Languages)-1 {
			childInfo.Languages = append(childInfo.Languages, "English")
		}
	}

	updatedStudentBody.Languages = &childInfo.Languages

	updatedLangsResp, err := h.RestClient.PutApiStudentsSchoolIdWithResponse(ctx, childInfo.SchoolId, nil, updatedStudentBody)
	if err != nil || updatedLangsResp.StatusCode() != 200 {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
		return
	}

	// trips consent

	tripsConsent := "false"
	if childInfo.TripsConsent {
		tripsConsent = "true"
	}

	tripsConsentBody := client.PatchApiStudentsSchoolIdCustomFieldsCustomFieldIdJSONRequestBody{
		Value: &tripsConsent,
	}

	tripsConsentResp, err := h.RestClient.PatchApiStudentsSchoolIdCustomFieldsCustomFieldIdWithResponse(ctx, childInfo.SchoolId, 29, nil, tripsConsentBody)
	if err != nil {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
		return
	}

	if tripsConsentResp.StatusCode() != 200 {
		middleware.LogAndError(r, w, "Unable to update condition", string(tripsConsentResp.Body), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "OK")
}

func (h *Handler) HandleFormPost(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	cookie, err := r.Cookie("session_id")
	if err != nil {
		// No cookie, not logged in
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	const maxMemory = 10 << 20
	if err := r.ParseMultipartForm(maxMemory); err != nil {
		middleware.LogAndError(r, w, "Error parsing form", err.Error(), http.StatusBadRequest)
		return
	}

	sessionData, exists := session.Get(cookie.Value)

	if !exists {
		// Session not found (e.g., server was restarted)
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	hc := auth.Cfg.Client(ctx, sessionData.ISAMSToken.Token)

	c, err := client.NewClientWithResponses(client.BASE_URL, client.WithHTTPClient(hc))
	if err != nil {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
	}

	parentChildren, err := client.GetUserChildren(sessionData)
	if err != nil {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
	}

	contactInfo, err := client.GetContactInfo(sessionData.ISAMSSessionId)
	if err != nil {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
	}

	formType := r.Form.Get("contactType")

	updatedContactBody := client.PatchApiPortalsAccountStudentsSchoolIdContactsContactIdJSONRequestBody{}
	op := "replace"

	for k, v := range r.Form {
		switch strings.TrimPrefix(k, "secondary") {
		case "forename":
			trimmed := strings.TrimSpace(v[0])
			path := "/forename"
			if formType == "secondaryContactForm" {
				path = "/secondaryForename"
			}
			operation := client.PatchApiPortalsAccountStudentsSchoolIdContactsContactIdJSONRequestBody{
				{
					Op:    &op,
					Path:  &path,
					Value: &trimmed,
				},
			}

			updatedContactBody = append(updatedContactBody, operation[0])
		case "surname":
			trimmed := strings.TrimSpace(v[0])
			path := "/surname"
			if formType == "secondaryContactForm" {
				path = "/secondarySurname"
			}
			operation := client.PatchApiPortalsAccountStudentsSchoolIdContactsContactIdJSONRequestBody{
				{
					Op:    &op,
					Path:  &path,
					Value: &trimmed,
				},
			}

			updatedContactBody = append(updatedContactBody, operation[0])
		case "title":
			trimmed := strings.TrimSpace(v[0])
			path := "/title"
			if formType == "secondaryContactForm" {
				path = "/secondaryTitle"
			}
			operation := client.PatchApiPortalsAccountStudentsSchoolIdContactsContactIdJSONRequestBody{
				{
					Op:    &op,
					Path:  &path,
					Value: &trimmed,
				},
			}

			updatedContactBody = append(updatedContactBody, operation[0])
		case "addressLineOne":
			path := "/address1"
			operation := client.PatchApiPortalsAccountStudentsSchoolIdContactsContactIdJSONRequestBody{
				{
					Op:    &op,
					Path:  &path,
					Value: &v[0],
				},
			}

			updatedContactBody = append(updatedContactBody, operation[0])
		case "addressLineTwo":
			path := "/address2"
			operation := client.PatchApiPortalsAccountStudentsSchoolIdContactsContactIdJSONRequestBody{
				{
					Op:    &op,
					Path:  &path,
					Value: &v[0],
				},
			}

			updatedContactBody = append(updatedContactBody, operation[0])
		case "addressLineThree":
			path := "/address3"
			operation := client.PatchApiPortalsAccountStudentsSchoolIdContactsContactIdJSONRequestBody{
				{
					Op:    &op,
					Path:  &path,
					Value: &v[0],
				},
			}

			updatedContactBody = append(updatedContactBody, operation[0])
		case "postalCode":
			path := "/postCode"
			operation := client.PatchApiPortalsAccountStudentsSchoolIdContactsContactIdJSONRequestBody{
				{
					Op:    &op,
					Path:  &path,
					Value: &v[0],
				},
			}

			updatedContactBody = append(updatedContactBody, operation[0])
		case "county":
			path := "/county"
			operation := client.PatchApiPortalsAccountStudentsSchoolIdContactsContactIdJSONRequestBody{
				{
					Op:    &op,
					Path:  &path,
					Value: &v[0],
				},
			}

			updatedContactBody = append(updatedContactBody, operation[0])
		case "phone":
			path := "/mobilePhone"
			if formType == "secondaryContactForm" {
				path = "/secondaryMobilePhone"
			}
			operation := client.PatchApiPortalsAccountStudentsSchoolIdContactsContactIdJSONRequestBody{
				{
					Op:    &op,
					Path:  &path,
					Value: &v[0],
				},
			}

			updatedContactBody = append(updatedContactBody, operation[0])
		}
	}

	updateStudentResp, err := c.PatchApiPortalsAccountStudentsSchoolIdContactsContactIdWithResponse(ctx, parentChildren[0].SchoolId, contactInfo[0].Id, nil, updatedContactBody)
	if err != nil {
		middleware.LogAndError(r, w, "Unable to update iSAMS record", err.Error(), http.StatusInternalServerError)
		return
	}

	if updateStudentResp.StatusCode() != http.StatusOK {
		middleware.LogAndError(r, w, "You are not linked to your child's account. Please ensure you can see your child at fideliscollege.parents.isams.cloud", string(updateStudentResp.Body), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "OK")
}
