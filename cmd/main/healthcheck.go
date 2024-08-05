package main

import (
	"fmt"
	"net/http"
)

func (app *application) handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
"status": "available",
"environment": app.config.env,
"version": version,
}
// Pass the map to the json.Marshal()
js, err := json.Marshal(data)
if err != nil {
app.logger.Println(err)
http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
return
}

js = append(js, '\n')
// encoding the data
w.Header().Set("Content-Type", "application/json")
// Use w.Write() to send the []byte slice containing the JSON as the response body.
w.Write(js)
}
