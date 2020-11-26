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
import Indices from "@/components/Indices.vue";

export default {
  name: "App",
  data() {
    return { user: document.querySelector("#app").dataset.user };
  },
  components: { Indices },
};
</script>
