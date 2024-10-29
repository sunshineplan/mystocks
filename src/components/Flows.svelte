<script lang="ts">
  import { onMount } from "svelte";
  import { Chart, type ChartDataset, type ScatterDataPoint } from "chart.js";
  import AutoComplete from "./AutoComplete.svelte";
  import { dateStr, checkTradingTime, getColor, capitalflows } from "../misc";

  interface Flows {
    sector: string;
    chart: ScatterDataPoint[];
  }

  let autoUpdate: number;
  let chart: Chart<"line">;
  let datasets: ChartDataset<"line">[] = [];
  let show = $state<number[]>([]);
  let date = $state(getDate(0));
  let today = $state(getDate(0));
  let last = $state("");
  let loading = $state(0);
  let status = $state(0);
  let controller: AbortController;
  let hover = $state(false);
  let dayChange = false;

  function getDate(n: -1 | 0 | 1, setDate?: boolean) {
    let day: Date;
    if (n == 0) day = new Date();
    else {
      day = new Date(date);
      day.setDate(day.getDate() + n);
    }
    const ymd = dateStr(day);
    if (setDate) {
      date = ymd;
      goto(date);
    }
    return ymd;
  }

  const updateDate = () => {
    if (date == today && date != getDate(0)) {
      dayChange = true;
      today = getDate(0, true);
    } else {
      today = getDate(0);
    }
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

  const load = async (force?: boolean, date?: string) => {
    let url = "/flows";
    if (date && date != today) {
      if (new Date(date) > new Date()) {
        status = -1;
        return;
      }
      url = url + `?date=${date}`;
    }
    if (force) updateChart(true);
    if (force || (await checkTradingTime())) {
      updateDate();
      loading++;
      let array: Flows[];
      try {
        controller = new AbortController();
        setTimeout(() => controller.abort("fetch flows timeout"), 50000);
        const resp = await fetch(url, { signal: controller.signal });
        if (!resp.ok) {
          status = 0;
          datasets.length = 0;
          chart.data.datasets.length = 0;
          loading--;
          return;
        }
        array = await resp.json();
      } catch (e) {
        console.log(e);
        status = 0;
        datasets.length = 0;
        chart.data.datasets.length = 0;
        loading--;
        return;
      }
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
        if (force || !hover) {
          chart.data.datasets = [...datasets];
          updateChart();
        }
      } else status = -1;
      loading--;
    }
  };

  const goto = (date: string) => {
    if (dayChange) {
      dayChange = false;
      return;
    }
    if (!chart) return;
    if (controller) controller.abort();
    show.length = 0;
    updateChart(true);
    if (date != today) {
      if (autoUpdate) clearInterval(autoUpdate);
      load(true, date);
    } else {
      last = "";
      load(true);
      autoUpdate = setInterval(load, 60000);
    }
  };

  const updateChart = (empty?: boolean) => {
    const datasets = chart.data.datasets;
    if (empty) {
      datasets.length = 0;
      chart.update();
      return;
    }
    if (show.length)
      datasets.forEach((_, i) => {
        const meta = chart.getDatasetMeta(i);
        if (show.includes(i)) meta.hidden = false;
        else meta.hidden = true;
      });
    chart.update();
    if (!date || (date && date == today)) last = new Date().toLocaleString();
  };

  onMount(() => {
    chart = new Chart(
      document.querySelector<HTMLCanvasElement>("#flowsChart"),
      capitalflows,
    );
    goto(date);
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
  <div>
    <div class="input-group">
      <button
        class="input-group-text"
        disabled={loading ? true : false}
        onclick={() => getDate(-1, true)}
      >
        -
      </button>
      <input
        class="form-control"
        type="date"
        disabled={loading ? true : false}
        bind:value={date}
        onchange={() => goto(date)}
      />
      <button
        class="input-group-text"
        disabled={loading ? true : false}
        onclick={() => getDate(1, true)}
      >
        +
      </button>
    </div>
    <button
      class="btn btn-danger"
      onclick={() => {
        if (chart.data.datasets && date == today)
          chart.data.datasets.forEach((_, i) => {
            const meta = chart.getDatasetMeta(i);
            meta.hidden = false;
          });
        else chart.data.datasets = [];
        show.length = 0;
        chart.update();
        getDate(0, true);
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
      <i class="material-icons text-danger" onclick={() => load(true, date)}>
        close
      </i>
    {/if}
  </div>
  {#if date == today && last}
    <small>Last update: {last}</small>
  {/if}
</header>
<!-- svelte-ignore a11y_no_static_element_interactions -->
<div
  class="chart"
  onmouseenter={() => (hover = true)}
  onmouseleave={() => {
    hover = false;
    if (date == today) {
      chart.data.datasets = [...datasets];
      updateChart();
    }
  }}
>
  <canvas id="flowsChart"></canvas>
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
