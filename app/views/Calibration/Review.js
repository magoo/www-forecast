var pageData = getPageData();

function prepGraphData() {
    // each answer has [correct?, confidence]
    var answers = pageData.answers;

    var buckets = [.6, .7, .8, .9, .95, .99].map(function (confidence) {
        return {
            confidence: confidence,
            numCorrect: 0,
            numIncorrect: 0,
        }
    })

    // Tally up the answers/outcomes
    answers.forEach(function (answer) {
        console.log(JSON.stringify(buckets), answer.Confidence)
        bucket = buckets.find(function(b){ return b.confidence == answer.Confidence })
        if (answer.Outcome == true) {
            bucket.numCorrect += 1
        } else {
            bucket.numIncorrect += 1
        }
    })

    // grant each bucket a fractionCorrect value for the graph
    buckets.forEach(function (bucket) {
        numTotal = bucket.numCorrect + bucket.numIncorrect;
        bucket.fractionCorrect = (1.0 * bucket.numCorrect) / numTotal;
    })

    buckets = buckets.filter(function(bucket){
        return !('NaN' == '' + bucket.fractionCorrect)
    })

    return buckets
}

function drawGraph(data) {
    var margin = {top: 20, right: 20, bottom: 20, left: 30};

    var width = 300// - margin.right - margin.left;
    var height = 300// - margin.top - margin.bottom;

    var x = d3.scaleLinear()
        .domain([.6, 1.0])
        .range([0, width]);

    var y = d3.scaleLinear()
        .domain([0, 1.0])
        .range([height, 0]);

    var xAxis = d3.axisBottom(x)
        .ticks(5);

    var yAxis = d3.axisLeft(y)
        .ticks(5);

    var svg = d3.select('#confidence-graph')
        .append('svg')
        .attr('width', width + margin.left + margin.right)
        .attr('height', height + margin.top + margin.bottom)
        .append("g")
        .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

    svg.append('g')
        .attr('class', 'x axis')
        .attr('transform', 'translate(0,' + height + ')')
        .call(xAxis);

    svg.append('g')
        .attr('class', 'y axis')
        .call(yAxis);

    var line = d3.line()
        .x(d => x(d.confidence))
        .y(d => y(d.fractionCorrect))

    console.log(JSON.stringify(data))

    svg.append("path")
        .datum(data)
        .attr("fill", "none")
        .attr("stroke", "steelblue")
        .attr("stroke-width", 1.5)
        .attr("stroke-linejoin", "round")
        .attr("stroke-linecap", "round")
        .attr("d", line);

    svg.selectAll("dot")
        .data(data)
        .enter().append("circle")
        .attr("stroke", "steelblue")
        .attr("fill", "steelblue")
        .attr("r", 3.5)
        .attr("cx", function(d) { return x(d.confidence )})
        .attr("cy", function(d) { return y(d.fractionCorrect )})

    var idealData = [
        { confidence: 0.6, fractionCorrect: 0.6 },
        { confidence: 0.7, fractionCorrect: 0.7 },
        { confidence: 0.8, fractionCorrect: 0.8 },
        { confidence: 0.9, fractionCorrect: 0.9 },
        { confidence: 0.95, fractionCorrect: 0.95 },
        { confidence: 0.99, fractionCorrect: 0.99 },
    ]

    svg.append("path")
        .datum(idealData)
        .attr("fill", "none")
        .attr("stroke", "steelblue")
        .attr("stroke-width", 1.5)
        .attr("stroke-linejoin", "round")
        .attr("stroke-linecap", "round")
        .style("stroke-dasharray", ("3, 3"))
        .attr("d", line);
}

$(function(){
    drawGraph(prepGraphData())
})
