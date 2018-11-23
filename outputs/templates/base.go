package templates

var HTML_BASE = `
<!doctype html>
<html>

<head>
    <title>Line Chart</title>
    <script src="https://code.highcharts.com/highcharts.js"></script>

    <style>
        canvas{
            -moz-user-select: none;
            -webkit-user-select: none;
            -ms-user-select: none;
        }
    </style>
</head>

<body>
<div style="width:95%;">
    <div id="container"></div>
	<div id="raw" style="display: none;">{{.Raw}}</div>
</div>

<script src="charts.js"></script>

</body>

</html>

`
