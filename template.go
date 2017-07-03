package main

import "html/template"

var sleepTemplate = template.Must(template.New("main").Parse(`
<head>
<title>the sleep of ai</title>
</head>
<body>
  <div>{{ . }}</div>
</body>
`))
