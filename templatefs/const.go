package templatefs

const pageTemplate = `
<!DOCTYPE html>
<html>
<head>
  <title>{{ .Title }}</title>
  <meta charset="utf-8">
</head>
<body>
{{ .Body | safehtml }}
</body>
</html>
`
