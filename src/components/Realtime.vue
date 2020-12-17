<template>
  <div>
    <div style="display: flex; font-size: 2rem">
      <i
        class="material-icons star"
        :class="stared ? 'stared' : ''"
        @click="star()"
      >
        {{ stared ? "star" : "star_border" }}
      </i>
      <span>
        {{ stock.name }}
      </span>
      (
      <span>{{ stock.code }}</span>
      )
      <i class="material-icons open" @click="open">open_in_new</i>
      &nbsp;&nbsp;&nbsp;
      <span :style="addColor(stock, 'now')">
        {{ stock.now }}
      </span>
      &nbsp;&nbsp;&nbsp;
      <span :style="addColor(stock, 'percent')">{{ stock.percent }}</span>
    </div>
    <div style="min-height: 52px">
      <table style="float: left; table-layout: fixed" :style="{ width: width }">
        <tbody>
          <tr>
            <td>
              昨收:
              <span>{{ stock.last }}</span>
            </td>
            <td>
              涨跌:
              <span :style="addColor(stock, 'change')">{{ stock.change }}</span>
            </td>
            <td>
              涨幅:
              <span :style="addColor(stock, 'percent')">
                {{ stock.percent }}
              </span>
            </td>
          </tr>
          <tr>
            <td>
              最高:
              <span :style="addColor(stock, 'high')">{{ stock.high }}</span>
            </td>
            <td>
              最低:
              <span :style="addColor(stock, 'low')">{{ stock.low }}</span>
            </td>
            <td>
              开盘:
              <span :style="addColor(stock, 'open')">{{ stock.open }}</span>
            </td>
          </tr>
        </tbody>
      </table>
      <table v-if="stock.sell5 || stock.buy5">
        <tbody>
          <tr>
            <td>
              <span style="display: inline-flex">
                卖盘:&nbsp;
                <div
                  class="sellbuy"
                  style="color: red"
                  v-for="(sell, index) in stock.sell5"
                  :key="index"
                >
                  {{ sell.Price }}-{{ sell.Volume }}
                </div>
              </span>
            </td>
          </tr>
          <tr>
            <td>
              <span style="display: inline-flex">
                买盘:&nbsp;
                <div
                  class="sellbuy"
                  style="color: green"
                  v-for="(buy, index) in stock.buy5"
                  :key="index"
                >
                  {{ buy.Price }}-{{ buy.Volume }}
                </div>
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
    <small>
      更新时间: <span class="update">{{ stock.update }}</span>
    </small>
  </div>
</template>

<script>
import { checkTime, post } from "@/misc.js";
import Cookies from "js-cookie";

export default {
  name: "Realtime",
  data() {
    return {
      refresh: Cookies.get("Refresh") || 3,
      stared: false,
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
    stock() {
      return this.$store.state.stock;
    },
    width() {
      return !this.stock.sell5 && !this.stock.buy5 ? "480px" : "360px";
    },
  },
  watch: {
    async $route(to) {
      this.$store.commit("stock", {});
      if (to.name == "stock" && this.code != "n/a") await this.load(true);
    },
  },
  async created() {
    const resp = await fetch("/star");
    if ((await resp.text()) == "1") this.stared = true;
  },
  async mounted() {
    await this.start();
  },
  beforeUnmount() {
    clearInterval(this.autoUpdate);
  },
  methods: {
    async star() {
      let resp;
      if (!this.stared) resp = await post("/star", {});
      else resp = await post("/star", { action: "unstar" });
      if ((await resp.text()) == "1") this.stared = !this.stared;
    },
    async start() {
      if (this.code != "n/a") {
        await this.load(true);
        this.autoUpdate = setInterval(this.load, this.refresh * 1000);
      }
    },
    async load(force) {
      if (checkTime() || (force && this.code)) {
        await this.$store.dispatch("stock", {
          index: this.index,
          code: this.code,
        });
      }
    },
    open() {
      window.open("http://stockpage.10jqka.com.cn/" + this.stock.code);
    },
  },
};
</script>

<style scoped>
.star {
  color: #f4b400;
  width: 50px;
  height: 50px;
  font-size: 40px;
  cursor: default;
  display: flex;
  justify-content: center;
  align-items: center;
}

.star:hover {
  background-color: #fcebbaa1;
  border-radius: 50%;
}

.open {
  margin-left: 0.5rem;
  color: gray;
  cursor: pointer;
  display: flex;
  align-items: center;
}

.sellbuy {
  min-width: 115px;
  padding-right: 6px;
}
</style>
