package main

import (
	"encoding/json"
"errors"
"net/http"
"strconv"
"github.com/julienschmidt/httprouter"
)

func (app *application) handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
"status": "available",
"environment": app.config.env,
"version": version,
}

err := app.writeJSON(w, http.StatusOK, data, nil)

if err != nil {
app.logger.Println(err)
http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
}
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {

js, err := json.Marshal(data)
if err != nil {
return err
}


js = append(js, '\n')

for key, value := range headers {
w.Header()[key] = value
}

w.Header().Set("Content-Type", "application/json")
w.WriteHeader(status)
w.Write(js)

return nil
}
