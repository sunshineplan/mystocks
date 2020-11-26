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
      indices: {},
      names: { 沪: "上证指数", 深: "深证成指", 创: "创业板指", 中: "中小板指" },
      fields: ["now", "change", "percent"],
    };
  },
  created() {
    this.start();
  },
  methods: {
    start() {
      this.load(true);
      setInterval(this.load, 10000);
    },
    load(force) {
      if (checkTime() || force) {
        fetch("/indices")
          .then((response) => response.json())
          .then((json) => {
            this.indices = json;
          });
      }
    },
  },
};
</script>
