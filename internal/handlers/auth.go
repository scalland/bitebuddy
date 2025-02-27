package handlers

import (
	"fmt"
	"github.com/scalland/bitebuddy/pkg/utils"
	"github.com/spf13/viper"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type LoginPage struct {
	IsLoggedIn bool
	UserEmail  string
	Errors     []string
}

type LoginUserData struct {
	ID       int64
	Email    string
	UserType int
}

type OTPValidationData struct {
	OTPRequestID  int64
	OTPUserID     int64
	OTP           string
	OTPValidTill  int64
	OTPSessionID  string
	OTPUserEmail  string
	OTPUserTypeID int
}

func (l *LoginUserData) UserSendOTP(mode string, sessID string, length int, wh *WebHandlers) error {
	wh.Log.Debugf("Generating OTP for session with ID: %s", sessID)
	if mode == "" {
		mode = "email"
	}
	// code to generate and store OTP in DB
	otp, otpErr := wh.u.GenerateAlNumOTP(length)
	if otpErr != nil {
		return otpErr
	}

	err := wh.ReconnectDB()
	if err != nil {
		wh.Log.Errorf("handlers.LoginUserData.SendOTP: error connecting to DB: %s", err.Error())
		return err
	}

	q0SQL := "INSERT INTO otp_requests(user_id,otp_code,delivery_method,session_id) VALUES(?,?,?,?)"

	stmt, stmtErr := wh.db.Prepare(q0SQL)
	if stmtErr != nil {
		wh.Log.Errorf("handlers.LoginUserData.SendOTP: error preparing DB Query for storing OTP in DB: %s", stmtErr.Error())
		return stmtErr
	}

	switch mode {
	case "email":
		// code to actually send OTP using EMAIL

		res, resErr := stmt.Exec(l.ID, otp, mode, sessID)
		if resErr != nil {
			wh.Log.Errorf("handlers.LoginUserData.SendOTP: error executing SQL: %s", q0SQL)
			return resErr
		}

		lastInsertID, liIDErr := res.LastInsertId()
		if liIDErr != nil {
			wh.Log.Errorf("handlers.LoginUserData.SendOTP: error fetching last insert ID: %s", liIDErr.Error())
		}
		emailer := wh.u.NewSMTPEmailWithConfig(viper.GetInt("smtp_port"), viper.GetString("smtp_server"), viper.GetString("smtp_user"), viper.GetString("smtp_pass"))
		otpEmailSendError := emailer.Send(utils.APP_NAME+" Tech", "noreply@am.scalland.com", "LOGIN OTP for "+utils.APP_NAME, fmt.Sprintf("<p>Your OTP to login to https://%s is <b>%s</b></p>", utils.APP_NAME, otp), []string{l.Email}, []string{}, []string{}, []string{})
		if otpEmailSendError != nil {

			wh.Log.Infof("handlers.LoginUserData.SendOTP: LastInsertID: %d", lastInsertID)

			return otpEmailSendError
		}

		return nil
	case "sms":
		// code to actually send OTP using SMS
	case "both":
		// code to send OTP using both email and sms
	case "combined":
		// code to send OTP in combined mode where first half of the OTP is sent over email and the 2nd half over SMS
	default:
		return fmt.Errorf("unknown OTP sending mode %s", mode)
	}

	dbCloseErr := wh.db.Close()
	if dbCloseErr != nil {
		wh.Log.Errorf("handlers.LoginUserData.SendOTP: error closing DB connection for %s", utils.APP_NAME)
	}

	return nil
}

// LoginHandler renders the login page (GET) and processes login (POST).
func (wh *WebHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	lp := LoginPage{
		IsLoggedIn: wh.IsLoggedIn(r, w),
	}

	session, sessErr := wh.GetSession(r)
	if sessErr != nil {
		wh.Log.Errorf("handlers.WebHandlers.LoginHandler: error getting session: %s", sessErr.Error())
		http.Error(w, sessErr.Error(), http.StatusInternalServerError)
		return
	}

	// If already logged in, redirect to dashboard.
	if lp.IsLoggedIn != false {
		wh.Log.Debugf("handlers.LoginHandler: session.Values[is_logged_in]: %t", lp.IsLoggedIn)
		wh.Log.Debugf("handlers.LoginHandler: redirecting to dashboard")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if r.Method == http.MethodGet {
		wh.Log.Debugf("handlers.LoginHandler: requested via GET. Presenting login page")
		// Render the login template.
		data, err := wh.ExecuteTemplate("login", lp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wh.WriteHTML(w, data, http.StatusOK)
		return
	}

	wh.Log.Debugf("handlers.LoginHandler: requested via POST")

	// POST: Process login.
	email := r.FormValue("email")

	// (If you plan to add password support later, retrieve r.FormValue("password") as well.)

	err := wh.ReconnectDB()
	if err != nil {
		wh.Log.Errorf("handlers.LoginUserData.SendOTP: error connecting to DB: %s", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// For this example, we assume that if the email exists, the login is successful.
	var user LoginUserData
	err = wh.db.QueryRow("SELECT user_id, email, user_type_id FROM users WHERE email = ?", email).
		Scan(&user.ID, &user.Email, &user.UserType)

	if err != nil {
		wh.Log.Debugf("handlers.LoginHandler: user specified by email %s does not exist", email)
		wh.Log.Errorf("handlers.LoginHandler: error finding user with email %s: ", err.Error())
		lp.Errors = append(lp.Errors, "Invalid email")
		// Render the login template.
		data, err := wh.ExecuteTemplate("login", lp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wh.WriteHTML(w, data, http.StatusUnauthorized)
		//http.Error(w, "Invalid email", http.StatusUnauthorized)
		return
	}

	wh.Log.Debugf("handlers.LoginHandler: user specified by email %s does exist", email)

	// If execution reaches this point, it means that a valid user exists and has been found
	lp.UserEmail = user.Email // set the lon=gin page data with the email address

	// Read the configured OTP Length for this application
	configOTPLength := viper.GetInt("otp_length")

	wh.Log.Debugf("handlers.LoginHandler: OTP Length: %d", configOTPLength)

	// if the OTP length read from configuration is zero or less, then
	// set the OTP Legth configuration to the default one
	if configOTPLength <= 0 {
		wh.Log.Debugf("handlers.LoginHandler: illegal OTP Length configured: %d. Using default: %d", configOTPLength, utils.DefaultOTPLength)
		configOTPLength = utils.DefaultOTPLength
	}

	wh.Log.Debugf("handlers.LoginHandler: Reading OTP provided by the user")
	otp := r.FormValue("otp")
	if len(otp) == 0 {
		// user has been found but OTP has not been provided by the user
		// so call the code to send the OTP over email and then
		// load the login page with the email address filled-in and an empty OTP field

		wh.Log.Debugf("handlers.LoginHandler: No OTP was provided by the user, sending OTP to their registered email")

		// Generate and send the OTP
		otpSendErr := user.UserSendOTP("email", session.ID, configOTPLength, wh)
		if otpSendErr != nil {
			wh.Log.Debugf("handlers.LoginHandler: OTP sending Error: %s", otpSendErr.Error())
			http.Error(w, otpSendErr.Error(), http.StatusInternalServerError)
			return
		}

		wh.Log.Debugf("handlers.LoginHandler: OTP sent successfully. Rendering Login Page")

		// Render the login template with email address filled
		data, err := wh.ExecuteTemplate("login", lp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wh.WriteHTML(w, data, http.StatusUnauthorized)
		return
	}

	// if the OTP supplied by the user is less than the configured OTP length
	// i.e., the OTP is wrong, load the login template
	if len(otp) < configOTPLength {
		wh.Log.Debugf("handlers.LoginHandler: OTP provided by the user is too short")
		lp.Errors = append(lp.Errors, "Illegal credentials")
		// Render the login template
		data, err := wh.ExecuteTemplate("login", lp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wh.WriteHTML(w, data, http.StatusUnauthorized)
		return
	}

	wh.Log.Debugf("handlers.LoginHandler: validating OTP provided by the user from DB")
	// check the OTP validity and authenticity using DB
	// If the OTP is wrong, load the login template
	// if OTP is correct and valid, set lp.IsLoggedIn = true in session
	var ovd OTPValidationData
	unixTS := time.Now().Unix()
	wh.Log.Debugf("handlers.LoginHandler: SQLEXEC: SELECT otr.otp_request_id, otr.user_id, otr.otp_code, otr.valid_till, otr.session_id, u.email, u.user_type_id FROM otp_requests AS otr LEFT JOIN users AS u ON otr.user_id=u.user_id WHERE u.email = %s AND otp_code = %s AND otr.session_id = %s AND otr.valid_till < %d", email, otp, session.ID, unixTS)
	err = wh.db.QueryRow("SELECT otr.otp_request_id, otr.user_id, otr.otp_code, otr.valid_till, otr.session_id, u.email, u.user_type_id FROM otp_requests AS otr LEFT JOIN users AS u ON otr.user_id=u.user_id WHERE u.email = ? AND otp_code = ? AND otr.session_id = ? AND otr.valid_till < ?", email, otp, session.ID, unixTS).
		Scan(&ovd.OTPRequestID, &ovd.OTPUserID, &ovd.OTP, &ovd.OTPValidTill, &ovd.OTPSessionID, &ovd.OTPUserEmail, &ovd.OTPUserTypeID)

	if err != nil {
		wh.Log.Debugf("handlers.LoginHandler: OTP could not be validated from the DB: %s", err.Error())
		lp.Errors = append(lp.Errors, "Invalid Credentials")
		// Render the login template.
		data, err := wh.ExecuteTemplate("login", lp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		wh.WriteHTML(w, data, http.StatusUnauthorized)
		//http.Error(w, "Invalid email", http.StatusUnauthorized)
		return
	}

	wh.Log.Debugf("handlers.LoginHandler: OTP validated successfully")

	// Save the user info in the session.
	session.Values["user_id"] = ovd.OTPUserID
	session.Values["user_type_id"] = ovd.OTPUserTypeID
	session.Values["is_logged_in"] = true

	lp.IsLoggedIn = true

	wh.Log.Debugf("handlers.LoginHandler: setting session values now")
	wh.Log.Debugf("handlers.LoginHandler: user_id = %d, user_type_id = %d, is_logged_in = %t", session.Values["user_id"], session.Values["user_type_id"], session.Values["is_logged_in"])

	if err = session.Save(r, w); err != nil {
		http.Error(w, fmt.Sprintf("Failed to save session: %s", err.Error()), http.StatusInternalServerError)
		return
	}

	wh.Log.Debugf("handlers.LoginHandler: Session set successfully. Redirecting to home page ---> /")

	// Redirect to dashboard.
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// LogoutHandler clears the session and logs the user out.
func (wh *WebHandlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, err := wh.GetSession(r)
	if err != nil {
		wh.Log.Debugf("handlers.LogoutHandler: error getting session: %s", err.Error())
		http.Error(w, fmt.Sprintf("failed to get session: %s", err.Error()), http.StatusInternalServerError)
	}
	delete(session.Values, "user_id")
	delete(session.Values, "user_type_id")
	delete(session.Values, "is_logged_in")
	wh.isLoggedIn = false
	wh.isAdmin = false

	err = session.Save(r, w)
	if err != nil {
		wh.Log.Debugf("handlers.LogoutHandler: error saving session: %s", err.Error())
		http.Error(w, fmt.Sprintf("failed to save session: %s", err.Error()), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// RequireAuth is a middleware that ensures the user is logged in.
func (wh *WebHandlers) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !wh.IsLoggedIn(r, w) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAdmin is a middleware that ensures the logged-in user is an admin (user_type_id == 1).
func (wh *WebHandlers) RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// check if user is logged-in or not. If they are not, redirect them to login page
		if !wh.IsLoggedIn(r, w) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// If execution reaches here, it means that the user is logged-in. Check if their userTypeID is 1 or not
		if !wh.IsLoggedInAdmin(r, w) {
			http.Error(w, "Unauthorized: Admins only", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
