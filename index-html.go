package main

import "html/template"

var indexTemplate = template.Must(template.New("index").Parse(`<!doctype html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>Blog</title>
  </head>
  <body>
  <h1>Blog</h1>
  <div><a href="/posts/new">New Post</a></div>
  <h2>Posts</h2>
  {{range .Posts}}<div><a href="/posts/{{.ID}}">{{ .Title }}</a></div>{{end}}
  </body>
</html>
`))
