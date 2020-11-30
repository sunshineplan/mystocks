<script context="module">
  export async function preload(page) {
    const { index, code } = page.params;
    return { index, code };
  }
</script>

<script>
  import AutoComplete from "../../components/AutoComplete.svelte";
  import Realtime from "../../components/Realtime.svelte";
  import StockChart from "../../components/Chart.svelte";
  import { onMount } from "svelte";
  import { goto } from "@sapper/app";
  import Chart from "chart.js";
  import annotation from "chartjs-plugin-annotation";
  import { stock, refresh } from "../../stores.js";
  import { checkTime, post, intraday } from "../../misc.js";

  export let index, code;

  Chart.defaults.global.maintainAspectRatio = false;
  Chart.defaults.global.legend.display = false;
  Chart.defaults.global.hover.mode = "index";
  Chart.defaults.global.hover.intersect = false;
  Chart.defaults.global.tooltips.mode = "index";
  Chart.defaults.global.tooltips.intersect = false;
  Chart.defaults.global.tooltips.displayColors = false;
  Chart.defaults.global.animation.duration = 0;
  Chart.plugins.register({ annotation });

  let autoUpdate = [];
  let data = [];
  let chart, update;

  async function start() {
    chart = new Chart(document.querySelector("#stockChart"), intraday);
    if (code != "n/a") {
      await load();
      autoUpdate.push(setInterval(loadRealtime, refresh * 1000));
      autoUpdate.push(setInterval(loadChart, 60000));
    }
  }

  async function load() {
    update = "";
    await loadRealtime(true);
    await loadChart(true);
    chart.options.scales.yAxes[0].ticks.suggestedMin = $stock.last / 1.01;
    chart.options.scales.yAxes[0].ticks.suggestedMax = $stock.last * 1.01;
    chart.annotation.elements.PreviousClose.options.value = $stock.last;
    updateChart();
  }

  async function loadRealtime(force) {
    if (checkTime() || (force && code)) {
      let response = await post("/get", {
        index: index,
        code: code,
        q: "realtime",
      });
      let json = await response.json();
      if (json.name) {
        stock.set(json);
        update = json.update;
        if (update) updateChart();
        document.title = `${json.name} ${json.now} ${json.percent}`;
      }
    }
  }

  async function loadChart(force) {
    if (checkTime() || (force && code)) {
      let resp = await post("/get", { index, code, q: "chart" });
      let json = await resp.json();
      if (json.chart) data = json.chart;
    }
  }

  function updateChart() {
    if (data.length && $stock.now) {
      data[data.length - 1].y = $stock.now;
      chart.data.datasets[0].data = data;
      chart.update();
    }
  }

  onMount(() => {
    start();
    return () => {
      for (; autoUpdate.length > 0; ) clearInterval(autoUpdate.pop());
      chart.destroy();
    };
  });
</script>

<style>
  .home {
    display: inline-flex;
    cursor: pointer;
  }

  .home a {
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
</style>

<svelte:head>
  <title>My Stocks</title>
</svelte:head>

<div class="content">
  <header>
    <AutoComplete />
    <div class="home" on:click={goto('/')}>
      <div class="icon"><i class="material-icons">home</i></div>
      <a>Home</a>
    </div>
    <Realtime />
  </header>
  <StockChart />
</div>
