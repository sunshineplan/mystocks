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
import { computed } from "vue";
import { checkTime, post, intraday } from "@/misc.js";
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

const refresh = Cookies.get("Refresh") ? Cookies.get("Refresh") : 3;

export default {
  name: "Stock",
  components: { AutoComplete, Realtime, StockChart },
  provide() {
    return { Stock: computed(() => this.stock) };
  },
  data() {
    return {
      refresh: Number(refresh) + 1,
      autoUpdate: [],
      stock: {},
      chart: "",
      data: [],
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
        let response = await post("/get", {
          index: this.index,
          code: this.code,
          q: "realtime",
        });
        let stock = await response.json();
        if (stock.name) {
          this.stock = stock;
          this.update = stock.update;
          document.title = `${stock.name} ${stock.now} ${stock.percent}`;
        }
      }
    },
    async loadChart(force) {
      if (checkTime() || (force && this.code)) {
        let response = await post("/get", {
          index: this.index,
          code: this.code,
          q: "chart",
        });
        let json = await response.json();
        if (json.chart) this.data = json.chart;
      }
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
