<script lang="ts">
  import Nav from "./components/Nav.svelte";
  import Login from "./components/Login.svelte";
  import Setting from "./components/Setting.svelte";
  import Stocks from "./components/Stocks.svelte";
  import Stock from "./components/Stock.svelte";
  import Indices from "./components/Indices.svelte";
  import { username as user, component, refresh } from "./stores";

  const getInfo = async () => {
    const resp = await fetch("/info");
    const info = await resp.json();
    if (Object.keys(info).length) {
      $user = info.username;
      $refresh = info.refresh;
    }
  };
  const promise = getInfo();

  const components = {
    login: Login,
    setting: Setting,
    stocks: Stocks,
    stock: Stock,
  } as { [component: string]: any };
</script>

<Nav bind:user={$user} />
{#await promise then _}
  <div class="content">
    <svelte:component this={components[$component]} on:info={getInfo} />
  </div>
{/await}
<Indices />
