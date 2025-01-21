<script lang="ts">
  import { Chart, type ChartDataset, type ScatterDataPoint } from "chart.js";
  import { onMount } from "svelte";
  import { capitalflows, checkTradingTime, dateStr, getColor } from "../misc";
  import AutoComplete from "./AutoComplete.svelte";

  interface Flows {
    sector: string;
    chart: ScatterDataPoint[];
  }

  let chart: Chart<"line">;
  let datasets: ChartDataset<"line">[] = [];
  let show = $state<number[]>([]);
  let current = $state(dateStr(new Date()));
  let today = $state(dateStr(new Date()));
  let last = $state("");
  let loading = $state(0);
  let status = $state(0);
  let hover = $state(false);
  let canvas: HTMLCanvasElement;
  let timer: number;
  let controller: AbortController;

  const gotoDate = (n: -1 | 0 | 1) => {
    let date: Date;
    if (!n) date = new Date();
    else {
      date = new Date(current);
      date.setDate(date.getDate() + n);
    }
    const newDate = dateStr(date);
    if (current == newDate) return;
    current = newDate;
    reset();
  };

  const reset = () => {
    show.length = 0;
    updateChart(true);
    if (current != today) last = "";
    abort();
    subscribe();
  };

  const legend = capitalflows.options?.plugins?.legend;

  legend.onClick = (_, legendItem) => {
    const index = legendItem.datasetIndex;
    if (typeof index === "number") {
      display(index);
      chart.update();
    }
  };

  const display = (index: number) => {
    const datasets = chart.data.datasets;
    if (!show.length) {
      datasets.forEach((_, i) => {
        const meta = chart.getDatasetMeta(i);
        if (i !== index) meta.hidden = true;
      });
      show.push(index);
    } else if (show.includes(index) && show.length == 1) {
      datasets.forEach((_, i) => {
        const meta = chart.getDatasetMeta(i);
        meta.hidden = false;
      });
      show.length = 0;
    } else if (show.includes(index) && show.length > 1) {
      chart.getDatasetMeta(index).hidden = true;
      show.splice(show.indexOf(index), 1);
    } else {
      datasets.forEach((_, i) => {
        const meta = chart.getDatasetMeta(i);
        if (i == index) meta.hidden = false;
      });
      show.push(index);
    }
  };

  const updateChart = (empty?: boolean) => {
    if (!chart) return;
    if (empty) datasets.length = 0;
    chart.data.datasets = [...datasets];
    if (show.length)
      chart.data.datasets.forEach((_, i) => {
        const meta = chart.getDatasetMeta(i);
        if (show.includes(i)) meta.hidden = false;
        else meta.hidden = true;
      });
    if (!hover) chart.update();
    if (empty) return;
    if (current == today) last = new Date().toLocaleString();
  };

  const subscribe = () => {
    controller = new AbortController();
    const fetchData = async (force?: boolean) => {
      const now = dateStr(new Date());
      if (today != now) today = now;
      let url = "/flows";
      if (current != today) {
        if (!force) {
          timer = setTimeout(fetchData, 1000);
          return;
        } else if (new Date(current) > new Date()) {
          status = -1;
          return;
        }
        url += `?date=${current}`;
      }
      if (force) updateChart(true);
      let resp: Response;
      let array: Flows[];
      loading++;
      try {
        if (force || (await checkTradingTime()))
          resp = await fetch(url, { signal: controller.signal });
        else resp = new Response(null, { status: 400 });
      } catch (e) {
        if (e instanceof DOMException && e.name === "AbortError") return;
        console.error(e);
        resp = new Response(null, { status: 500 });
        status = 0;
        datasets.length = 0;
      }
      let timeout = 30000;
      if (resp.ok) {
        array = await resp.json();
        if (array && array.length) {
          status = 1;
          datasets.length = 0;
          array.forEach((e: Flows, i: number) => {
            datasets.push({
              label: e.sector,
              fill: false,
              tension: 0,
              borderWidth: 1.5,
              borderColor: getColor(i),
              backgroundColor: getColor(i),
              pointRadius: 0,
              pointHoverRadius: 3,
              data: e.chart.map((i) => i.y / 100000000),
            });
          });
        } else status = -1;
        timeout = 60000;
      } else if (resp.status == 400) timeout = 1000;
      loading--;
      updateChart();
      timer = setTimeout(fetchData, timeout);
    };
    fetchData(true);
  };

  const abort = () => {
    if (timer) clearTimeout(timer);
    if (controller) controller.abort();
  };

  onMount(() => {
    chart = new Chart(canvas, capitalflows);
    subscribe();
    return () => {
      abort();
      chart.destroy();
    };
  });
</script>

<svelte:head>
  <title>My Stocks</title>
</svelte:head>

<header style="height:80px">
  <AutoComplete />
  <div>
    <div class="input-group">
      <button
        class="input-group-text"
        disabled={loading ? true : false}
        onclick={() => gotoDate(-1)}
      >
        -
      </button>
      <input
        class="form-control"
        type="date"
        disabled={loading ? true : false}
        bind:value={current}
        onchange={reset}
      />
      <button
        class="input-group-text"
        disabled={loading ? true : false}
        onclick={() => gotoDate(1)}
      >
        +
      </button>
    </div>
    <button
      class="btn btn-danger"
      onclick={() => {
        if (chart.data.datasets && current == today)
          chart.data.datasets.forEach((_, i) => {
            const meta = chart.getDatasetMeta(i);
            meta.hidden = false;
          });
        else chart.data.datasets = [];
        show.length = 0;
        chart.update();
        gotoDate(0);
      }}
    >
      Reset
    </button>
    {#if loading}
      <div class="spinner-border text-secondary" role="status">
        <span class="visually-hidden">Loading...</span>
      </div>
    {:else if status}
      {#if status == 1}
        <i class="material-icons text-success">done</i>
      {:else}
        <i class="material-icons text-warning" title="No data of this date">
          warning_amber
        </i>
      {/if}
    {:else}
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <i
        class="material-icons text-danger"
        onclick={() => {
          abort();
          subscribe();
        }}
      >
        close
      </i>
    {/if}
  </div>
  {#if current == today && last}
    <small>Last update: {last}</small>
  {/if}
</header>
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  class="chart"
  onmouseenter={() => (hover = true)}
  onmouseleave={() => {
    hover = false;
    if (current == today) {
      chart.data.datasets = [...datasets];
      updateChart();
    }
  }}
>
  <canvas bind:this={canvas}></canvas>
</div>

<style>
  .input-group {
    display: inline-flex;
    width: 240px;
    vertical-align: middle;
  }

  .input-group-text {
    width: 35px;
    justify-content: center;
  }

  .spinner-border {
    height: 32px;
    width: 32px;
    vertical-align: middle;
  }

  .material-icons {
    vertical-align: middle;
    font-size: 36px;
    cursor: default;
  }

  .chart {
    height: calc(100% - 80px);
  }
</style>
