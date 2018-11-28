package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"

	"KBTGCourse/project/pq/post"
)

type Display struct {
	data []post.Event
}

func indexPageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	indexTmpl, err := template.ParseFiles("html/layout.html", "html/index.html", "html/navbar.html", "html/footer.html")

	err = indexTmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		http.Error(w, "blog: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func eventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	registerTmpl, err := template.ParseFiles("html/layout.html", "html/eventManagement.html", "html/navbar.html", "html/footer.html")

	data, err := post.All()
	if err != nil {
		http.Error(w, "[eventHandler] select all data got error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	display := map[string]interface{}{
		"events": data,
	}

	err = registerTmpl.ExecuteTemplate(w, "layout", display)
	if err != nil {
		http.Error(w, "[eventHandler] load html page got error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func showAddHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	registerTmpl, err := template.ParseFiles("html/layout.html", "html/addEvent.html", "html/navbar.html", "html/footer.html")

	err = registerTmpl.ExecuteTemplate(w, "layout", nil)
	if err != nil {
		http.Error(w, "blog: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func addEventDataToDBHandler(w http.ResponseWriter, r *http.Request) {
	var evnt post.Event

	evnt.Name = r.FormValue("name")
	evnt.Place = r.FormValue("place")
	evnt.Speaker = r.FormValue("speaker")
	evnt.Detail = r.FormValue("detail")

	err := post.Insert(&evnt)
	if err != nil {
		fmt.Fprintf(w, "[Insert] got error: %s", err)
		return
	}

	evnt.EventId = evnt.Id

	data := r.PostForm
	for i := 0; i < len(data["limit"]); i++ {
		evnt.StartDate = data["start_date"][i]
		evnt.EndDate = data["end_date"][i]
		evnt.StartTime = data["start_time"][i]
		evnt.EndTime = data["end_time"][i]
		evnt.Limit = data["limit"][i]
		post.InsertEventDetails(&evnt)
	}
}

func DeleteEventFromDbHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/eventmanagement/delete/"))
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	err = post.Delete(id)
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	http.Redirect(w, r, "/eventmanagement/", http.StatusSeeOther)
}

func showUpdateHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/eventmanagement/update/"))
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}
	// fmt.Println(id)
	rs, err := post.FindEventByID(id)
	if err != nil {
		fmt.Fprintf(w, "[showUpdateHandler] FindByID got error: %s", err)
		return
	}

	display := map[string]interface{}{
		"id":      rs.Id,
		"name":    rs.Name,
		"place":   rs.Place,
		"speaker": rs.Speaker,
		"detail":  rs.Detail,
		"details": rs.Detail,
	}

	registerTmpl, err := template.ParseFiles("html/layout.html", "html/updateEvent.html", "html/navbar.html", "html/footer.html")

	err = registerTmpl.ExecuteTemplate(w, "layout", display)
	if err != nil {
		http.Error(w, "blog: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func getDataToShowHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/eventmanagement/getDataToShow/"))
	if err != nil {
		fmt.Fprintf(w, "%v", err)
		return
	}

	fmt.Println(id)
}

func eventsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	registerTmpl, err := template.ParseFiles("html/layout.html", "html/event.html", "html/navbar.html", "html/footer.html")

	data, err := post.All()
	if err != nil {
		http.Error(w, "[eventHandler] select all data got error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	display := map[string]interface{}{
		"events": data,
	}

	err = registerTmpl.ExecuteTemplate(w, "layout", display)
	if err != nil {
		http.Error(w, "[eventHandler] load html page got error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	id, err := strconv.Atoi(strings.TrimPrefix(r.URL.Path, "/register/"))
	if err != nil {
		fmt.Fprintf(w, "[registerHandler] got error: %v", err)
		return
	}

	rs, err := post.FindDetailByID(id)
	if err != nil {
		http.Error(w, "[registerHandler] FindDetailByID got error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	display := map[string]interface{}{
		"main":   rs.Main,
		"option": rs.Detail,
	}
	funcs := template.FuncMap{"add": add, "showOption": showOption}
	registerTmpl := template.Must(template.New("foo").Funcs(funcs).ParseFiles("html/layout.html", "html/register.html", "html/navbar.html", "html/footer.html"))

	err = registerTmpl.ExecuteTemplate(w, "layout", display)
	if err != nil {
		http.Error(w, "[registerHandler] execute got error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func addRegisterToDbHandler(w http.ResponseWriter, r *http.Request) {
	var reg post.Register

	reg.FName = r.FormValue("fName")
	reg.LName = r.FormValue("lName")
	reg.UserId = r.FormValue("empId")
	reg.Tel = r.FormValue("tel")
	reg.Event = r.FormValue("eventSelected")

	err := post.InsertRegister(&reg)
	if err != nil {
		fmt.Fprintf(w, "[addRegisterToDbHandler] got error: %s", err)
		return
	}
}

func add(x, y int) int {
	return x + y
}

func showOption(count, limit int) bool {
	if count < limit {
		return true
	}
	return false
}

func startServer() error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.Method, r.URL.Path)
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/":
			indexPageHandler(w, r)
		case r.Method == http.MethodGet && r.URL.Path == "/eventmanagement/":
			eventHandler(w, r)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/eventmanagement/") && strings.HasSuffix(r.URL.Path, "/add/"):
			showAddHandler(w, r)
		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/eventmanagement/") && strings.HasSuffix(r.URL.Path, "/add/toDB/"):
			addEventDataToDBHandler(w, r)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/eventmanagement/delete/"):
			DeleteEventFromDbHandler(w, r)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/eventmanagement/update/"):
			showUpdateHandler(w, r)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/eventmanagement/getDataToShow/"):
			getDataToShowHandler(w, r)
		case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/register/"):
			registerHandler(w, r)
		case r.Method == http.MethodGet && r.URL.Path == "/events/":
			eventsHandler(w, r)
		case r.Method == http.MethodPost && strings.HasPrefix(r.URL.Path, "/register/add/toDB/"):
			addRegisterToDbHandler(w, r)
		}
	})

	http.Handle("/css/", http.FileServer(http.Dir("./html/assets/")))
	http.Handle("/fonts/", http.FileServer(http.Dir("./html/assets/")))
	http.Handle("/img/", http.FileServer(http.Dir("./html/assets/")))
	http.Handle("/js/", http.FileServer(http.Dir("./html/assets/")))
	http.Handle("/scss/", http.FileServer(http.Dir("./html/assets/")))
	http.Handle("/vendor/", http.FileServer(http.Dir("./html/assets/")))

	return http.ListenAndServe(":8000", nil)
}
