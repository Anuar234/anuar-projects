package handlers

import (
	"net/http"
)

func SearchHandler(w http.ResponseWriter, r *http.Request) {
    city := r.URL.Query().Get("city")
    mode := r.URL.Query().Get("mode")

    if city == "" {
        http.Redirect(w, r, "/", http.StatusSeeOther)
        return
    }

    if mode == "forecast" {
        http.Redirect(w, r, "/forecast?city="+city, http.StatusSeeOther)
    } else {
        http.Redirect(w, r, "/weather?city="+city, http.StatusSeeOther)
    }
}
