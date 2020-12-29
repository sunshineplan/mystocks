<template>
  <nav class="navbar navbar-light topbar">
    <router-link class="brand" to="/">My Stocks</router-link>
    <div class="navbar-nav flex-row">
      <div class="navbar-nav flex-row" v-if="user">
        <a class="nav-link" v-text="user"></a>
        <router-link class="nav-link link" to="/setting">Setting</router-link>
        <a class="nav-link link" href="/logout">Log out</a>
      </div>
      <div class="navbar-nav flex-row" v-else>
        <router-link class="nav-link link" to="/login">Log in</router-link>
      </div>
    </div>
  </nav>
  <router-view />
  <Indices />
</template>

<script>
import { defineAsyncComponent } from "vue";

export default {
  name: "App",
  components: {
    Indices: defineAsyncComponent(() =>
      import(/* webpackChunkName: "stock" */ "@/components/Indices.vue")
    ),
  },
};
</script>

<style scoped>
.topbar {
  position: fixed;
  top: 0;
  z-index: 2;
  width: 100%;
  height: 60px;
  background-color: #1a73e8;
}

.brand {
  padding-left: 20px;
  font-size: 25px;
  letter-spacing: 0.3px;
  color: white;
}

.brand:hover {
  color: white;
  text-decoration: none;
}

.topbar .nav-link {
  padding-left: 8px;
  padding-right: 8px;
  color: white !important;
}

.topbar .link:hover {
  background: rgba(255, 255, 255, 0.2);
  border-radius: 5px;
}
</style>

<style>
body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
    "Helvetica Neue", Arial, "Noto Sans", "Microsoft YaHei New",
    "Microsoft Yahei", 微软雅黑, 宋体, SimSun, STXihei, 华文细黑, sans-serif;
  white-space: nowrap;
}

header {
  padding: 10px 20px;
}

small {
  color: gray;
  padding-left: 5px;
}

button + button {
  margin-left: 0.3em;
}

.content {
  position: fixed;
  top: 0;
  padding: 60px 0 80px;
  width: 100%;
  height: 100%;
}

.icon {
  flex-direction: column;
  display: flex;
  justify-content: center;
  padding-left: 20px;
}
</style>
