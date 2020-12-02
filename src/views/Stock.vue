<template>
  <div class="content">
    <header>
      <AutoComplete />
      <div class="home" @click="$router.push('/')">
        <div class="icon">
          <i class="material-icons">home</i>
        </div>
        <a>Home</a>
      </div>
      <Realtime />
    </header>
    <StockChart />
  </div>
</template>

<script>
import AutoComplete from "@/components/AutoComplete.vue";
import Realtime from "@/components/Realtime.vue";
import StockChart from "@/components/Chart.vue";
import { checkTime, intraday } from "@/misc.js";
import Chart from "chart.js";
import annotation from "chartjs-plugin-annotation";
import Cookies from "js-cookie";

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
  name: "Stock",
  components: { AutoComplete, Realtime, StockChart },
  data() {
    return {
      refresh: Cookies.get("Refresh") ? Cookies.get("Refresh") : 3,
      autoUpdate: [],
      chart: "",
      update: "",
    };
  },
  computed: {
    index() {
      return this.$route.params.index;
    },
    code() {
      return this.$route.params.code;
    },
    stock() {
      return this.$store.state.stock;
    },
    data() {
      return this.$store.state.chart;
    },
  },
  watch: {
    update(now, last) {
      if (now && last) this.updateChart();
    },
    async $route(to) {
      if (to.name == "stock" && this.code != "n/a") await this.load();
    },
  },
  mounted() {
    document.title = "My Stocks";
    this.start();
  },
  beforeUnmount() {
    for (; this.autoUpdate.length > 0; ) clearInterval(this.autoUpdate.pop());
    this.chart.destroy();
  },
  methods: {
    async start() {
      this.chart = new Chart(document.querySelector("#stockChart"), intraday);
      if (this.code != "n/a") {
        await this.load();
        this.autoUpdate.push(
          setInterval(this.loadRealtime, this.refresh * 1000)
        );
        this.autoUpdate.push(setInterval(this.loadChart, 60000));
      }
    },
    async load() {
      this.update = "";
      await this.loadRealtime(true);
      await this.loadChart(true);
      this.chart.options.scales.yAxes[0].ticks.suggestedMin =
        this.stock.last / 1.01;
      this.chart.options.scales.yAxes[0].ticks.suggestedMax =
        this.stock.last * 1.01;
      this.chart.annotation.elements.PreviousClose.options.value = this.stock.last;
      this.updateChart();
    },
    async loadRealtime(force) {
      if (checkTime() || (force && this.code)) {
        await this.$store.dispatch("stock", {
          index: this.index,
          code: this.code,
        });
        if (this.stock.name) {
          this.update = this.stock.update;
          document.title = `${this.stock.name} ${this.stock.now} ${this.stock.percent}`;
        }
      }
    },
    async loadChart(force) {
      if (checkTime() || (force && this.code))
        await this.$store.dispatch("chart", {
          index: this.index,
          code: this.code,
        });
    },
    updateChart() {
      if (this.data.length && this.stock.now) {
        let data = this.data;
        data[data.length - 1].y = this.stock.now;
        this.chart.data.datasets[0].data = data;
        this.chart.update();
      }
    },
  },
};
</script>

<style scoped>
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
