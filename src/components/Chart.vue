<template>
  <canvas class="chart" id="stockChart" />
</template>

<script>
import { checkTime, intraday } from "@/misc.js";
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
