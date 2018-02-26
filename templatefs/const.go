package templatefs

const pageTemplate = `
<!DOCTYPE html>
<html>
<head>
  <title>{{ .Title }}</title>
  <meta charset="utf-8">
  {{- if .CSS}}
  <link rel="stylesheet" type="text/css" href="{{ .CSS }}">
  {{- end}}
</head>
<body>
{{ .Body | safehtml }}
</body>
</html>
`
