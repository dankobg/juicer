package server

import (
	"net/http"
)

func (a *ApiHandler) registrationAfterPassword(w http.ResponseWriter, r *http.Request) {
	a.wh.RegistrationAfterPassword(w, r)
}

func (a *ApiHandler) registrationAfterOidc(w http.ResponseWriter, r *http.Request) {
	a.wh.RegistrationAfterOidc(w, r)
}
