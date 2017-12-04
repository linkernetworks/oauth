package app

import (
	"net/http"

	"bitbucket.org/linkernetworks/aurora/src/oauth/validator"
)

func UserLogout(w http.ResponseWriter, r *http.Request, appService *ServiceProvider) {
	locale := MustLocaleFunc(r)
	if len(r.URL.Query()["user_email"]) < 1 {
		toJson(w, FormActionResponse{
			Error:   true,
			Message: locale(ErrParseQuery),
		})
		return
	}
	userEmail := r.URL.Query()["user_email"][0]

	var validations = validator.ValidationMap{}
	emailValidate, err := validator.ValidateEmail(userEmail)
	if err != nil {
		validations["email"] = emailValidate
	}

	if validations.HasError() {
		toJson(w, FormActionResponse{
			Error:       true,
			Validations: validations,
		})
		return
	}

	session, err := appService.SessionStore.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	delete(session.Values, userEmail)
	session.Save(r, w)

	// logout successfully,
	// redirect to /signin
	http.Redirect(w, r, "/signin", http.StatusMovedPermanently)
}
