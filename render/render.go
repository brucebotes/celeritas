package render

import (
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
)

type Render struct {
	Renderer   string
	RootPath   string
	Secure     bool
	Port       string
	ServerName string
	JetViews   *jet.Set
	Session    *scs.SessionManager
}

type TemplateData struct {
	IsAuthenticated bool
	IntMap          map[string]int
	StringMap       map[string]string
	FloatMap        map[string]float32
	Data            map[string]interface{}
	CSRFToken       string
	Port            string
	ServerName      string
	Secure          bool
	Error           string
	Flash           string
}

func (c *Render) defaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.Secure = c.Secure
	td.ServerName = c.ServerName
	td.Port = c.Port
	td.CSRFToken = nosurf.Token(r)
	if c.Session.Exists(r.Context(), "userID") {
		td.IsAuthenticated = true
	}
	// Get the Error and Flash messages from the session
	td.Error = c.Session.PopString(r.Context(), "error")
	td.Flash = c.Session.PopString(r.Context(), "flash")

	return td
}

func (c *Render) Page(w http.ResponseWriter, r *http.Request, view string, variables, data interface{}) error {
	switch strings.ToLower(c.Renderer) {
	case "go":
		return c.GoPage(w, r, view, data)
	case "jet":
		return c.JetPage(w, r, view, variables, data)
	default:

	}
	return errors.New("no rendering engine specified")
}

// GoPage renders a template using the Go templating engine
// Note: this function cannot reference other templates (or fragments) 
//       relative to it self( as is the case of jet templates)
//       Thus it is most suited to SPA which do not require
//       multiple pages with the same header or footers templates
//       that are re-used.
//      ( ie header/footer components templates that are shared )
//
//      The generic data is passed via the TemplateData struct. 
//      Which is available in the Celeritas and Render structs
//      as pointers fields.
func (c *Render) GoPage(w http.ResponseWriter, r *http.Request, view string, data interface{}) error {
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", c.RootPath, view))
	if err != nil {
		return err
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}
	td = c.defaultData(td, r)

	err = tmpl.Execute(w, td)
	if err != nil {
		return err
	}

	return nil
}

// JetPage renders a template using the Jet templating engine
// Note: The jet templates can receive their data to be rendered
//      via the variables and/or data structuers. 
//      - The variables structure is the *jet.VarMaps struct defined 
//        by jet.
//      - The data structure uses the generic data is passed via the 
//        TemplateData struct. 
//        Which is available in the Celeritas and Render structs
//        as pointers fields.
func (c *Render) JetPage(w http.ResponseWriter, r *http.Request, templateName string, variables, data interface{}) error {
	var vars jet.VarMap

	if variables == nil {
		vars = make(jet.VarMap)
	} else {
		vars = variables.(jet.VarMap)
	}

	td := &TemplateData{}
	if data != nil {
		td = data.(*TemplateData)
	}

	td = c.defaultData(td, r)

	t, err := c.JetViews.GetTemplate(fmt.Sprintf("%s.jet", templateName))
	if err != nil {
		log.Println(err)
		return err
	}

	if err = t.Execute(w, vars, td); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
