import * as d3 from 'd3';
import WebSocket from 'reconnecting-websocket';

/* Data */

const list = [];

/* Visualization */

const width = window.innerWidth - 100;
const height = window.innerHeight - 100;

const svg = d3.select('#viz');

const xAxis = svg.append('g');
const yAxis = svg.append('g');

const path = svg.append('path');

function viz() {
  if (list.length === 0) {
    return;
  }

  const x = d3
    .scaleTime()
    .domain(d3.extent(list, (d) => d.time))
    .range([0, width]);

  xAxis.attr('transform', 'translate(50,' + (50 + height) + ')').call(d3.axisBottom(x));

  const y = d3
    .scaleLinear()
    .domain([0, d3.max(list, (d) => d3.max(d.sample))])
    .range([height, 0]);

  yAxis.attr('transform', 'translate(50,50)').call(d3.axisLeft(y));

  // Add the line
  path
    .datum(list)
    .attr('fill', 'none')
    .attr('stroke', 'steelblue')
    .attr('stroke-width', 1.5)
    .attr(
      'd',
      d3
        .line()
        .x(function (d) {
          return 50 + x(d.time);
        })
        .y(function (d) {
          return 50 + y(d.sample[0]);
        })
    );
}

/* WebSocket */

const ws = new WebSocket('ws://0.0.0.0:8080');

ws.onmessage = (event) => {
  // parse sample
  const sample = JSON.parse(event.data.toString());

  // prepare entry
  const entry = {
    sample: sample,
    time: new Date(),
  };

  // add to list
  list.push(entry);

  // trim list
  while (list.length > 1000) {
    list.shift();
  }

  // update
  viz();
};
