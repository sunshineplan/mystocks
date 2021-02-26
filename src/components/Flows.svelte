<script lang="ts">
  import { onMount } from "svelte";
  import Chart from "chart.js";
  import AutoComplete from "./AutoComplete.svelte";
  import { getToday, checkTime, labels, getColor } from "../misc";

  const onClick = (
    event: MouseEvent,
    legendItem: Chart.ChartLegendLabelItem
  ) => {
    const index = legendItem.datasetIndex as number;
    const current = chart.getDatasetMeta(index);
    const alreadyHidden =
      current.hidden === undefined || current.hidden === null
        ? false
        : current.hidden;
    let anyOthersAlreadyHidden = false;
    let allOthersHidden = true;

    (chart.data.datasets as Chart.ChartDataSets[]).forEach((e, i) => {
      const meta = chart.getDatasetMeta(i);
      if (i !== index) {
        if (meta.hidden) {
          anyOthersAlreadyHidden = true;
        } else {
          allOthersHidden = false;
        }
      }
    });
    if (alreadyHidden) {
      current.hidden = undefined;
    } else {
      (chart.data.datasets as Chart.ChartDataSets[]).forEach((e, i) => {
        const meta = chart.getDatasetMeta(i);
        if (i !== index) {
          if (anyOthersAlreadyHidden && !allOthersHidden) {
            meta.hidden = true;
          } else {
            meta.hidden =
              meta.hidden === undefined || meta.hidden === null
                ? !meta.hidden
                : undefined;
          }
        } else {
          meta.hidden = undefined;
        }
      });
    }
    chart.update();
  };

  const capitalflows = {
    type: "line",
    data: {
      labels,
    },
    options: {
      maintainAspectRatio: false,
      legend: {
        position: "right",
        onClick,
      },
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

  const today = getToday();

  let autoUpdate = 0;
  let chart: Chart;
  let date = today;
  let last = "";

  const load = async (force?: boolean, date?: string) => {
    var url = "/flows";
    if (date) url = url + `?date=${date}`;
    if (checkTime() || force) {
      const resp = await fetch(url);
      const array = await resp.json();
      if (array && array.length) {
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
        if (!date || (date && date == today))
          last = new Date().toLocaleString();
      }
    }
  };

  const goto = (day?: string) => {
    if (day && day != today) {
      if (autoUpdate) clearInterval(autoUpdate);
      load(true, day);
    } else {
      last = "";
      date = today;
      load(true);
      autoUpdate = setInterval(load, 60000);
    }
  };

  onMount(() => {
    chart = new Chart(
      document.querySelector("#flowsChart") as HTMLCanvasElement,
      capitalflows
    );

    load(true);
    autoUpdate = setInterval(load, 60000);

    return () => {
      if (autoUpdate) clearInterval(autoUpdate);
      chart.destroy();
    };
  });
</script>

<svelte:head>
  <title>My Stocks</title>
</svelte:head>

<header style="height:80px">
  <AutoComplete />
  <div class="input-group">
    <div class="input-group-prepend">
      <label class="input-group-text" for="date">Date</label>
    </div>
    <input class="form-control" type="date" bind:value={date} id="date" />
    <div class="input-group-append">
      <button class="btn btn-primary" on:click={() => goto(date)}>Go</button>
    </div>
    <div class="input-group-append">
      <button class="btn btn-danger" on:click={() => goto()}>Reset</button>
    </div>
  </div>
  {#if date == today && last}
    <small>Last update: {last}</small>
  {/if}
</header>
<div class="chart">
  <canvas id="flowsChart" />
</div>

<style>
  .input-group {
    width: 350px;
  }

  .chart {
    height: calc(100% - 80px);
  }
</style>
