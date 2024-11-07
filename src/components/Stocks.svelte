<script lang="ts">
  import Sortable from "sortablejs";
  import { onMount } from "svelte";
  import AutoComplete from "./AutoComplete.svelte";
  import { checkTradingTime, post, addColor } from "../misc";
  import { mystocks } from "../stock.svelte";

  const columns: { [key: string]: keyof Stock } = {
    指数: "index",
    代码: "code",
    名称: "name",
    最新: "now",
    涨跌: "change",
    涨幅: "percent",
    最高: "high",
    最低: "low",
    开盘: "open",
    昨收: "last",
  };

  let stocks: Stock[] = $state([]);
  let autoUpdate: number;
  let fetching: AbortController;
  let table: HTMLElement;

  const start = async () => {
    await load(true);
    autoUpdate = setInterval(load, mystocks.refresh * 1000);
  };

  const stop = () => {
    fetching.abort();
    clearInterval(autoUpdate);
  };

  const load = async (force?: boolean) => {
    if (force || (await checkTradingTime())) {
      fetching = new AbortController();
      const resp = await fetch("/mystocks", { signal: fetching.signal });
      stocks = await resp.json();
    }
  };

  const goto = (stock: Stock) => {
    mystocks.current = stock;
    window.history.pushState({}, "", `/stock/${stock.index}/${stock.code}`);
    mystocks.component = "stock";
  };

  const onUpdate = async (evt: Sortable.SortableEvent) => {
    await post("/reorder", {
      old: `${stocks[evt.oldIndex].index} ${stocks[evt.oldIndex].code}`,
      new: `${stocks[evt.newIndex].index} ${stocks[evt.newIndex].code}`,
    });
  };

  onMount(() => {
    start();
    const sortable = new Sortable(table, {
      animation: 150,
      delay: 400,
      swapThreshold: 0.5,
      onStart: stop,
      onEnd: start,
      onUpdate,
    });
    return () => {
      stop();
      sortable.destroy();
    };
  });
</script>

<svelte:head>
  <title>My Stocks</title>
</svelte:head>

<header>
  <AutoComplete />
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <span style="padding-left:10px" onclick={() => post("/refresh")}>
    <i class="material-icons refresh">refresh</i>
  </span>
</header>
<div class="table-responsive">
  <table class="table table-hover table-sm">
    <thead>
      <tr>
        {#each Object.keys(columns) as key (key)}
          <th>{key}</th>
        {/each}
      </tr>
    </thead>
    <tbody bind:this={table}>
      {#each stocks as stock (stock.index + stock.code)}
        <tr onclick={() => goto(stock)}>
          {#each Object.entries(columns) as [key, val] (key)}
            <td style={addColor(stock, val)}>{stock[val]}</td>
          {/each}
        </tr>
      {/each}
    </tbody>
  </table>
</div>

<style>
  .refresh {
    transition: transform 0.4s linear;
    font-size: 30px;
    color: #1a73e8;
    cursor: pointer;
  }

  .refresh:hover {
    transform: rotate(360deg);
  }

  .table-responsive {
    height: calc(100% - 54px);
    padding: 0px 30px;
    cursor: default;
    width: 100%;
  }

  th {
    position: sticky;
    top: 0;
    background-color: white;
  }
</style>
