
# web.go

web.go - is a web framework written in go with the main purpose of writing less javascript (or any other web language tbh) and more go when building a project that requires a frontend. This could technically be extended by any language if you add some grpc glue in between. Please be aware, this is current a poc that can and will break very easily. You are warned.

## Structure

Let's start of by taking a look at the actual structure of a project using web.go. This makes it easier to understand the rest of the process.

```
├── foo
│   ├── build -> auto generated folder that contains generated html, js, build scripts, … . 
|   |   ├── runtime.json -> required by web server to be able to serve the correct endpoints and data.
|   |   ├── store.html -> formatted into valid html of your version.
|   |   └── wg.js -> contains default & generated javascript code.
│   ├── main.go -> file containing your main function (can also contain routes & actions)
│   ├── settings.yaml
│   ├── index.html -> index template, that contains things such as stylehseets, custom js links, … .
│   ├── store.html -> (*.html) other html files are pages or components. I guess these could be given a custom file extension, but meh.
│   ├── go.mod
│   └── go.sum
```

## How does it work?
Everything starts with the parsing of the html pages. Since these can contain tokens. The tokens represent data structures, and they can contain any type of data; strings, integers, boolean’s, html components. Basically anything you want and can do in go and as in a javascript framework. 

```html
-articles
 <div>
    <p>{articles.id}</p>
    <p>{articles.name}</p>
    <p>{articles.price}</p>
    
    <p>{action:opt}</p>
    <button type="button" onclick="{remove(articles.id)}">Buy</button>	
</div>
```

These tokens will be used to build an internal structure from the application that will go through a few processes. Parsing tokens, retrieving relationships between structures, formatting the html code to be valid, mapping parent tokens to functions and actions, generating library code which includes the generation of the entire javascript side of things. 

All of this generates the content of your build directory, which contains all your app’s static files that can be hosted with any web server you prefer (nginx for example). 

Of course you still need to your run your actual logic, which is written entirely the way you prefer.  All your functions need are 2 or 3 parameters (depending if your function represents a data retrieval or action)

```go
func RetrieveArticles(w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	...
	return "articles", []Articles, nil
}

func Buy(args map[string]interface{}, w http.ResponseWriter, r *http.Request) (string, interface{}, error) {
	...
	return "result", "Thank you for ordering.", nil
}
```

Afterwards your main function needs to run the following function to start the backend. It requires the function pointers of each function and action. 

```go
func main() {
    ...

    functions := []web.Function{RetrieveArticles}
    actions := []web.Action{Buy}

    web.Start(os.Args, functions, actions)
}
```

Now your backend will serve the data to your static (but dynamic :thinking:) frontend, by http calls. The data will be retrieved and rendered in the correct position within the html code.

This needs to be ran using the wg (web.go cli command), since we have three runtime options; generate, production run & development mode. Generate obviously generates the static files and backend runtime file. Production run requires the main.go and location of your runtime file. Development mode allows peacefully write code and test your application locally without having to run multiple commands. 

Here you have a small overview of what web.go is and how it kind of works. Now please stay stuned for some more examples and functionality.
