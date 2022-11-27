package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"text/template"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Page struct {
	Title string
	Body  string
}

var templates = template.Must(template.ParseFiles("./templates/edit.html", "./templates/view.html"))
var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

func (p *Page) save(coll *mongo.Collection) {
	opts := options.Update().SetUpsert(true)
	doc := bson.D{{"title", p.Title}, {"content", p.Body}}
	update := bson.D{{"$set", doc}}
	_, err := coll.UpdateOne(context.TODO(), bson.D{{"title", p.Title}}, update, opts)
	if err != nil {
		return
	}
}

func loadPage(title string, coll *mongo.Collection) *Page {
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{"title", title}}).Decode(&result)

	if err != nil {
		return nil
	}
	fmt.Printf("found document %v", result)
	return &Page{Title: result["title"].(string), Body: result["content"].(string)}
}

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invalid Page Title")
	}
	return m[2], nil // The title is the second subexpression.
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string, *mongo.Collection), coll *mongo.Collection) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2], coll)
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string, coll *mongo.Collection) {
	p := loadPage(title, coll)
	if p == nil {
		http.Redirect(w, r, "/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string, coll *mongo.Collection) {
	p := loadPage(title, coll)
	if p == nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string, coll *mongo.Collection) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: body}
	p.save(coll)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	http.Redirect(w, r, "/view/"+title, http.StatusFound)
}

func main() {
	fmt.Println("Works!")

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_CONNECTION_STRING")))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(ctx)

	coll := client.Database(os.Getenv("MONGO_DB_NAME")).Collection(os.Getenv("MONGO_COLLECTION_NAME"))

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/view/", makeHandler(viewHandler, coll))
	http.HandleFunc("/edit/", makeHandler(editHandler, coll))
	http.HandleFunc("/save/", makeHandler(saveHandler, coll))
	log.Fatal(http.ListenAndServe(":8085", nil))
}
