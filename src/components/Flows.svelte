<script lang="ts">
  import { onMount } from "svelte";
  import Chart from "chart.js";
  import AutoComplete from "./AutoComplete.svelte";
  import { checkTime, labels, getColor } from "../misc";

  const capitalflows = {
    type: "line",
    data: {
      labels,
    },
    options: {
      legend: { position: "right" },
      animation: { duration: 0 },
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
            id: "zero",
            type: "line",
            mode: "horizontal",
            scaleID: "y-axis-0",
            value: 0,
            borderColor: "black",
            borderWidth: 0.75,
          },
        ],
      },
    },
  } as Chart.ChartConfiguration;

  let autoUpdate = 0;
  let chart: Chart;

  const start = () => {
    chart = new Chart(
      document.querySelector("#flowsChart") as HTMLCanvasElement,
      capitalflows
    );
    load();
    autoUpdate = setInterval(load, 60000);
  };

  const load = async () => {
    if (checkTime()) {
      const resp = await fetch("/flows");
      const array = await resp.json();
      if (array.length) {
        const datasets = chart.data.datasets as Chart.ChartDataSets[];
        datasets.length = 0;
        array.forEach((e: any, i: number) => {
          datasets.push({
            label: e.sector,
            fill: false,
            lineTension: 0,
            borderWidth: 1.5,
            borderColor: getColor(i),
            backgroundColor: getColor(i),
            pointRadius: 0,
            pointHoverRadius: 3,
            data: e.chart,
          });
        });
        chart.update();
      }
    }
  };

  onMount(() => {
    start();
    return () => {
      clearInterval(autoUpdate);
      chart.destroy();
    };
  });
</script>

<svelte:head>
  <title>My Stocks</title>
</svelte:head>

<header style="height:60px">
  <AutoComplete />
</header>
<canvas class="chart" id="flowsChart" />

<style>
  .chart {
    max-width: 1440px;
    max-height: 720px;
    height: calc(100% - 210px);
  }
</style>
