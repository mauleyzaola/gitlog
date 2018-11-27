function sorter(a,b){
    return (a.y < b.y) ? 1 : -1;
}

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


function groupCommitsByMonth(commits){
    var accum = {};
    (commits || []).forEach(function(x) {
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


function calcCommitPerMonthAccum(data){
    var series = [];
    data.forEach(function(s) {
        var res = groupCommitsByMonth(s.commits || []).map(function(i){
            return Object.assign({}, {
                x: (moment(dayToDate(i.key)).unix()) * 1000,
                y: i.value,
            });
        })
        var sum = 0;
        series.push(
            {
                name: s.name,
                data: res.map(function(i){
                    sum += i.y;
                    return {
                        x: i.x,
                        y: sum,
                    }
                }),
            },
        );
    });
    return series;
}


function calcCommitPerMonth(data){
    var series = [];
    data.forEach(function(s) {
        var res = groupCommitsByMonth(s.commits || []).map(function(i){
            return Object.assign({}, {
                x: (moment(dayToDate(i.key)).unix()) * 1000,
                y: i.value,
            });
        })
        series.push(
            {
                name: s.name,
                data: res,
            },
        );
    });
    return series;
}

function calcAuthorDistributionSummary(data){
    // TODO: implement others when there are too many authors
    var res = [];
    var authors = {};
    data.forEach(function(r){
        (r.commits || []).forEach(function(c){
            var author = authors[c.author.email];
            if(!author){
                author = { count: 0};
                authors[c.author.email] = author;
            }
            author.count++;
        })
    })

    for(var key in authors){
        res.push({
            name: key, // TODO: split @
            y: authors[key].count,
        })
    }
    res.sort(sorter);
    return res;
}

function calcFileTypeDistribution(data){
    var res = [];
    var extensions = {};
    data.forEach(function(r){
        for(let ext in r.fileStat){
            var extension = extensions[ext];
            if(!extension){
                extension = { size:0, count: 0 };
                extensions[ext] = extension;
            }
            extension.size += r.fileStat[ext].size;
            extension.count += r.fileStat[ext].count;
        }
    });

    var sizes = [];
    var counts = [];

    for(let ext in extensions){
        sizes.push({
            name: ext,
            y: extensions[ext].size,
        });
        counts.push({
            name: ext,
            y: extensions[ext].count,
        })
    }

    sizes.sort(sorter);
    counts.sort(sorter);

    return {
        sizes,
        counts,
    }
}

function drawCommitsYearMonthTimeline(params){
    Highcharts.chart(params.element, {
        credits: false,
        chart: {
            type: 'spline'
        },
        title: {
            text: params.title,
        },
        subtitle: {
            text: params.subtitle,
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
        series: params.series,
    });
}

function drawAuthorSummary(params){
    Highcharts.chart(params.element, {
        credits: false,
        chart: {
            plotBackgroundColor: null,
            plotBorderWidth: null,
            plotShadow: false,
            type: 'pie'
        },
        title: {
            text: params.title
        },
        subtitle: {
            text: params.subtitle,
        },
        tooltip: {
            pointFormat: '{series.name}: <b>{point.percentage:.1f}%</b>'
        },
        plotOptions: {
            pie: {
                allowPointSelect: true,
                cursor: 'pointer',
                dataLabels: {
                    enabled: true,
                    format: '<b>{point.name}</b>: {point.percentage:.1f} %',
                    style: {
                        color: (Highcharts.theme && Highcharts.theme.contrastTextColor) || 'black'
                    }
                }
            },
        },
        series: [{
            name: 'Commits',
            colorByPoint: true,
            data: params.data,
        }]
    });
}

function drawFileTypeCountSummary(params){
    Highcharts.chart(params.element, {
        credits: false,
        chart: {
            plotBackgroundColor: null,
            plotBorderWidth: null,
            plotShadow: false,
            type: 'pie'
        },
        title: {
            text: params.title
        },
        subtitle: {
            text: params.subtitle,
        },
        tooltip: {
            pointFormat: '{series.name}: <b>{point.percentage:.1f}%</b>'
        },
        plotOptions: {
            pie: {
                allowPointSelect: true,
                cursor: 'pointer',
                dataLabels: {
                    enabled: true,
                    format: '<b>{point.name}</b>: {point.percentage:.1f} %',
                    style: {
                        color: (Highcharts.theme && Highcharts.theme.contrastTextColor) || 'black'
                    }
                }
            }
        },
        series: [{
            name: 'Types',
            colorByPoint: true,
            data: params.data,
        }]
    });
}

var raw = document.getElementById('raw');
var data = JSON.parse(raw.innerHTML) || [];

drawCommitsYearMonthTimeline({
    series: calcCommitPerMonth(data),
    title: 'Commits by Year/Month',
    subtitle: 'for each repository',
    element: 'container-monthly-commits',
});

drawCommitsYearMonthTimeline({
    series: calcCommitPerMonthAccum(data),
    title: 'Commits by Year/Month',
    subtitle: 'Aggregate sum for each repository',
    element: 'container-monthly-commits-accum',
});

drawAuthorSummary({
    element: 'container-author-distribution',
    title: 'Distribution of Authors',
    subtitle: 'by counting the number of commits in all repositories',
    data: calcAuthorDistributionSummary(data),
});

var typeData = calcFileTypeDistribution(data);
drawFileTypeCountSummary({
    element: 'container-file-type-count-distribution',
    title: 'Distribution of File Types',
    subtitle: 'by counting the number of files for each type',
    data: typeData.counts,
});
drawFileTypeCountSummary({
    element: 'container-file-type-size-distribution',
    title: 'Distribution of File Types',
    subtitle: 'by adding the size for each file',
    data: typeData.sizes,
});