<script lang="ts">
  import Nav from "./components/Nav.svelte";
  import Login from "./components/Login.svelte";
  import Setting from "./components/Setting.svelte";
  import Stocks from "./components/Stocks.svelte";
  import Flows from "./components/Flows.svelte";
  import Stock from "./components/Stock.svelte";
  import Indices from "./components/Indices.svelte";
  import { username, current, component, info, flows } from "./stores";

  const components: {
    [component: string]: typeof Login | typeof Setting | typeof Stock;
  } = {
    login: Login,
    setting: Setting,
    stock: Stock,
  };

  const init = async () => {
    await info();
    if (/^\/stock\/[A-Z]{3,4}\/\d{6}$/.test(window.location.pathname)) {
      const stock = window.location.pathname.split("/");
      $current = { index: stock[2], code: stock[3] };
      window.history.pushState({}, "", `/stock/${stock[2]}/${stock[3]}`);
      $component = "stock";
    }
  };
  const promise = init();
</script>

<Nav bind:user={$username} />
{#await promise then _}
  <div class="content">
    <svelte:component
      this={$component != "stocks"
        ? components[$component]
        : $flows
          ? Flows
          : Stocks}
      on:info={info}
    />
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

  :global(:root) {
    --sk-color: #1a73e8;
  }

  :global(body) {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
      "Helvetica Neue", Arial, "Noto Sans", "Microsoft YaHei New",
      "Microsoft Yahei", 微软雅黑, 宋体, SimSun, STXihei, 华文细黑, sans-serif;
    white-space: nowrap;
  }

  :global(header) {
    padding: 10px 20px;
  }

  :global(small) {
    color: gray;
    padding-left: 5px;
  }

  :global(button + button) {
    margin-left: 0.3em;
  }

  :global(.swal) {
    margin: 8px 6px;
  }

  @media (max-width: 900px) {
    :global(.content) {
      padding-left: 0 !important;
    }
  }
</style>
