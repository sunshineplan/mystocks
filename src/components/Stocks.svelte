<script lang="ts">
  import Sortable from "sortablejs";
  import { onMount } from "svelte";
  import { addColor, checkTradingTime, post } from "../misc";
  import { mystocks } from "../stock.svelte";
  import AutoComplete from "./AutoComplete.svelte";

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
  let controller: AbortController;
  let table: HTMLElement;

  const subscribe = async (force?: boolean) => {
    controller = new AbortController();
    let resp: Response;
    try {
      if (force || (await checkTradingTime()))
        resp = await fetch("/mystocks", { signal: controller.signal });
      else resp = new Response(null, { status: 400 });
    } catch (e) {
      if (e instanceof DOMException && e.name === "AbortError") return;
      console.error(e);
      resp = new Response(null, { status: 500 });
    }
    let timeout = 30000;
    if (resp.ok) {
      stocks = await resp.json();
      timeout = mystocks.refresh * 1000;
    } else if (resp.status == 400) timeout = mystocks.refresh * 1000;
    controller.abort("sleep");
    await new Promise((sleep) => setTimeout(sleep, timeout));
    if (controller.signal.aborted && controller.signal.reason == "sleep")
      await subscribe();
  };

  const abort = () => {
    if (controller.signal.aborted) controller = new AbortController();
    controller.abort();
  };

  onMount(() => {
    subscribe(true);
    const sortable = new Sortable(table, {
      animation: 150,
      delay: 400,
      swapThreshold: 0.5,
      onStart: abort,
      onEnd: () => subscribe(true),
      onUpdate: async (evt: Sortable.SortableEvent) => {
        await post("/reorder", {
          old: `${stocks[evt.oldIndex].index} ${stocks[evt.oldIndex].code}`,
          new: `${stocks[evt.newIndex].index} ${stocks[evt.newIndex].code}`,
        });
      },
    });
    return () => {
      abort();
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
        <tr onclick={() => mystocks.goto(stock)}>
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
