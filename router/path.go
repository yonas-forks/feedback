package router

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/andrewarrow/feedback/models"
)

type LayoutVars struct {
	Title   string
	User    *models.User
	Phone   string
	Content template.HTML
	Flash   string
}

func (r *Router) PlaceContentInLayoutVars(flash string, user *models.User, filename string, vars any) *LayoutVars {
	content := new(bytes.Buffer)
	r.Template.ExecuteTemplate(content, filename, vars)

	lvars := LayoutVars{}
	lvars.Title = "Feedback"
	lvars.Phone = r.Site.Phone
	lvars.User = user
	lvars.Flash = flash
	lvars.Content = template.HTML(content.String())
	return &lvars
}

func (r *Router) SendContentInLayout(flash string, user *models.User, writer http.ResponseWriter,
	filename string, contentVars any, status int) {
	vars := r.PlaceContentInLayoutVars(flash, user, filename, contentVars)
	writer.WriteHeader(status)
	r.Template.ExecuteTemplate(writer, "application_layout.html", vars)
}

func (r *Router) RouteFromRequest(writer http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	cookie, err := request.Cookie("user")
	var user *models.User
	if err == nil && cookie.Value != "" {
		user = r.LookupUser(cookie.Value)
	}
	cookie, err = request.Cookie("flash")
	flash := ""
	if err == nil && cookie.Value != "" {
		flash = cookie.Value
		DestroyFlash(writer)
	}

	if path == "/" {
		r.SendContentInLayout(flash, user, writer, "welcome.html",
			WelcomeIndexVars(r.Db), 200)
	} else if strings.HasPrefix(path, "/assets") {
		r.HandleAsset(path, writer)
	} else if !strings.HasSuffix(path, "/") {
		http.Redirect(writer, request, fmt.Sprintf("%s/", path), 301)
	} else {
		c := Context{}
		c.writer = writer
		c.request = request
		c.flash = flash
		c.method = request.Method
		c.router = r
		c.user = user
		c.path = path
		c.db = r.Db
		c.tokens = strings.Split(path, "/")
		c.userRequired = r.IsUserRequired(path, c.method)
		if c.userRequired && c.user == nil {
			http.Redirect(c.writer, c.request, "/sessions/new/", 302)
			return
		}
		if c.method == "POST" {
			c.ReadFormPost()
		}
		handleContext(&c)
		if c.notFound {
			r.SendContentInLayout("", user, writer, "404.html", nil, 404)
		}
	}
}
