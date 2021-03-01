<script lang="ts">
  import { onMount } from "svelte";
  import Chart from "chart.js";
  import AutoComplete from "./AutoComplete.svelte";
  import { checkTime, getColor, capitalflows } from "../misc";
  import type { Flows } from "../stores";

  let autoUpdate = 0;
  let chart: Chart;
  let show: number[] = [];
  let date = "";
  let last = "";
  let loading = 0;
  let status = 0;

  $: date && goto();

  const getDate = (n: -1 | 0 | 1) => {
    let day: Date;
    if (n == 0) day = new Date();
    else day = new Date(date);
    day.setDate(day.getDate() + n);
    const dd = String(day.getDate()).padStart(2, "0");
    const mm = String(day.getMonth() + 1).padStart(2, "0");
    const yyyy = day.getFullYear();
    date = `${yyyy}-${mm}-${dd}`;
    return date;
  };
  const today = getDate(0);

  ((capitalflows.options as Chart.ChartOptions)
    .legend as Chart.ChartLegendOptions).onClick = (
    event: MouseEvent,
    legendItem: Chart.ChartLegendLabelItem
  ) => {
    display(legendItem.datasetIndex as number);
    chart.update();
  };

  const display = (index: number) => {
    if (!show.length) {
      (chart.data.datasets as Chart.ChartDataSets[]).forEach((e, i) => {
        const meta = chart.getDatasetMeta(i);
        if (i !== index) meta.hidden = true;
      });
      show.push(index);
    } else if (show.includes(index) && show.length == 1) {
      (chart.data.datasets as Chart.ChartDataSets[]).forEach((e, i) => {
        const meta = chart.getDatasetMeta(i);
        meta.hidden = undefined;
      });
      show.length = 0;
    } else if (show.includes(index) && show.length > 1) {
      chart.getDatasetMeta(index).hidden = true;
      show.splice(show.indexOf(index), 1);
    } else {
      (chart.data.datasets as Chart.ChartDataSets[]).forEach((e, i) => {
        const meta = chart.getDatasetMeta(i);
        if (i == index) meta.hidden = false;
      });
      show.push(index);
    }
  };

  const load = async (force?: boolean, date?: string) => {
    let url = "/flows";
    if (date) {
      if (new Date(date) > new Date()) {
        status = -1;
        return;
      }
      url = url + `?date=${date}`;
    }
    const datasets = chart.data.datasets as Chart.ChartDataSets[];
    if (force) {
      datasets.length = 0;
      chart.update();
    }
    if (checkTime() || force) {
      loading++;
      const resp = await fetch(url);
      if (!resp.ok) {
        status = 0;
        loading--;
        return;
      }
      const array = await resp.json();
      if (array && array.length) {
        status = 1;
        datasets.length = 0;
        array.forEach((e: Flows, i: number) => {
          datasets.push({
            label: e.sector,
            fill: false,
            lineTension: 0,
            borderWidth: 1.5,
            borderColor: getColor(i),
            backgroundColor: getColor(i),
            pointRadius: 0,
            pointHoverRadius: 3,
            data: e.chart.map((i) => (i.y as number) / 100000000),
          });
        });
        if (show.length)
          datasets.forEach((e, i) => {
            const meta = chart.getDatasetMeta(i);
            if (show.includes(i)) meta.hidden = false;
            else meta.hidden = true;
          });
        chart.update();
        if (!date || (date && date == today))
          last = new Date().toLocaleString();
      } else status = -1;
      loading--;
    }
  };

  const goto = () => {
    if (!chart) return;
    show.length = 0;
    (chart.data.datasets as Chart.ChartDataSets[]).length = 0;
    chart.update();
    if (date != today) {
      if (autoUpdate) clearInterval(autoUpdate);
      load(true, date);
    } else {
      last = "";
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
  <div>
    <div class="input-group">
      <div class="input-group-prepend">
        <button
          class="input-group-text"
          disabled={loading ? true : false}
          on:click={() => getDate(-1)}
        >
          -
        </button>
      </div>
      <input
        class="form-control"
        type="date"
        disabled={loading ? true : false}
        bind:value={date}
      />
      <div class="input-group-append">
        <button
          class="input-group-text"
          disabled={loading ? true : false}
          on:click={() => getDate(1)}
        >
          +
        </button>
      </div>
    </div>
    <button
      class="btn btn-danger"
      on:click={() => {
        getDate(0);
        show.length = 0;
      }}
    >
      Reset
    </button>
    {#if loading}
      <div class="spinner-border text-secondary" role="status">
        <span class="sr-only">Loading...</span>
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
      <i class="material-icons text-danger" on:click={() => load(true, date)}>
        close
      </i>
    {/if}
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
    display: inline-flex;
    width: 240px;
    vertical-align: middle;
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
