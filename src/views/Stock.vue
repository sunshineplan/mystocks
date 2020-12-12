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

export default {
  name: "Stock",
  components: { AutoComplete, Realtime, StockChart },
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
