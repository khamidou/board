package main

import (
    "fmt"
    "time"
    "net/http"
    "strconv"
    "net/url"
    "html/template"
    "encoding/json"
    "os"
)

type Post struct {
    Score int
    Title string
    Contents string
    Replies []Post; // Post []replies;
}

type BoardData struct {
    Posts []Post;
}

var bData *BoardData
var templates = template.Must(template.ParseFiles("templates/newpost.html", "templates/index.html", "templates/post.html"))

func renderTemplate(w http.ResponseWriter, name string, data interface{}) {
    realName := name + ".html";
    err := templates.ExecuteTemplate(w, realName, data)

    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
    if "GET" == r.Method {
        w.Header().Set("Content-Type", "text/html");
        renderTemplate(w, "index", bData);
    }
}

func postHandler(w http.ResponseWriter, r *http.Request) {
    if "GET" == r.Method {
        w.Header().Set("Content-Type", "text/html");
        queryForm, _ := url.ParseQuery(r.URL.RawQuery)
        idparam := queryForm.Get("id")
        if("" == idparam) {
            http.Error(w, "Bad request", http.StatusBadRequest)
            return;
        }
        id, _ := strconv.Atoi(idparam);

        if id < len(bData.Posts) {
            renderTemplate(w, "post", bData.Posts[id]);
        } else {
            http.Error(w, "Post not found", http.StatusNotFound)
        }
    }
}

func newpostHandler(w http.ResponseWriter, r *http.Request) {
    if "GET" == r.Method {

        w.Header().Set("Content-Type", "text/html");
        renderTemplate(w, "newpost", nil);

    } else if "POST" == r.Method {
        p := Post{Score: 1, Title: r.FormValue("title"), Contents: r.FormValue("contents")};
        bData.Posts = append(bData.Posts, p);

        http.Redirect(w, r, "/post?id=" + string(len(bData.Posts)), 303);
        fmt.Fprintf(os.Stdout, "/post?id=" + string(len(bData.Posts)));
    }
}

func saveDb(filename string) {

    for ;; {
        time.Sleep(60000000000)
        b, _ := json.Marshal(bData);
        os.Stdout.Write(b);
    }
}

func main() {
    var boardData *BoardData = &BoardData{};
    bData = boardData;
    go saveDb("stuff");
    http.HandleFunc("/newpost", newpostHandler)
    http.HandleFunc("/post", postHandler)
    http.HandleFunc("/", rootHandler)
    fmt.Println("Listening on http://localhost:8080/");
    http.ListenAndServe(":8080", nil)
}
