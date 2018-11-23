package main

import "html/template"

var newTemplate = template.Must(template.New("new").Parse(`<!doctype html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>New Post</title>
  </head>
  <body>
    <h2>New Post</h2>
    <form method="POST" action="/posts/">
      <h3>Title</h3>
      <input name="title" />
      <h3>Body</h3>
      <textarea name="body"></textarea></br>
      <input type="submit" />
    </form>
  </body>
</html>
`))
