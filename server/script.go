package server

import (
	"net/http"
	"text/template"
)

func (s *Server) scriptHandler(w http.ResponseWriter, r *http.Request) {
	base := r.URL.Query().Get("base")
	res, err := s.createResponse(base)
	if err != nil {
		http.Error(w, "Error creating response", http.StatusInternalServerError)
		return
	}

	t, err := template.New("script").Parse(scriptTemplate)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/javascript")
	err = t.Execute(w, res)
	if err != nil {
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
	}
}

var scriptTemplate = `
(function() {
	var currency = {
		currencies: [{{ range .Rates }}'{{.Name}}',{{end}}],

		rates: {
			{{ range .Rates }}
			'{{.Name}}': {{.Rate}},
			{{ end }}
		},

		convert: function(amount, target) {
			return this.rates[target] * amount
		}
	}

	window.currency = currency
})();
`
