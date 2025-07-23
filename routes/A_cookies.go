package  routes


import {
	"net/http"
}

func CreateCookie(user_name, user_id string, w http.ResponseWriter, r *http.Request){

	// Validate existing cookie if present
	if cookie, err := r.Cookie("user_info"); err == nil {
		// Cookie exists - verify its value matches current user
		parts := strings.Split(cookie.Value, ":")
		if len(parts) == 2 && parts[0] == user_id && parts[1] == user_name {
			// Valid cookie exists, no need to set a new one
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Login successful"))
			return
		}
		// Cookie exists but doesn't match - we'll overwrite it
	}

	// Set new cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "user_info",
		Value:    fmt.Sprintf("%s:%s", user_id, user_name),
		Expires:  time.Now().Add(30 * 24 * time.Hour), // 30 days
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}


func GetUserName(r *http.Request) (string, error) {
	cookie, err := r.Cookie("user_info")
	if err != nil {
		return "", fmt.Errorf("cookie not found")
	}


	parts := strings.Split(cookie.Value, ":")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid cookie format")
	}

	return parts[2], nil
}
