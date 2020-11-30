<script>
  import { onMount } from "svelte";
  import AutoComplete from "../components/AutoComplete.svelte";
  import Sortable from "sortablejs";
  import { refresh } from "../stores.js";
  import { checkTime, post, gotoStock, addColor } from "../misc.js";

  let columns = {
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
  let stocks = [];
  let autoUpdate, fetching;

  function start() {
    load(true);
    autoUpdate = setInterval(load, $refresh * 1000);
  }

  function stop() {
    fetching.abort();
    clearInterval(autoUpdate);
  }

  function load(force) {
    if (checkTime() || force) {
      fetching = new AbortController();
      fetch("/mystocks", { signal: fetching.signal })
        .then((response) => response.json())
        .then((json) => {
          stocks = json;
        });
    }
  }

  onMount(() => {
    start();
    const sortable = new Sortable(document.querySelector("#sortable"), {
      animation: 150,
      delay: 500,
      swapThreshold: 0.5,
      onStart: stop(),
      onEnd: start(),
      onUpdate: (evt) => {
        post("/reorder", {
          old: `${stocks[evt.oldIndex].index} ${stocks[evt.oldIndex].code}`,
          new: `${stocks[evt.newIndex].index} ${stocks[evt.newIndex].code}`,
        });
      },
    });
    return () => {
      stop();
      sortable.destroy();
    };
  });
</script>

<style>
  .table-responsive {
    height: calc(100% - 54px);
    padding: 0px 30px;
    cursor: default;
  }

  th {
    position: sticky;
    top: 0;
    border-top: 0 !important;
    border-bottom: 0 !important;
    background-color: white;
  }

  .sortable-ghost {
    opacity: 0;
  }
</style>

<svelte:head>
  <title>My Stocks</title>
</svelte:head>

<div class="content">
  <header>
    <AutoComplete />
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
          <tr on:click={gotoStock(stock)}>
            {#each Object.entries(columns) as [key, val] (key)}
              <td style={addColor(stock, val)}>{stock[val]}</td>
            {/each}
          </tr>
        {/each}
      </tbody>
    </table>
  </div>
</div>
