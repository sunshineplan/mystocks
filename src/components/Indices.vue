<template>
  <div class="indices" v-if="Object.keys(indices).length !== 0">
    <div
      v-for="(val, key) in names"
      :key="key"
      :id="key"
      @click="gotoStock(indices[key])"
    >
      <span class="short">{{ key }}</span>
      <span class="full">{{ val }}</span>
      <span
        v-for="field in fields"
        :key="field"
        :style="addColor(indices[key], field)"
      >
        &nbsp;&nbsp;{{ indices[key][field] }}
      </span>
    </div>
  </div>
</template>

<script>
import { checkTime } from "@/misc.js";

export default {
  name: "Indices",
  data() {
    return {
      names: { 沪: "上证指数", 深: "深证成指", 创: "创业板指", 中: "中小板指" },
      fields: ["now", "change", "percent"],
    };
  },
  computed: {
    indices() {
      return this.$store.state.indices;
    },
  },
  created() {
    this.start();
  },
  methods: {
    start() {
      this.load(true);
      setInterval(this.load, 10000);
    },
    async load(force) {
      if (checkTime() || force) {
        this.$store.dispatch("indices");
      }
    },
  },
};
</script>

<style scoped>
.indices {
  position: fixed;
  z-index: 100;
  bottom: 0;
  width: 100%;
  height: 70px;
  display: flex;
  align-items: center;
  background-color: white;
  box-shadow: 0 -1px 2px 0 #e7e7e7;
  white-space: normal;
}

#沪,
#深,
#创,
#中 {
  color: black;
  max-width: 25%;
  flex: 0 0 25%;
  cursor: default;
  text-align: center;
  font-size: 20px;
}

#沪,
#深,
#创,
#中:hover {
  text-decoration: none;
}

.short {
  display: none;
}

@media (max-width: 1360px) {
  .short {
    display: inline;
  }

  .full {
    display: none;
  }
}
</style>
