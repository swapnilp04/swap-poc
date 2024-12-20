package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strconv"
	"net/http"
	"swapnil-ex/constants"
	"swapnil-ex/models"
	"swapnil-ex/swapErr"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/scrypt"
)

func GetUsers(c echo.Context) error { 
	u := &models.User{}
	users, err := u.All()
	if err != nil {
		fmt.Println("s.ALL(GetUsers)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, users)	
}

func GetUser(c echo.Context) error { 
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	user := &models.User{ID: int(newId)}
	err = user.Find()
	if err != nil {
		fmt.Println("s.Find(GetUser)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, user)
} 

func GetCurrentUser(c echo.Context) error {
	cc := c.(CustomContext)
	if cc.session == nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"error": swapErr.ErrForbidden.Error()})
	}

	session := cc.session

	if ok := session.Valid(); !ok {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"error": swapErr.ErrForbidden.Error()})
	}

	u := &models.User{ID: session.UserID}
	err := u.Find()
	if err != nil {
		fmt.Println("u.Find(GetUser)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, u)	
}


func UpdatePassword(c echo.Context) error {
	cc := c.(CustomContext)
	if cc.session == nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"error": swapErr.ErrForbidden.Error()})
	}

	session := cc.session

	if ok := session.Valid(); !ok {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"error": swapErr.ErrForbidden.Error()})
	}

	user := &models.User{ID: session.UserID}
	err := user.Find()
	if err != nil {
		fmt.Println("u.Find(GetUser)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	
	userData := make(map[string]interface{})
	
	if err := c.Bind(&userData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	if _, ok := userData["password"]; !ok {
		fmt.Println("userData[\"password\"] not present")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	if _, ok := userData["confirm_password"]; !ok {
		fmt.Println("userData[\"confirm_password\"] not present")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		fmt.Println("rand.Read(salt)", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer})
	}

	hash, err := scrypt.Key([]byte(userData["password"].(string)), salt, 32768, 8, 1, 32)
	if err != nil {
		fmt.Println("scrypt.Key()", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer})
	}

	confirmHash, err := scrypt.Key([]byte(userData["confirm_password"].(string)), salt, 32768, 8, 1, 32)
	if err != nil {
		fmt.Println("scrypt.Key()", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer})
	}

	user.Password = hex.EncodeToString(hash)
	user.ConfirmPassword = hex.EncodeToString(confirmHash)
	user.Salt = hex.EncodeToString(salt)

	if err := user.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	err = user.Save()
	if err != nil {
		fmt.Println("user.Save()", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "user Update successfully"})
}

func Register(c echo.Context) error {
	var userData map[string]string
	if err := c.Bind(&userData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	if _, ok := userData["username"]; !ok {
		fmt.Println("userData[\"username\"] not present")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	if _, ok := userData["password"]; !ok {
		fmt.Println("userData[\"password\"] not present")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	if _, ok := userData["confirm_password"]; !ok {
		fmt.Println("userData[\"confirm_password\"] not present")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	if _, ok := userData["role"]; !ok {
		fmt.Println("userData[\"role\"] not present")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		fmt.Println("rand.Read(salt)", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer})
	}

	hash, err := scrypt.Key([]byte(userData["password"]), salt, 32768, 8, 1, 32)
	if err != nil {
		fmt.Println("scrypt.Key()", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer})
	}

	confirmHash, err := scrypt.Key([]byte(userData["confirm_password"]), salt, 32768, 8, 1, 32)
	if err != nil {
		fmt.Println("scrypt.Key()", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer})
	}

	var user models.User
	user.Username = userData["username"]
	user.Password = hex.EncodeToString(hash)
	user.ConfirmPassword = hex.EncodeToString(confirmHash)
	user.Role = userData["role"]
	user.Salt = hex.EncodeToString(salt)

	if err := user.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	err = user.Save()
	if err != nil {
		fmt.Println("user.Save()", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "user created successfully", "user": user})
}

func Login(c echo.Context) error {
	var loginData map[string]string

	if err := c.Bind(&loginData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}
	if _, ok := loginData["username"]; !ok {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrInvalidUser.Error()})
	}

	if _, ok := loginData["password"]; !ok {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrInvalidUser.Error()})
	}
	
	user := &models.User{}
	err := user.FindUserByUsername(loginData["username"])
	if err != nil {
		fmt.Println("username not found", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrInvalidUser.Error()})
	}

	fmt.Println("salt", user.Salt)
	salt, err := hex.DecodeString(user.Salt)
	if err != nil {
		fmt.Println("salt err", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrInvalidUser.Error()})
	}

	hash, err := scrypt.Key([]byte(loginData["password"]), salt, 32768, 8, 1, 32)
	if err != nil {
		fmt.Println("scrypt.Key", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer})
	}

	if user.Password != hex.EncodeToString(hash) {
		fmt.Println("password mismatch")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrInvalidUser.Error()})
	}

	sessionID, err := uuid.NewUUID()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer})
	}

	session := &models.Session{
		ID:        sessionID.String(),
		UserID:    user.ID,
		Role:				user.Role,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(constants.SESSION_EXPIRY * time.Hour),
	}

	err = session.Save()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer})
	}

	// we can also set the session ID as token back
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "user loggedin successfully", "token": session.ID, "username": user.Username, "role": user.Role})
}

func Logout(c echo.Context) error {
	cc := c.(CustomContext)
	if cc.session == nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"error": swapErr.ErrForbidden.Error()})
	}

	session := cc.session

	if ok := session.Valid(); !ok {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"error": swapErr.ErrForbidden.Error()})
	}

	if err := session.Destroy(); err != nil {
		fmt.Println("session.Destroy()", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "user logged out"})

}

func UpdateUser(c echo.Context) error {
	cc := c.(CustomContext)
	if cc.session == nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"error": swapErr.ErrForbidden.Error()})
	}

	session := cc.session

	if ok := session.Valid(); !ok {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"error": swapErr.ErrForbidden.Error()})
	}

	var updateData map[string]string

	if err := c.Bind(&updateData); err != nil {
		fmt.Println("c.Bind()", err)
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	if _, ok := updateData["password"]; !ok {
		fmt.Println("updateData[\"password\"] not present")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	if _, ok := updateData["confirm_password"]; !ok {
		fmt.Println("updateData[\"confirm_password\"] not present")
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": swapErr.ErrBadData.Error()})
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		fmt.Println("rand.Read(salt)", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	hash, err := scrypt.Key([]byte(updateData["password"]), salt, 32768, 8, 1, 32)
	if err != nil {
		fmt.Println("scrypt.Key()", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	confirmHash, err := scrypt.Key([]byte(updateData["confirm_password"]), salt, 32768, 8, 1, 32)
	if err != nil {
		fmt.Println("scrypt.Key()", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}

	user := &models.User{}
	user.ID = session.UserID
	if err := user.Load(); err != nil {
		fmt.Println("user.Load()", err)
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"error": swapErr.ErrInternalServer.Error()})
	}
	user.Password = hex.EncodeToString(hash)
	user.ConfirmPassword = hex.EncodeToString(confirmHash)
	user.Salt = hex.EncodeToString(salt)

	if err := user.Validate(); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": err.Error()})
	}

	user.Save()
	user.Password = ""
	user.ConfirmPassword = ""
	user.Salt = ""
	return c.JSON(http.StatusOK, map[string]interface{}{"message": "user update successfully", "user": user})
}


func DeactivateUser(c echo.Context) error { 
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	user := &models.User{ID: int(newId)}
	err = user.Find()
	if err != nil {
		fmt.Println("s.Find(GetUser)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	if user.Role == "Admin" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Can not Deactive Admin User"})
	}
	err = user.DeactiveUser()
	if err != nil {
		fmt.Println("s.Find(GetUser)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, user)
} 

func ActivateUser(c echo.Context) error { 
	id := c.Param("id")
	newId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("strconv.Atoi failed", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"message": swapErr.ErrBadData.Error()})
	}

	user := &models.User{ID: int(newId)}
	err = user.Find()
	if err != nil {
		fmt.Println("s.Find(GetUser)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}
	if user.Role == "Admin" {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Can not Deactive Admin User"})
	}
	err = user.ActiveUser()
	if err != nil {
		fmt.Println("s.Find(GetUser)", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": swapErr.ErrInternalServer.Error()})
	}

	return c.JSON(http.StatusOK, user)
} 

// package main

// import (
// 	"database/sql"
// 	"encoding/hex"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"time"

// 	"github.com/gorilla/mux"
// 	"github.com/gorilla/sessions"
// 	_ "github.com/mattn/go-sqlite3"
// 	"github.com/satori/go.uuid"
// 	"golang.org/x/crypto/scrypt"
// )

// type User struct {
// 	ID       int    `json:"id"`
// 	Username string `json:"username"`
// 	Password string `json:"-"`
// 	Salt     string `json:"-"`
// }

// type Session struct {
// 	ID        string
// 	UserID    int
// 	CreatedAt time.Time
// 	ExpiresAt time.Time
// }

// var db *sql.DB
// var store *SqliteStore

// func main() {
// 	var err error
// 	db, err = sql.Open("sqlite3", "./users.db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	store = NewSqliteStore(db, "sessions", []byte("secret-key"))

// 	router := mux.NewRouter()

// 	router.HandleFunc("/users", createUser).Methods("POST")
// 	router.HandleFunc("/users/{id}", getUser).Methods("GET")
// 	router.HandleFunc("/users/{id}", updateUser).Methods("PUT")
// 	router.HandleFunc("/users/{id}", deleteUser).Methods("DELETE")
// 	router.HandleFunc("/login", login).Methods("POST")
// 	router.HandleFunc("/logout", logout).Methods("POST")

// 	log.Fatal(http.ListenAndServe(":8000", router))
// }

// func createUser(w http.ResponseWriter, r *http.Request) {
// 	// ...
// }

// func getUser(w http.ResponseWriter, r *http.Request) {
// 	// ...
// }

// func updateUser(w http.ResponseWriter, r *http.Request) {
// 	// ...
// }

// func deleteUser(w http.ResponseWriter, r *http.Request) {
// 	// ...
// }

// func login(w http.ResponseWriter, r *http.Request) {
// 	var user User
// 	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	row := db.QueryRow("SELECT * FROM users WHERE username=?", user.Username)
// 	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Salt)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
// 			return
// 		}
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	salt, err := hex.DecodeString(user.Salt)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	hash, err := scrypt.Key([]byte(user.Password), salt, 32768, 8, 1, 32)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if hex.EncodeToString(hash) != user.Password {
// 		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
// 		return
// 	}

// 	sessionID, err := uuid.NewV4()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	session := Session{
// 		ID:        sessionID.String(),
// 		UserID:    user.ID,
// 		CreatedAt: time.Now(),
// 		ExpiresAt: time.Now().Add(24 * time.Hour),
// 	}

// 	err = store.Save(&session)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "session_id",
// 		Value:    session.ID,
// 		Expires:  session.Ex
