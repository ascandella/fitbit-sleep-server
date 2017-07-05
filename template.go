package main

import "html/template"

var sleepTemplate = template.Must(template.New("main").Parse(`
<head>
<title>the sleep of ai</title>
<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css" integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">

<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap-theme.min.css" integrity="sha384-rHyoN1iRsVXV4nD0JutlnGaslCJuC7uwjduW9SVrLvRYooPp2bWYgmgJQIXwl/Sp" crossorigin="anonymous">

<meta property="og:title" content="aiden's sleep" />
<meta property="og:description" content="{{ .FriendlyDuration }}" />

</head>
<body>

	<div class="container">
		<h3 class="page-header">How much sleep did aiden get?</h3>
		<p class="lead">{{ .FriendlyDuration }}</p>

		<h3>Start</h3>
		<p class="lead">{{ .Start }}</p>
	</div>
	<script>
		(function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
		(i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
		m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
		})(window,document,'script','https://www.google-analytics.com/analytics.js','ga');

		ga('create', 'UA-102024160-1', 'auto');
		ga('send', 'pageview');
	</script>
</body>
`))
