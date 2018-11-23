package templates

var JS_COMMITS = `

// data must be an array of gitlog commits
function transform(data){
    data = JSON.parse(data);
    var accum = {};

    var monthKey = function(row){
        var date = new Date(row);
        var year = date.getUTCFullYear().toString();
        var month = (date.getUTCMonth() + 1).toString();
        if(month.length === 1){
            month = "0" + month;
        }
        return year + "-" + month;
    }

    data.forEach(function(x) {
        var m = monthKey(x.date);
        var curr = accum[m] || 0;
        accum[m] = ++curr;
    });

    var categories = [];
    var series = [];

    for(var key in accum){
        categories.push(key);
        series.push(accum[key]);
    }

    return {
        categories:categories,
        series:series,
    }
}

function draw(data){
    Highcharts.chart('container', {
        chart: {
            type: 'line'
        },
        title: {
            text: 'Commits Summary'
        },
        subtitle: {
            text: '.'
        },
        xAxis: {
            categories: data.categories,
        },
        yAxis: {
            title: {
                text: 'Number of Commits'
            }
        },
        plotOptions: {
            line: {
                dataLabels: {
                    enabled: true
                },
                enableMouseTracking: false
            }
        },
        series: [{
            name: 'Everyone',
            data: data.series,
        }, ]
    });
}

var raw = document.getElementById('raw');
draw(transform(raw.innerHTML));
`
