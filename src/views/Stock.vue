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
      <Realtime :stock="stock" />
    </header>
    <StockChart @mouseenter="enter" @mouseleave="leave" />
  </div>
</template>

<script>
import { defineAsyncComponent } from "vue";
import Cookies from "js-cookie";
import { checkTime } from "@/misc.js";

export default {
  name: "Stock",
  components: {
    AutoComplete: defineAsyncComponent(() =>
      import(/* webpackChunkName: "stock" */ "@/components/AutoComplete.vue")
    ),
    Realtime: defineAsyncComponent(() =>
      import(/* webpackChunkName: "stock" */ "@/components/Realtime.vue")
    ),
    StockChart: defineAsyncComponent(() =>
      import(/* webpackChunkName: "chart" */ "@/components/Chart.vue")
    ),
  },
  data() {
    return {
      refresh: Cookies.get("Refresh") || 3,
      autoUpdate: [],
      update: "",
      hover: false,
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
  },
  watch: {
    update(now, last) {
      if (now && last && !this.hover) this.$store.dispatch("updateChart");
    },
    async $route(to) {
      if (to.name == "stock" && this.code != "n/a") await this.reload();
    },
  },
  async mounted() {
    document.title = "My Stocks";
    if (this.code != "n/a") {
      await this.reload();
      this.autoUpdate.push(setInterval(this.loadRealtime, this.refresh * 1000));
      this.autoUpdate.push(setInterval(this.loadChart, 60000));
    }
  },
  beforeUnmount() {
    for (; this.autoUpdate.length > 0; ) clearInterval(this.autoUpdate.pop());
    this.$store.dispatch("destroyChart");
  },
  methods: {
    enter() {
      this.hover = true;
    },
    leave() {
      setTimeout(() => (this.hover = false), 200);
    },
    async reload() {
      this.$store.dispatch("resetChart");
      await this.loadRealtime(true);
      await this.loadChart(true);
      this.$store.dispatch("updateChart");
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
        await this.$store.dispatch("line", {
          index: this.index,
          code: this.code,
        });
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
