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

			err := os.Chdir("src")
			check(err)

			componentName := os.Args[3]
			if _, err := os.Stat("components"); os.IsNotExist(err) {
				err = os.Mkdir("components", os.ModePerm)
				check(err)
			}

			err = os.Chdir("components")
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
			packagejson := []byte(`{
  "name": "` + projectName + `",
  "version": "1.0.0",
  "description": "An eisen Project",
  "main": "index.ts",
  "scripts": {
    "dev": "npx parcel src/index.html",
    "build": "npx parcel build src/index.html",
    "clean": "rm -rf dist/ && rm -rf lib/" },
  "author": "",
  "license": "ISC",
    "dependencies": {
    "@kloudsoftware/eisen": "1.0.28",
    "postcss": "^7.0.18",
    "tailwindcss": "^1.1.2"
  },
  "devDependencies": {
    "parcel-bundler": "^1.12.4",
    "parcel-plugin-purgecss": "^2.1.2",
    "parcel-plugin-static-files-copy": "^2.0.0",
    "sass": "^1.18.0",
    "typescript": "^3.4.2"
  }
}`)
			gitignore := []byte("node_modules/\nlib/\ndist\n.cache/\n.#*")

			err = ioutil.WriteFile("package.json", packagejson, 0644)
			check(err)
			err = ioutil.WriteFile(".gitignore", gitignore, 0644)
			check(err)

			if _, err := os.Stat("static"); os.IsNotExist(err) {
				err = os.Mkdir("static", os.ModePerm)
				check(err)
			}

			dockerfile := []byte(`FROM voidlinux/voidlinux
RUN xbps-install -Syu nodejs nginx curl
WORKDIR /app
ADD src/ src/
ADD package.json .
ADD static/ static/
RUN npm install && npm run build &&  mv /app/dist/* /usr/share/nginx/html
RUN curl https://gist.githubusercontent.com/larsgrah/6003204aaf2b32b885682e7fe94c0ed8/raw/62b96d8820b62220d618fb4edb00223862c66127/nginx.conf -o /etc/nginx/nginx.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]`)

			err = ioutil.WriteFile("Dockerfile", dockerfile, 0644)
			check(err)

			postcss := []byte(`module.exports = {
    plugins: [
        require('tailwindcss'),
        require('autoprefixer'),
    ]
};`)

			purgecss := []byte(`module.exports = {
    content: ["**/*.ts"],
};`)
			tailwindcss := []byte(`module.exports = {
};`)
			tsconfig := []byte(`{
  "compilerOptions": {
    "module": "commonjs",
    "moduleResolution": "node",
    "newLine": "LF",
    "outDir": "./lib/",
    "target": "es5",
    "sourceMap": true,
    "declaration": true,
    "jsx": "preserve",
    "lib": [
      "es2017",
      "dom",
      "esnext"
    ]
  },
  "include": [
    "src/**/*"
  ],
  "exclude": [
    ".git",
    "node_modules"
  ]
}`)

			err = ioutil.WriteFile("postcss.config.js", postcss, 0644)
			check(err)
			err = ioutil.WriteFile("purgecss.config.js", purgecss, 0644)
			check(err)
			err = ioutil.WriteFile("tailwind.config.js", tailwindcss, 0644)
			check(err)
			err = ioutil.WriteFile("tsconfig.json", tsconfig, 0644)
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
			stylescss := []byte(`
@tailwind base;

@tailwind components;

@tailwind utilities;
`)

			indexts := []byte(`//vendor
import { VApp, Renderer, cssClass, Props } from "@kloudsoftware/eisen"

//own
import { HelloEisen } from './components/helloeisen/HelloEisen';

const app = new VApp("target", new Renderer());
app.init();

const container = app.createElement("div", undefined, app.rootNode, [cssClass("container")]);

const routerMnt = app.createElement("div", undefined, container);

const router = app.useRouter(routerMnt);
router.registerRoute("/", new HelloEisen())
router.resolveRoute(document.location.pathname).catch(() => router.resolveRoute("/"));

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
            const div = app.k("div", { attrs: [cssClass("contentDiv")] }, [
                app.k("h1", { value: "Hello, eisen!" }),
            ]);
            
            root.appendChild(div);
            return {
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
			err = ioutil.WriteFile("index.ts", indexts, 0644)
			check(err)

		}
	}
}
