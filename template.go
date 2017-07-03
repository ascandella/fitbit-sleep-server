package main

import "html/template"

var sleepTemplate = template.Must(template.New("main").Parse(`
<head>
<title>the sleep of ai</title>
<link rel="stylesheet" type="text/css" href="https://necolas.github.io/normalize.css/latest/normalize.css" />
</head>
<body>
  <h2>Date<h2>
  <p>{{ .Date }} </p>

  <h2>Time</h2>
  <div>{{ .MostRecent }}</div>
</body>
`))
