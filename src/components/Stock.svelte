<script lang="ts">
  import type { ScatterDataPoint } from "chart.js";
  import { Chart } from "chart.js";
  import type { LineAnnotationOptions } from "chartjs-plugin-annotation";
  import { onMount } from "svelte";
  import { checkTradingTime, intraday, post } from "../misc";
  import { mystocks } from "../stock.svelte";
  import AutoComplete from "./AutoComplete.svelte";
  import Realtime from "./Realtime.svelte";

  let data: ScatterDataPoint[] = [];
  let chart: Chart<"line">;
  let hover = $state(false);
  let canvas: HTMLCanvasElement;
  let timer: number;
  let controller: AbortController;

  $effect(() => {
    mystocks.current.index;
    mystocks.current.code;
    subscribe();
    return abort;
  });

  $effect(() => {
    mystocks.current.update;
    updateChart();
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
    if (mystocks.current.stock.last)
      return `${
        Math.round(
          ((Number(value) - mystocks.current.stock.last) /
            mystocks.current.stock.last) *
            10000,
        ) / 100
      }%`;
    return null;
  };

  callbacks.label = (tooltipItem) => {
    const value = tooltipItem.parsed.y;
    const percent =
      Math.round(
        ((value - mystocks.current.stock.last) / mystocks.current.stock.last) *
          10000,
      ) / 100;
    return `${value}   ${percent}%`;
  };
  callbacks.labelTextColor = (tooltipItem) => {
    const change = tooltipItem.parsed.y - mystocks.current.stock.last;
    if (change > 0) return "red";
    else if (change < 0) return "green";
    return "black";
  };

  const updateChart = (force?: boolean) => {
    if (!chart) return;
    const yAxes = chart.options?.scales?.y;
    yAxes.suggestedMin = mystocks.current.stock.last / 1.01;
    yAxes.suggestedMax = mystocks.current.stock.last * 1.01;
    const annotations = chart.options.plugins?.annotation
      ?.annotations as Record<string, LineAnnotationOptions>;
    annotations.last.value = mystocks.current.stock.last;
    if (data.length && mystocks.current.stock.now && (force || !hover)) {
      data[data.length - 1].y = mystocks.current.stock.now;
      chart.data.datasets[0].data = data;
      chart.update();
    } else if (!data.length) {
      chart.data.datasets[0].data = [];
      chart.update();
    }
  };

  const subscribe = () => {
    controller = new AbortController();
    const fetchData = async (force?: boolean) => {
      let resp: Response;
      try {
        if (mystocks.current.code && (force || (await checkTradingTime())))
          resp = await fetch(
            `/chart?index=${mystocks.current.index}&code=${mystocks.current.code}`,
            { signal: controller.signal },
          );
        else resp = new Response(null, { status: 400 });
      } catch (e) {
        if (e instanceof DOMException && e.name === "AbortError") return;
        console.error(e);
        resp = new Response(null, { status: 500 });
      }
      let timeout = 30000;
      if (resp.ok) {
        const res = await resp.json();
        if (res.chart) {
          data = res.chart;
          updateChart(true);
        }
        timeout = 60000;
      } else if (resp.status == 400) timeout = 1000;
      timer = setTimeout(fetchData, timeout);
    };
    fetchData(true);
  };

  const abort = () => {
    if (timer) clearTimeout(timer);
    if (controller) controller.abort();
  };

  onMount(() => {
    chart = new Chart(canvas, intraday);
    return () => chart.destroy();
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
  <Realtime />
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
