function dateToDay(row){
    var date = new Date(row);
    var year = date.getUTCFullYear().toString();
    var month = (date.getUTCMonth() + 1).toString();
    var day = (date.getDate()).toString();
    month = ((month.length === 1) ? '0' : '') + month;
    // day = ((day.length === 1) ? '0' : '') + day;
    day = '01';
    return year + month + day;
}

function dayToDate(val){
    var year = val.toString().substr(0,4);
    var month = val.toString().substr(4,2);
    var day = val.toString().substr(6,2);
    return new Date(Date.parse(month + '/' + '/' + day + '/' + year));
}


function transformCollection(commits){
    var accum = {};
    commits.forEach(function(x) {
        var key = dateToDay(x.date);
        var curr = accum[key] || 0;
        accum[key] = ++curr;
    });

    var series = [];

    for(var key in accum){
        series.push({
            key,
            value: accum[key],
        });
    }

    return series;
}


function transform(data){
    data = JSON.parse(data);
    var series = [];
    data.forEach(function(s) {
        var res = transformCollection(s.commits).map(function(i){
            return Object.assign({}, {
                x: (moment(dayToDate(i.key)).unix()) * 1000,
                y: i.value,
            });
        })
        series.push({
            name: s.name,
            data: res,
        });
    });
    return {
        series,
    };
}

function draw(data){
    Highcharts.chart('container', {
        credits: false,
        chart: {
            type: 'spline'
        },
        title: {
            text: 'Commits by Year/Month'
        },
        subtitle: {
            text: 'Aggregate sum for each repository'
        },
        xAxis: {
            type: 'datetime',
            title: {
                text: '',
            },
        },
        yAxis: {
            title: {
                text: 'Number of Commits'
            },
        },
        plotOptions: {
            spline: {
                marker: {
                    enabled: true
                }
            },
        },
        series: data.series,
    });
}

var raw = document.getElementById('raw');
draw(transform(raw.innerHTML));