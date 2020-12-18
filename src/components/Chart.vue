<template>
  <canvas class="chart" id="chart" />
</template>

<script>
import { timeLabels } from "@/misc.js";

export default {
  name: "StockChart",
};

export const intraday = {
  type: "line",
  data: {
    labels: timeLabels(9 * 60 + 30, 11 * 60 + 30).concat(
      timeLabels(13 * 60 + 1, 15 * 60)
    ),
    datasets: [
      {
        label: "Price",
        fill: false,
        lineTension: 0,
        borderWidth: 2,
        borderColor: "red",
        backgroundColor: "red",
        pointRadius: 0,
        pointHoverRadius: 3,
      },
    ],
  },
  options: {
    scales: {
      xAxes: [
        {
          gridLines: { drawTicks: false },
          ticks: {
            padding: 10,
            autoSkipPadding: 100,
            maxRotation: 0,
          },
        },
      ],
      yAxes: [
        {
          gridLines: { drawTicks: false },
          ticks: { padding: 12 },
        },
      ],
    },
    annotation: {
      annotations: [
        {
          id: "PreviousClose",
          type: "line",
          mode: "horizontal",
          scaleID: "y-axis-0",
          borderColor: "black",
          borderWidth: 0.75,
        },
      ],
    },
  },
};
</script>

<style scoped>
.chart {
  max-width: 1000px;
  max-height: 500px;
  height: calc(100% - 210px);
}
</style>
