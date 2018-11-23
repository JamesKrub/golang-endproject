package main

import "html/template"

var editTemplate = template.Must(template.New("edit").Parse(`<!doctype html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>Edit Post</title>
  </head>
  <h2>Edit Post</h2>
  <form method="POST" action="/posts/{{.ID}}/update">
    <h3>Title</h3>
    <input name="title" value="{{.Title}}"/>
    <h3>Body</h3>
    <textarea name="body">{{.Body}}</textarea></br>
    <input type="submit" />
  </form>
</html>
`))
