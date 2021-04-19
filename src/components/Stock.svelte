<script lang="ts">
  import { onMount } from "svelte";
  import { Chart } from "chart.js";
  import AutoComplete from "./AutoComplete.svelte";
  import Realtime from "./Realtime.svelte";
  import { checkTime, post, intraday } from "../misc";
  import { component, current, refresh } from "../stores";
  import type { Stock } from "../stores";
  import type {
    ChartItem,
    ScatterDataPoint,
    TooltipCallbacks,
    LinearScaleOptions,
  } from "chart.js";
  import type { AnnotationOptions } from "chartjs-plugin-annotation";

  let autoUpdate: number[] = [];
  let stock: Stock = {
    index: "n/a",
    code: "n/a",
    name: "n/a",
    now: 0,
    change: 0,
    percent: "-",
    high: 0,
    low: 0,
    open: 0,
    last: 0,
    sell5: [],
    buy5: [],
    update: "",
  };
  let data: ScatterDataPoint[] = [];
  let chart: Chart<"line">;
  let update = "";
  let hover = false;

  $: $current, load();

  const y = intraday.options?.scales?.y as LinearScaleOptions;
  const y2 = intraday.options?.scales?.y2 as LinearScaleOptions;

  y.afterFit = (axis) => {
    if (chart) {
      const y2 = chart.options?.scales?.y2 as LinearScaleOptions;
      y2.min = axis.min;
      y2.max = axis.max;
    }
  };

  y2.ticks.callback = (value) => {
    if (stock.last)
      return `${
        Math.round((((value as number) - stock.last) / stock.last) * 10000) /
        100
      }%`;
    return null;
  };

  const callbacks = intraday.options?.plugins?.tooltip
    ?.callbacks as TooltipCallbacks<"line">;

  callbacks.label = (tooltipItem) => {
    const value = tooltipItem.parsed.y;
    const percent =
      Math.round(((value - stock.last) / stock.last) * 10000) / 100;
    return `${value}   ${percent}%`;
  };
  callbacks.labelTextColor = (tooltipItem) => {
    const change = parseFloat(tooltipItem.formattedValue) - stock.last;
    if (change > 0) return "red";
    else if (change < 0) return "green";
    return "black";
  };

  const start = () => {
    chart = new Chart(
      document.querySelector("#stockChart") as ChartItem,
      intraday
    );
    if ($current.code != "n/a") {
      autoUpdate.push(setInterval(loadRealtime, $refresh * 1000));
      autoUpdate.push(setInterval(loadChart, 60000));
    }
  };

  const load = async () => {
    update = "";
    await loadRealtime(true);
    await loadChart(true);
    const yAxes = chart.options?.scales?.y as LinearScaleOptions;
    yAxes.suggestedMin = stock.last / 1.01;
    yAxes.suggestedMax = stock.last * 1.01;
    const annotations = chart.options.plugins?.annotation
      ?.annotations as AnnotationOptions[];
    annotations[0].value = stock.last;
    updateChart(true);
  };

  const loadRealtime = async (force?: boolean) => {
    if (checkTime() || (force && $current.code)) {
      const resp = await post("/get", {
        index: $current.index,
        code: $current.code,
        q: "realtime",
      });
      const json = await resp.json();
      if (json.name) {
        stock = json;
        update = json.update;
        if (update && !force) updateChart();
        document.title = `${json.name} ${json.now} ${json.percent}`;
      }
    }
  };

  const loadChart = async (force?: boolean) => {
    if (checkTime() || (force && $current.code)) {
      const resp = await post("/get", {
        index: $current.index,
        code: $current.code,
        q: "chart",
      });
      const json = await resp.json();
      if (json.chart) data = json.chart;
    }
  };

  const updateChart = (force?: boolean) => {
    if (data.length && stock.now && (force || !hover)) {
      data[data.length - 1].y = stock.now;
      chart.data.datasets[0].data = data;
      chart.update();
    } else if (!data.length) {
      chart.data.datasets[0].data = [];
      chart.update();
    }
  };

  onMount(() => {
    start();
    return () => {
      for (; autoUpdate.length > 0; ) clearInterval(autoUpdate.pop());
      chart.destroy();
    };
  });
</script>

<svelte:head>
  <title>My Stocks</title>
</svelte:head>

<header>
  <AutoComplete />
  <div
    class="home"
    on:click={() => {
      window.history.pushState({}, "", "/");
      $component = "stocks";
    }}
  >
    <div class="icon"><i class="material-icons">home</i></div>
    <span>Home</span>
  </div>
  <Realtime bind:stock />
</header>
<div
  class="chart"
  on:mouseenter={() => (hover = true)}
  on:mouseleave={() => {
    hover = false;
    updateChart(true);
  }}
>
  <canvas id="stockChart" />
</div>

<style>
  .home {
    display: inline-flex;
    cursor: pointer;
  }

  .home > span {
    color: #1da1f2 !important;
    font-weight: bold;
    padding-right: 20px;
    padding-top: 15px;
    padding-bottom: 15px;
  }

  .home:hover {
    background: rgba(29, 161, 242, 0.1);
    border-radius: 9999px;
  }

  .home .material-icons {
    font-size: 36px;
    color: #1da1f2;
  }

  .icon {
    flex-direction: column;
    display: flex;
    justify-content: center;
    padding-left: 20px;
  }

  .chart {
    max-width: 1000px;
    max-height: 500px;
    height: calc(100% - 210px);
  }
</style>
