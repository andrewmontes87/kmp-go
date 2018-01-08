
package main

import (
    "html/template"
    // "io/ioutil"
    "net/http"
    "regexp"
    // "strings"
    "bytes"
    "os"
    // "html"
    // "fmt"
)

type Page struct {
    Title string
    P string
    T string
    RenderedMatches template.HTML
}

func computePrefix(P string, lenP int) []int {
    s := make([]int, lenP)
    border := 0
    s[0] = 0
    for i := 1; i < lenP; i++ {
        for (border > 0 && P[i] != P[border]) {
            border = s[border - 1]
        }
        if P[i] == P[border] {
            border++
        } else {
            border = 0
        }
        s[i] = border
    }
    return s
}

func findPattern(P string, T string) []int {
    S := P + "%" + T
    lenP := len(P)
    lenT := len(T)
    lenS := lenP + lenT + 1
    sPrefixed := computePrefix(S, lenS)
    var result []int
    for i := lenP + 1; i < lenS; i++ {
        if sPrefixed[i] == lenP {
            result = append(result, i - (2 * lenP))
        }
    }
    return result
}


func renderMatches(M []int, T string, P string) string {
    lenP := len(P)
    lenT := len(T)
    lenM := len(M)
    spentMatches := 0
    end := -1

    var buffer bytes.Buffer

    for i := 0; i < lenT; i++ {
        C := T[i]
        if (spentMatches < lenM && i == M[spentMatches]) {
            if end == -1 {
                buffer.WriteString("<span style=\"color: red; font-weight: bold;\">")
            }
            buffer.WriteString(string(C))
            end = i + lenP
            spentMatches++
        } else if i == end {
            buffer.WriteString("</span>")
            buffer.WriteString(string(C))
            end = -1
        } else {
            buffer.WriteString(string(C))
        }
    }

    return buffer.String()
}


func kmpHandler(w http.ResponseWriter, r *http.Request, title string) {

    pValue := r.FormValue("P")
    tValue := r.FormValue("T")

    if (tValue == "" && pValue == "") {
        pValue = "fox"
        tValue = "The quick brown fox jumped over the lazy dog"
    }

    result := findPattern(pValue, tValue)

    renderedMatchesValue := template.HTML(renderMatches(result, tValue, pValue))

    p := &Page{Title: "KMP", P: pValue, T: tValue, RenderedMatches: renderedMatchesValue}

    renderTemplate(w, "kmp", p)
}

var templates = template.Must(template.ParseFiles("tmpl/kmp.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
    err := templates.ExecuteTemplate(w, tmpl+".html", p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

var validPath = regexp.MustCompile("^/$")

func makeHandler(fn func (http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        m := validPath.FindStringSubmatch(r.URL.Path)
        if m == nil {
            http.NotFound(w, r)
            return
        }
        title := "KMP"
        fn(w, r,title)
    }
}

func main() {
    port := os.Getenv("PORT")

    if port == "" {
        // log.Fatal("$PORT must be set")
        port = "8080"
    }

    http.HandleFunc("/", makeHandler(kmpHandler))
    http.ListenAndServe(":" + port, nil)
}
