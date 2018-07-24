package main

import (
	// "fmt"
	"github.com/gorilla/mux"
	"gopkg.in/russross/blackfriday.v2"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	t = template.Must(template.ParseGlob("views/*"))
)

type Post struct {
	Title   string
	Date    string
	Summary string
	Body    string
	File    string
	Cover   string
	Cate    string
}
type Home struct {
	Posts []Post
	Cates []string
}

func getPosts(cateid string) ([]Post, []string) {

	a := []Post{}
	cates := []string{}
	files, _ := filepath.Glob("posts/*")

	for _, f := range files {

		filename := strings.Replace(f, "posts/", "", -1)

		if filename == ".DS_Store" {
			continue
		}
		post := getPost(filename)

		if cateid == "" || cateid == post.Cate {
			a = append(a, post)
		}

		cates = append(cates, post.Cate)
	}
	return a, cates
}
func getPost(filename string) Post {

	filepath := "posts/" + filename + "/main.md"

	fileread, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(fileread), "\n")

	title := string(lines[1])
	title = strings.Replace(title, "标题:", "", -1)

	date := string(lines[2])
	date = strings.Replace(date, "日期:", "", -1)

	summary := string(lines[3])
	summary = strings.Replace(summary, "简介:", "", -1)

	cover := string(lines[4])
	cover = strings.Replace(cover, "封面:", "", -1)
	cover = replaceImagePath(cover, filename)

	cate := string(lines[5])
	cate = strings.Replace(cate, "分类:", "", -1)

	body := strings.Join(lines[6:len(lines)], "\n")
	body = replaceImagePath(body, filename)

	body = string(blackfriday.Run([]byte(body), blackfriday.WithNoExtensions()))

	post := Post{title, date, summary, body, filename, cover, cate}

	return post
}
func PostHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	postid := vars["postid"]
	post := getPost(postid)
	data := template.HTML(post.Body)
	renderTemplate(w, "detail.html", data)
}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	cateid, ok := vars["cateid"]
	if !ok {
		cateid = ""
	}
	posts, cates := getPosts(cateid)
	home := Home{posts, cates}

	renderTemplate(w, "index.html", home)
}
func ImageHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	imageid := vars["imgid"]
	file := vars["file"]

	imgpath := "/Users/ck/go/src/jiangjn/menu3/posts/" + file + "/" + imageid

	_, err := os.Stat(imgpath)
	if err != nil {
		w.Write([]byte("Error:Image Not Found."))
		return
	}
	http.ServeFile(w, r, imgpath)
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", HomeHandler).Methods("GET")
	r.HandleFunc("/cates/{cateid}", HomeHandler).Methods("GET")
	r.HandleFunc("/posts/{postid}", PostHandler).Methods("GET")
	r.HandleFunc("/images/{file}/{imgid}", ImageHandler).Methods("GET")
	http.ListenAndServe(":8000", r)
}
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	err := t.ExecuteTemplate(w, tmpl, data)
	if err != nil {
		http.Error(w, "error 500:"+" "+err.Error(), http.StatusInternalServerError)
	}
}

func replaceImagePath(text, file string) string {
	text = strings.Replace(text, "images:", "/images/"+file+"/", -1)
	return text
}
