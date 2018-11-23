package main

import "html/template"

var postTemplate = template.Must(template.New("post").Parse(`<!doctype html>
<html>
  <head>
    <meta charset="UTF-8">
    <title>{{.Title}}</title>
  </head>
  <body>
    <div><a href="/posts/{{.ID}}/edit">Edit</a></div>
    <h1>{{.Title}}</h1>
    <p>
    {{.Body}}
    </p>
    <h3>Add Comment</h3>
    <form method="POST" action="/posts/1/comment">
      <textarea name="body"></textarea><br />
      <input type="submit" value="Comment" />
    </form>
    <h2>Comments</h2>
    {{range .Comments}}
    <p>
      {{.Body}}
    </p>
    <hr/>
    {{end}}
  </body>
</html>
`))
