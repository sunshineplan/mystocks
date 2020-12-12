<template>
  <div class="content">
    <header>
      <AutoComplete />
    </header>
    <div class="table-responsive">
      <table class="table table-hover table-sm">
        <thead>
          <tr>
            <th v-for="(val, key) in columns" :key="key">{{ key }}</th>
          </tr>
        </thead>
        <tbody id="sortable">
          <tr
            v-for="stock in stocks"
            :key="stock.index + stock.code"
            @click="gotoStock(stock)"
          >
            <td
              v-for="(val, key) in columns"
              :key="key"
              :style="addColor(stock, val)"
            >
              {{ stock[val] }}
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<script>
import Sortable from "sortablejs";
import { defineAsyncComponent } from "vue";
import { checkTime, post } from "@/misc.js";
import Cookies from "js-cookie";

export default {
  name: "Stocks",
  components: {
    AutoComplete: defineAsyncComponent(() =>
      import(
        /* webpackChunkName: "stock" */ "@/components/AutoComplete.vue"
      )
    ),
  },
  data() {
    return {
      columns: {
        指数: "index",
        代码: "code",
        名称: "name",
        最新: "now",
        涨跌: "change",
        涨幅: "percent",
        最高: "high",
        最低: "low",
        开盘: "open",
        昨收: "last",
      },
      sortable: "",
      refresh: Cookies.get("Refresh") || 3,
      autoUpdate: "",
      fetching: "",
    };
  },
  computed: {
    stocks() {
      return this.$store.state.stocks;
    },
  },
  created() {
    this.start();
  },
  mounted() {
    document.title = "My Stocks";
    this.sortable = new Sortable(document.querySelector("#sortable"), {
      animation: 150,
      delay: 500,
      swapThreshold: 0.5,
      onStart: () => this.stop(),
      onEnd: () => this.start(),
      onUpdate: this.onUpdate,
    });
  },
  beforeUnmount() {
    this.stop();
    this.sortable.destroy();
  },
  methods: {
    start() {
      this.load(true);
      this.autoUpdate = setInterval(this.load, this.refresh * 1000);
    },
    stop() {
      this.fetching.abort();
      clearInterval(this.autoUpdate);
    },
    async load(force) {
      if (checkTime() || force) {
        this.fetching = new AbortController();
        const resp = await fetch("/mystocks", { signal: this.fetching.signal });
        this.$store.commit("stocks", await resp.json());
      }
    },
    onUpdate(evt) {
      post("/reorder", {
        old: `${this.stocks[evt.oldIndex].index} ${
          this.stocks[evt.oldIndex].code
        }`,
        new: `${this.stocks[evt.newIndex].index} ${
          this.stocks[evt.newIndex].code
        }`,
      });
    },
  },
};
</script>

<style scoped>
.table-responsive {
  height: calc(100% - 54px);
  padding: 0px 30px;
  cursor: default;
}

th {
  position: sticky;
  top: 0;
  border-top: 0 !important;
  border-bottom: 0 !important;
  background-color: white;
}

.sortable-ghost {
  opacity: 0;
}
</style>
