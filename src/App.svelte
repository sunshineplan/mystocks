<script lang="ts">
  import Nav from "./components/Nav.svelte";
  import Login from "./components/Login.svelte";
  import Setting from "./components/Setting.svelte";
  import Stocks from "./components/Stocks.svelte";
  import Stock from "./components/Stock.svelte";
  import Indices from "./components/Indices.svelte";
  import { username as user, current, component, refresh } from "./stores";

  const getInfo = async () => {
    const resp = await fetch("/info");
    const info = await resp.json();
    if (Object.keys(info).length) {
      $user = info.username;
      $refresh = info.refresh;
    }
    if (/^\/stock\/[A-Z]{3,4}\/\d{6}$/.test(window.location.pathname)) {
      const stock = window.location.pathname.split("/");
      $current = { index: stock[2], code: stock[3] };
      window.history.pushState({}, "", `/stock/${stock[2]}/${stock[3]}`);
      $component = "stock";
    }
  };
  const promise = getInfo();

  const components: {
    [component: string]:
      | typeof Login
      | typeof Setting
      | typeof Stocks
      | typeof Stock;
  } = {
    login: Login,
    setting: Setting,
    stocks: Stocks,
    stock: Stock,
  };
</script>

<Nav bind:user={$user} />
{#await promise then _}
  <div class="content">
    <svelte:component this={components[$component]} on:info={getInfo} />
  </div>
{/await}
<Indices />
