<script lang="ts">
  import type { ScatterDataPoint } from "chart.js";
  import { Chart } from "chart.js";
  import type { LineAnnotationOptions } from "chartjs-plugin-annotation";
  import { onMount } from "svelte";
  import { checkTradingTime, intraday, post } from "../misc";
  import { mystocks } from "../stock.svelte";
  import AutoComplete from "./AutoComplete.svelte";
  import Realtime from "./Realtime.svelte";

  let autoUpdate: number[] = [];
  let stock: Stock = $state({
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
  });
  let data: ScatterDataPoint[] = [];
  let chart: Chart<"line">;
  let update = "";
  let hover = $state(false);
  let canvas: HTMLCanvasElement;

  $effect(() => {
    load(mystocks.current);
  });

  const y2 = intraday.options?.scales?.y2;
  const callbacks = intraday.options?.plugins?.tooltip?.callbacks;

  y2.afterBuildTicks = (axis) => {
    if (chart) {
      axis.ticks = chart.scales.y.ticks.map(({ value }) => {
        return { value };
      });
      axis.max = chart.scales.y.max;
      axis.min = chart.scales.y.min;
    }
  };

  y2.ticks.callback = (value) => {
    if (stock.last)
      return `${
        Math.round(((Number(value) - stock.last) / stock.last) * 10000) / 100
      }%`;
    return null;
  };

  callbacks.label = (tooltipItem) => {
    const value = tooltipItem.parsed.y;
    const percent =
      Math.round(((value - stock.last) / stock.last) * 10000) / 100;
    return `${value}   ${percent}%`;
  };
  callbacks.labelTextColor = (tooltipItem) => {
    const change = tooltipItem.parsed.y - stock.last;
    if (change > 0) return "red";
    else if (change < 0) return "green";
    return "black";
  };

  const start = () => {
    chart = new Chart(canvas, intraday);
    if (mystocks.current.code) {
      autoUpdate.push(setInterval(loadRealtime, mystocks.refresh * 1000));
      autoUpdate.push(setInterval(loadChart, 60000));
    }
  };

  const load = async (stock: any) => {
    update = "";
    await loadRealtime(true);
    await loadChart(true);
    const yAxes = chart.options?.scales?.y;
    yAxes.suggestedMin = stock.last / 1.01;
    yAxes.suggestedMax = stock.last * 1.01;
    const annotations = chart.options.plugins?.annotation
      ?.annotations as Record<string, LineAnnotationOptions>;
    annotations.last.value = stock.last;
    updateChart(true);
  };

  const loadRealtime = async (force?: boolean) => {
    if ((force && mystocks.current.code) || (await checkTradingTime())) {
      const resp = await post("/realtime", {
        index: mystocks.current.index,
        code: mystocks.current.code,
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
    if ((force && mystocks.current.code) || (await checkTradingTime())) {
      const resp = await post("/chart", {
        index: mystocks.current.index,
        code: mystocks.current.code,
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
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div
    class="home"
    onclick={() => {
      window.history.pushState({}, "", "/");
      mystocks.component = "stocks";
    }}
  >
    <div class="icon"><i class="material-icons">home</i></div>
    <span>Home</span>
  </div>
  <Realtime bind:stock />
</header>
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  class="chart"
  onmouseenter={() => (hover = true)}
  onmouseleave={() => {
    hover = false;
    updateChart(true);
  }}
>
  <canvas bind:this={canvas}></canvas>
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
