package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
	if len(os.Args) < 2 {
		fmt.Println("Usage: - eisen new PROJECT_NAME")
		fmt.Println("       - eisen new component COMPONENT_NAME")
		os.Exit(-1)
	}

	arg := os.Args[1]

	if arg == "new" {
		projectName := os.Args[2]

		if projectName == "component" {
			if len(os.Args) < 3 {
				fmt.Println("Usage: eisen new component COMPONENT_NAME")
				os.Exit(-1)
			}

			if _, err := os.Stat("package.json"); os.IsNotExist(err) {
				fmt.Println("No package.json file present, not an eisen repository or root")
				os.Exit(-1)
			}

			componentName := os.Args[3]
			if _, err := os.Stat("components"); os.IsNotExist(err) {
				err = os.Mkdir("components", os.ModePerm)
				check(err)
			}

			err := os.Chdir("components")
			check(err)

			if _, err := os.Stat(strings.ToLower(componentName)); os.IsNotExist(err) {
				err = os.Mkdir(strings.ToLower(componentName), os.ModePerm)
				check(err)
			}

			err = os.Chdir(strings.ToLower(componentName))
			check(err)

			component := []byte(`import { Component, VApp, ComponentBuildFunc, Props, VNode, src, cssClass } from "@kloudsoftware/eisen"

export class ` + componentName + ` extends Component {
    build(app: VApp): ComponentBuildFunc {
        return (root: VNode, props: Props) => {
            return {
                mounted: () => {
                }
            };
        }
    }
}`)

			err = ioutil.WriteFile(componentName+".ts", component, 0644)

			os.Exit(0)

		}

		if _, err := os.Stat(projectName); os.IsNotExist(err) {
			err = os.Mkdir(projectName, os.ModePerm)
			check(err)

			resp, err := http.Get("https://api.npms.io/v2/search?q=@kloudsoftware/eisen")
			body, err := ioutil.ReadAll(resp.Body)
			var npmRes NpmjsResult

			_ = json.Unmarshal(body, &npmRes)

			err = os.Chdir(projectName)
			packagejson := []byte("{\n  \"name\": \"" + projectName + "\",\n  \"version\": \"1.0.0\",\n  \"description\": \"An eisen Project\",\n  \"main\": \"index.ts\",\n  \"scripts\": {\n    \"dev\": \"parcel src/index.html\",\n    \"clean\": \"rm -rf dist/ && rm -rf lib/\" },\n  \"author\": \"\",\n  \"license\": \"ISC\",\n  \"dependencies\": {\n    \"@kloudsoftware/eisen\": \"" + npmRes.Results[0].Package.Version + "\"\n  },\n  \"devDependencies\": {\n    \"parcel-plugin-static-files-copy\": \"^2.0.0\",\n    \"sass\": \"^1.18.0\",\n    \"typescript\": \"^3.4.2\"\n  }\n}")
			gitignore := []byte("node_modules/\nlib/\ndist\n.cache/\n.#*")

			err = ioutil.WriteFile("package.json", packagejson, 0644)
			check(err)
			err = ioutil.WriteFile(".gitignore", gitignore, 0644)
			check(err)
			err = os.Mkdir("src", os.ModePerm)
			err = os.Chdir("src")

			check(err)
			indexhtml := []byte(`<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <link href="./style.scss" rel="stylesheet"/>
</head>

<body>
    <div id="target"></div>
</body>

<script src="./index.ts"></script>`)
			stylescss := []byte(`@import 'reset.scss';

@mixin transform($property) {
    -webkit-transform: $property;
    -ms-transform: $property;
    -moz-transform: $property;
    transform: $property;
}

$transparent: transparent;

$white: #FFFFFF;
$grey: #B8C2CC;
$grey-lightest: #f8f8f8;
$grey-lighter: #efefef;
$grey-light: #dadada;
$grey-dark: #8795A1;
$grey-darker: #949ea8;
$grey-darkest: #3D4852;
$black: #22292F;

$blue: #3490DC;
$blue-lightest: #EFF8FF;
$blue-lighter: #ccf0ff;
$blue-light: #3490dc;
$blue-dark: #2779BD;
$blue-darker: #266fb3;
$blue-darkest: #12283A;

$red: #B94848;
$red-lightest: #fcebea;
$red-lighter: #f9acaa;
$red-light: #EF5753;
$red-dark: #9b4e4e;
$red-darker: #621B18;
$red-darkest: #3B0D0C;

$orange: #F6993F;
$orange-lightest: #FFF5EB;
$orange-lighter: #FCD9B6;
$orange-light: #FAAD63;
$orange-dark: #DE751F;
$orange-darker: #613B1F;
$orange-darkest: #462A16;

$green: #38C172;
$green-lightest: #E3FCEC;
$green-lighter: #A2F5BF;
$green-light: #51D88A;
$green-dark: #1F9D55;
$green-darker: #1A4731;
$green-darkest: #0F2F21;

$border-radius-none: 0;
$border-radius-sm: .125rem;
$border-radius-default: .25rem;
$border-radius-lg: .5rem;
$border-radius-full: 9999px;

.user-input-label {
    padding-bottom: .2em;
    padding-top:.4em;
}

.user-input {
    margin-bottom:.5em;
    -moz-appearance: none;
    -webkit-appearance: none;
    -webkit-box-align: center;
    -ms-flex-align: center;
    align-items: center;
    border: none;
    border-radius: 3px;
    box-shadow: none;
    display: -webkit-inline-box;
    display: -ms-inline-flexbox;
    display: inline-flex;
    font-size: 1rem;
    height: 2.25em;
    -webkit-box-pack: start;
    -ms-flex-pack: start;
    justify-content: flex-start;
    line-height: 1.25;
    padding-bottom: .5em;
    padding-left: .625em;
    padding-right: .625em;
    padding-top: .5em;
    position: relative;
    vertical-align: top;
    background-color: #fff;
    border: 1px solid #dbdbdb;
    color: #363636;
    //box-shadow: inset 0 1px 2px rgba(10,10,10,.1);
    max-width: 100%;
    width: 100%;
}

.card {
    box-shadow: 0 4px 8px 0 rgba(0,0,0,0.2);
}

.btn {
    border-radius: $border-radius-default;
    text-decoration:none;
	  text-shadow:0px 1px 0px $green-darker;
}

.btn:hover {
    filter: brightness(110%);
}

.btn-confirm {
    background-color: $green;
    border-color: darkgreen;
    color: white;
}

.error {
    border-color: red;
}
h1, h2, h3, h4, h5, h6 {
    font-weight: 100;
    margin: 0;
}

`)

			resetscss := []byte(`html, body, div, span, applet, object, iframe, p, blockquote, pre, a, abbr, acronym, address, big, cite, code, del, dfn, em, img, ins, kbd, q, s, samp, small, strike, strong, sub, sup, tt, var, b, u, i, center, dl, dt, dd, ol, ul, li, fieldset, form, label, legend, table, caption, tbody, tfoot, thead, tr, th, td, article, aside, canvas, details, embed, figure, figcaption, footer, header, hgroup, menu, nav, output, ruby, section, summary, time, mark, audio, video {
    margin: 0;
    padding: 0;
    border: 0;
    vertical-align: baseline; 
}


article, aside, details, figcaption, figure, footer, header, hgroup, menu, nav, section {
    display: block; 
}

body {
    line-height: 1; 
}

blockquote, q {
    quotes: none; 
}

blockquote {
    &:before, &:after {
    content: '';
    content: none; 
    } 
}

q {
    &:before, &:after {
    content: '';
    content: none; 
    } 
}

table {
    border-collapse: collapse;
    border-spacing: 0; 
}

button {
    background-color: transparent;
    outline: none;
    border: 0;
    cursor: pointer; 
}
`)

			indexts := []byte(`//vendor
import { VApp, Renderer, cssClass, Props } from "@kloudsoftware/eisen"

//own
import { HelloEisen } from './components/helloeisen/HelloEisen';

const app = new VApp("target", new Renderer());
app.init();

const div = app.k("div", { attrs: [cssClass("contentDiv")] }, [
    app.k("h1", { value: "Hello, eisen!" }),
]);

app.rootNode.appendChild(div);

const container = app.createElement("div", undefined, app.rootNode, [cssClass("container")]);

const routerMnt = app.createElement("div", undefined, container);

const router = app.useRouter(routerMnt);
router.registerRoute("/", new HelloEisen())
router.resolveRoute("/").catch(e => console.error(e));

`)

			componentName := "HelloEisen"
			if _, err := os.Stat("components"); os.IsNotExist(err) {
				err = os.Mkdir("components", os.ModePerm)
				check(err)
			}

			if _, err := os.Stat("components/" + strings.ToLower(componentName)); os.IsNotExist(err) {
				err = os.Mkdir("components/"+strings.ToLower(componentName), os.ModePerm)
				check(err)
			}

			component := []byte(`import { Component, VApp, ComponentBuildFunc, Props, VNode, src, cssClass } from "@kloudsoftware/eisen"

export class ` + componentName + ` extends Component {
    build(app: VApp): ComponentBuildFunc {
        return (root: VNode, props: Props) => {
            return {
                mounted: () => {
                }
            };
        }
    }
}`)

			err = ioutil.WriteFile("components/"+strings.ToLower(componentName)+"/"+componentName+".ts", component, 0644)
			check(err)
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
