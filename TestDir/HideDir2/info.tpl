<!DOCTYPE html>
<html lang="zh-CH">
    <head>
        <title></title>
        <meta charset="UTF-8">
        <meta name="viewport" content="width=device-width, initial-scale=1">
    </head>
    <body>
    <h2>{{ .Info }}</h2>
    <h3>{{ .Info }}</h3>
    <h4>{{ .Info }}</h4>
    </body>

    <script type="text/javascript">
        setTimeout(window.location.href='{{ .JumpTo }}',{{ .JTS }})
    </script>
</html>