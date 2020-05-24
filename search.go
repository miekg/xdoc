package main

const htmlForm = `<html>
<head>
<link rel="stylesheet" type="text/css" href="/a/style.css"/>
<title>Search</title>
</head>

<body>
<form action="/s/" method="GET">
	<input type="text" placeholder="query" name="q" id="search-input"/>
	<button type="submit">Search</button>
</form>
</body>
</html>
`
