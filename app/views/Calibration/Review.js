function prepGraphData() {
    var pageData = getPageData();
    if (!pageData.answers) {
        throw Error("There are no answers in the pageData!", JSON.stringify(pageData))
    }

    // each answer has [correct?, confidence]
    var answers = pageData.answers;

    var buckets = [0, .1, .2, .3, .4, .5, .6, .7, .8, .9, 1.0].map(function (confidence) {
        return {
            confidence: confidence,
            numCorrect: 0,
            numIncorrect: 0,
        }
    });

    answers.forEach(function (answer) {
        // Normalize values so false answers become inverted confidence of a true answer instead
        if (!answer.Answer) {
            answer.Answer = !answer.Answer;
            answer.Confidence = 1.0 - answer.Confidence;
        }

        // Tally up the answers/outcomes
        answer.Confidence = Math.round(answer.Confidence * 100) / 100; // Fix floating point precision errors
        var bucket = buckets.find(function(b){ return b.confidence === answer.Confidence })
        if (!bucket) { debugger }
        if (answer.Outcome === true) {
            bucket.numCorrect += 1
        } else {
            bucket.numIncorrect += 1
        }
    });

    // grant each bucket a fractionCorrect value for the graph
    buckets.forEach(function (bucket) {
        var numTotal = bucket.numCorrect + bucket.numIncorrect;
        bucket.fractionCorrect = bucket.numCorrect / numTotal;
    });

    buckets = buckets.filter(function(bucket){
        return !('NaN' === '' + bucket.fractionCorrect)
    });

    return buckets
}

function drawGraph(data) {
    var margin = {top: 20, right: 20, bottom: 20, left: 30};

    var width = 300; // - margin.right - margin.left;
    var height = 300; // - margin.top - margin.bottom;

    var x = d3.scaleLinear()
        .domain([0, 1.0])
        .range([0, width]);

    var y = d3.scaleLinear()
        .domain([0, 1.0])
        .range([height, 0]);

    var xAxis = d3.axisBottom(x)
        .ticks(10);

    var yAxis = d3.axisLeft(y)
        .ticks(10);

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
        { confidence: 0, fractionCorrect: 0 },
        { confidence: 0.1, fractionCorrect: 0.1 },
        { confidence: 0.2, fractionCorrect: 0.2 },
        { confidence: 0.3, fractionCorrect: 0.3 },
        { confidence: 0.4, fractionCorrect: 0.4 },
        { confidence: 0.5, fractionCorrect: 0.5 },
        { confidence: 0.6, fractionCorrect: 0.6 },
        { confidence: 0.7, fractionCorrect: 0.7 },
        { confidence: 0.8, fractionCorrect: 0.8 },
        { confidence: 0.9, fractionCorrect: 0.9 },
        { confidence: 1.0, fractionCorrect: 1.0 },
    ];

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
