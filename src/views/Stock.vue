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
import { defineAsyncComponent } from "vue";

export default {
  name: "Stock",
  components: {
    AutoComplete: defineAsyncComponent(() =>
      import(
        /* webpackChunkName: "stock" */ "@/components/AutoComplete.vue"
      )
    ),
    Realtime: defineAsyncComponent(() =>
      import(/* webpackChunkName: "stock" */ "@/components/Realtime.vue")
    ),
    StockChart: defineAsyncComponent(() =>
      import(/* webpackChunkName: "chart" */ "@/components/Chart.vue")
    ),
  },
  mounted() {
    document.title = "My Stocks";
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
