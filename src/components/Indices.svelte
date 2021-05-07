<script lang="ts">
  import { onMount } from "svelte";
  import { checkTime, addColor } from "../misc";
  import { current, component } from "../stores";
  import type { Stock } from "../index";

  const names = {
    沪: "上证指数",
    深: "深证成指",
    创: "创业板指",
    中: "中小板指",
  };
  const fields = ["now", "change", "percent"] as Array<keyof Stock>;

  let indices: { [key: string]: Stock } = {};

  const start = async () => {
    await load(true);
    setInterval(load, 10000);
  };

  const load = async (force?: boolean) => {
    if (checkTime() || force) {
      const resp = await fetch("/indices");
      indices = await resp.json();
    }
  };

  const goto = (stock: Stock) => {
    $current = stock;
    window.history.pushState({}, "", `/stock/${stock.index}/${stock.code}`);
    $component = "stock";
  };

  onMount(async () => {
    await start();
  });
</script>

{#if Object.keys(indices).length !== 0}
  <div class="indices">
    {#each Object.entries(names) as [key, val] (key)}
      <div id={key} on:click={() => goto(indices[key])}>
        <span class="short">{key}</span>
        <span class="full">{val}</span>
        {#each fields as field (field)}
          <span style={addColor(indices[key], field)}>
            &nbsp;&nbsp;{indices[key][field]}
          </span>
        {/each}
      </div>
    {/each}
  </div>
{/if}

<style>
  .indices {
    position: fixed;
    z-index: 100;
    bottom: 0;
    width: 100%;
    height: 70px;
    display: flex;
    align-items: center;
    background-color: white;
    box-shadow: 0 -1px 2px 0 #e7e7e7;
    white-space: normal;
  }

  #沪,
  #深,
  #创,
  #中 {
    color: black;
    max-width: 25%;
    flex: 0 0 25%;
    cursor: default;
    text-align: center;
    font-size: 20px;
  }

  #沪,
  #深,
  #创,
  #中:hover {
    text-decoration: none;
  }

  .short {
    display: none;
  }

  @media (max-width: 1360px) {
    .short {
      display: inline;
    }

    .full {
      display: none;
    }
  }
</style>
