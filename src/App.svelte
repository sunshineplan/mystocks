<script lang="ts">
  import type { Component } from "svelte";
  import Flows from "./components/Flows.svelte";
  import Indices from "./components/Indices.svelte";
  import Login from "./components/Login.svelte";
  import Nav from "./components/Nav.svelte";
  import Setting from "./components/Setting.svelte";
  import Stock from "./components/Stock.svelte";
  import Stocks from "./components/Stocks.svelte";
  import { mystocks } from "./stock.svelte";

  const components: {
    [component: string]: Component;
  } = {
    login: Login,
    setting: Setting,
    stock: Stock,
  };

  const Show = $derived(
    mystocks.component != "stocks"
      ? components[mystocks.component]
      : mystocks.isFlows
        ? Flows
        : Stocks,
  );

  const init = async () => {
    await mystocks.info();
    if (/^\/stock\/[A-Z]{3,4}\/\d{6}$/.test(window.location.pathname)) {
      const stock = window.location.pathname.split("/");
      mystocks.goto(stock[2], stock[3]);
    }
  };
  const promise = init();
</script>

<Nav />
{#await promise then _}
  <div class="content">
    <Show />
  </div>
{/await}
<Indices />

<style>
  .content {
    position: fixed;
    top: 0;
    padding: 60px 0 80px;
    width: 100%;
    height: 100%;
  }
</style>
