package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"parentscontactform/internal/auth"
	"parentscontactform/internal/client"
	"parentscontactform/internal/middleware"
	"parentscontactform/internal/models"
	"parentscontactform/internal/session"
	"strings"
)

func HandleFormGet(w http.ResponseWriter, r *http.Request) {
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

	t, err := template.ParseFiles("C:\\Users\\Patrick\\GolandProjects\\parentscontactform\\cmd\\server\\templates\\form.gohtml")
	if err != nil {
		middleware.LogAndError(r, w, "failed to parse gohtml", err.Error(), http.StatusInternalServerError)
		return
	}

	templateData := models.TemplateData{
		CurrentUserInfo: currentUserInfo,
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

func HandleFormPost(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	cookie, err := r.Cookie("session_id")
	if err != nil {
		// No cookie, not logged in
		http.Redirect(w, r, "/login", http.StatusFound)
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
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	hc := auth.Cfg.Client(ctx, sessionData.ISAMSToken.Token)

	c, err := client.NewClientWithResponses(client.BASE_URL, client.WithHTTPClient(hc))
	if err != nil {
		middleware.LogAndError(r, w, "Unable to connect to iSAMS", err.Error(), http.StatusInternalServerError)
	}

	parentChildren, err := client.GetUserChildren(sessionData)
	contactInfo, err := client.GetContactInfo(sessionData.ISAMSSessionId)

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

	resp, err := c.PatchApiPortalsAccountStudentsSchoolIdContactsContactIdWithResponse(ctx, parentChildren[0].SchoolId, contactInfo[0].Id, nil, updatedContactBody)
	if err != nil {
		middleware.LogAndError(r, w, "Unable to update iSAMS record", err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, string(resp.Body))
}
