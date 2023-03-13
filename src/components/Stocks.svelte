<script lang="ts">
  import Sortable from "sortablejs";
  import { onMount } from "svelte";
  import AutoComplete from "./AutoComplete.svelte";
  import { checkTime, post, addColor } from "../misc";
  import { current, component, refresh } from "../stores";

  const columns = {
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
  } as { [key: string]: keyof Stock };

  let stocks: Stock[] = [];
  let autoUpdate: NodeJS.Timeout;
  let fetching: AbortController;

  const start = async () => {
    await load(true);
    autoUpdate = setInterval(load, $refresh * 1000);
  };

  const stop = () => {
    fetching.abort();
    clearInterval(autoUpdate);
  };

  const load = async (force?: boolean) => {
    if (checkTime() || force) {
      fetching = new AbortController();
      const resp = await fetch("/mystocks", { signal: fetching.signal });
      stocks = await resp.json();
    }
  };

  const goto = (stock: Stock) => {
    $current = stock;
    window.history.pushState({}, "", `/stock/${stock.index}/${stock.code}`);
    $component = "stock";
  };

  const onUpdate = async (evt: Sortable.SortableEvent) => {
    await post("/reorder", {
      old: `${stocks[evt.oldIndex as number].index} ${
        stocks[evt.oldIndex as number].code
      }`,
      new: `${stocks[evt.newIndex as number].index} ${
        stocks[evt.newIndex as number].code
      }`,
    });
  };

  onMount(() => {
    start();
    const sortable = new Sortable(
      document.querySelector("#sortable") as HTMLElement,
      {
        animation: 150,
        delay: 400,
        swapThreshold: 0.5,
        onStart: stop,
        onEnd: start,
        onUpdate,
      }
    );
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
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <span style="padding-left:10px" on:click={() => post("/refresh")}>
    <i class="material-icons refresh">refresh</i>
  </span>
</header>
<div class="table-responsive">
  <table class="table table-hover table-sm">
    <thead>
      <tr>
        {#each Object.entries(columns) as [key, val] (key)}
          <th>{key}</th>
        {/each}
      </tr>
    </thead>
    <tbody id="sortable">
      {#each stocks as stock (stock.index + stock.code)}
        <tr on:click={() => goto(stock)}>
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

  :global(.sortable-ghost) {
    opacity: 0;
  }
</style>
