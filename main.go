package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type NpmjsResult struct {
	Total   int `json:"total"`
	Results []struct {
		Package struct {
			Name        string    `json:"name"`
			Scope       string    `json:"scope"`
			Version     string    `json:"version"`
			Description string    `json:"description"`
			Keywords    []string  `json:"keywords"`
			Date        time.Time `json:"date"`
			Links       struct {
				Npm        string `json:"npm"`
				Homepage   string `json:"homepage"`
				Repository string `json:"repository"`
				Bugs       string `json:"bugs"`
			} `json:"links"`
			Author struct {
				Name string `json:"name"`
			} `json:"author"`
			Publisher struct {
				Username string `json:"username"`
				Email    string `json:"email"`
			} `json:"publisher"`
			Maintainers []struct {
				Username string `json:"username"`
				Email    string `json:"email"`
			} `json:"maintainers"`
		} `json:"package"`
		Score struct {
			Final  float64 `json:"final"`
			Detail struct {
				Quality     float64 `json:"quality"`
				Popularity  float64 `json:"popularity"`
				Maintenance int     `json:"maintenance"`
			} `json:"detail"`
		} `json:"score"`
		SearchScore float64 `json:"searchScore"`
	} `json:"results"`
}

func main() {
	arg := os.Args[1]
	if len(os.Args) < 2 {
		fmt.Println("Usage: eisen new PROJECT_NAME")
		os.Exit(-1)
	}

	if arg == "new" {
		projectName := os.Args[2]
		if _, err := os.Stat(projectName); os.IsNotExist(err) {
			err = os.Mkdir(projectName, os.ModePerm)
			check(err)

			resp, err := http.Get("https://api.npms.io/v2/search?q=@kloudsoftware/eisen")
			body, err := ioutil.ReadAll(resp.Body)
			var npmRes NpmjsResult

			_ = json.Unmarshal(body, &npmRes)

			err = os.Chdir(projectName)
			packagejson := []byte("{\n  \"name\": \""+ projectName + "\",\n  \"version\": \"1.0.0\",\n  \"description\": \"An eisen Project\",\n  \"main\": \"index.ts\",\n  \"scripts\": {\n    \"dev\": \"parcel src/index.html\",\n    \"clean\": \"rm -rf dist/ && rm -rf lib/\" },\n  \"author\": \"\",\n  \"license\": \"ISC\",\n  \"dependencies\": {\n    \"@kloudsoftware/eisen\": \"" + npmRes.Results[0].Package.Version + "\"\n  },\n  \"devDependencies\": {\n    \"parcel-plugin-static-files-copy\": \"^2.0.0\",\n    \"sass\": \"^1.18.0\",\n    \"typescript\": \"^3.4.2\"\n  }\n}")
			gitignore := []byte("node_modules/\nlib/\ndist\n.cache/\n.#*")

			err = ioutil.WriteFile("package.json", packagejson, 0644)
			check(err)
			err = ioutil.WriteFile(".gitignore", packagejson, 0644)
			check(err)
			err = os.Mkdir("src", os.ModePerm)
			err = os.Chdir("src")

			check(err)
			indexhtml := []byte("<!DOCTYPE html>\n<head>\n<title>" + projectName + "</title>\n<link rel=\"stylesheet\" href=\"style.scss\">\n</head>\n<body>\n<div id=\"target\">\n</div>\n</body>\n<script src=\"index.ts\"></script>")
			stylescss := []byte("@import 'reset.scss';\n@mixin transform($property) {\n    -webkit-transform: $property;\n    -ms-transform: $property;\n    -moz-transform: $property;\n    transform: $property;\n  }\nbody {\n  font: 100% Helvetica, sans-serif;\n  background-color: #efefef;\n}")
			resetscss := []byte("html, body, div, span, applet, object, iframe, p, blockquote, pre, a, abbr, acronym, address, big, cite, code, del, dfn, em, img, ins, kbd, q, s, samp, small, strike, strong, sub, sup, tt, var, b, u, i, center, dl, dt, dd, ol, ul, li, fieldset, form, label, legend, table, caption, tbody, tfoot, thead, tr, th, td, article, aside, canvas, details, embed, figure, figcaption, footer, header, hgroup, menu, nav, output, ruby, section, summary, time, mark, audio, video {\n    margin: 0;\n    padding: 0;\n    border: 0;\n    font-size: 100%;\n    font: inherit;\n    vertical-align: baseline;\n}\narticle, aside, details, figcaption, figure, footer, header, hgroup, menu, nav, section {\n    display: block;\n}\nbody {\n    line-height: 1;\n}\nol, ul, li {\n    list-style: none;\n}\nblockquote, q {\n    quotes: none;\n}\nblockquote {\n    &:before, &:after {\n    content: '';\n    content: none;\n    }\n}\nq {\n    &:before, &:after {\n    content: '';\n    content: none;\n    }\n}\ntable {\n    border-collapse: collapse;\n    border-spacing: 0;\n}\nbutton {\n    background-color: transparent;\n    outline: none;\n    border: 0;\n    cursor: pointer;\n}")
			indexts := []byte("import { VApp, Renderer } from '@kloudsoftware/eisen';\nconst renderer = new Renderer();\nconst app = new VApp(\"target\", renderer);\napp.init();\napp.createElement(\"h1\", \"Hello eisen!\", app.rootNode);")
			err = ioutil.WriteFile("index.html", indexhtml, 0644)
			check(err)
			err = ioutil.WriteFile("style.scss", stylescss, 0644)
			check(err)
			err = ioutil.WriteFile("_reset.scss", resetscss, 0644)
			check(err)
			err = ioutil.WriteFile("index.ts", indexts, 0644)
			check(err)

		}
	}
}
