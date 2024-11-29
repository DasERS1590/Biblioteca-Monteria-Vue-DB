package main 

import (
    "net/http"
    "encoding/json"
)


func (app *application) loginHandler(w http.ResponseWriter , r *http.Request){
    var loginData struct {
        Username string `json:"gmail"`
        Password string `json:"password"`
    }
    err := json.NewDecoder(r.Body).Decode(&loginData)
    if err != nil {
        http.Error(w, "Datos inválidos", http.StatusBadRequest)
        return
    }

    userAuth := &models.UserAuth{}
    token, err := userAuth.Login(app.config.db , loginData.Username, loginData.Password)
    
    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    json.NewEncoder(w).Encode(map[string]string{
        "token": token,
        "role":  userAuth.Role,
    })
}

func (app *application) registerUserHandler(w http.ResponseWriter , r *http.Request){

}

func (app *application) updateUserHandler(w http.ResponseWriter , r *http.Request){

}

func (app *application) getUserHandler(w http.ResponseWriter , r *http.Request){

}
