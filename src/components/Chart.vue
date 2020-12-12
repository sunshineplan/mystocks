<template>
  <canvas class="chart" id="stockChart" />
</template>

<script>
import { checkTime, timeLabels } from "@/misc.js";
import Chart from "chart.js";
import annotation from "chartjs-plugin-annotation";

Chart.defaults.global.maintainAspectRatio = false;
Chart.defaults.global.legend.display = false;
Chart.defaults.global.hover.mode = "index";
Chart.defaults.global.hover.intersect = false;
Chart.defaults.global.tooltips.mode = "index";
Chart.defaults.global.tooltips.intersect = false;
Chart.defaults.global.tooltips.displayColors = false;
Chart.defaults.global.animation.duration = 0;
Chart.plugins.register({ annotation });

const intraday = {
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

export default {
  name: "StockChart",
  data() {
    return {
      autoUpdate: "",
    };
  },
  computed: {
    index() {
      return this.$route.params.index;
    },
    code() {
      return this.$route.params.code;
    },
  },
  watch: {
    async $route(to) {
      this.$store.dispatch("resetChart");
      if (to.name == "stock" && this.code != "n/a") await this.load(true);
    },
  },
  async mounted() {
    await this.start();
  },
  beforeUnmount() {
    clearInterval(this.autoUpdate);
    this.$store.dispatch("destroyChart");
  },
  methods: {
    async start() {
      this.$store.commit(
        "chart",
        new Chart(document.querySelector("#stockChart"), intraday)
      );
      if (this.code != "n/a") {
        await this.load(true);
        this.autoUpdate = setInterval(this.load, 60000);
      }
    },
    async load(force) {
      if (checkTime() || (force && this.code))
        await this.$store.dispatch("line", {
          index: this.index,
          code: this.code,
        });
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
