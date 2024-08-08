package data 


type Movie struct {
ID int64
CreatedAt time.Time
Title string
Year int32
Runtime int32
Genres []string
Version int32
}

// MarshalJSON() method satisfies the
// json.Marshaler interface.
func (m Movie) MarshalJSON() ([]byte, error) {
var runtime string

if m.Runtime != 0 {
runtime = fmt.Sprintf("%d mins", m.Runtime)
}

aux := struct {
ID int64 `json:"id"`
Title string `json:"title"`
Year int32 `json:"year,omitempty"`
Runtime string `json:"runtime,omitempty"` // This is a string.
Genres []string `json:"genres,omitempty"`
Version int32 `json:"version"`
}{

ID: m.ID,
Title: m.Title,
Year: m.Year,
Runtime: runtime,
Genres: m.Genres,
Version: m.Version,
}

// Encode the anonymous struct to JSON, 
return json.Marshal(aux)
}