<script lang="ts">
  import { onMount } from "svelte";
  import { addColor, checkTradingTime } from "../misc";
  import { mystocks } from "../stock.svelte";

  const names = {
    沪: "上证指数",
    深: "深证成指",
    创: "创业板指",
    中: "中小板指",
  };
  const fields: Array<keyof Stock> = ["now", "change", "percent"];

  let indices = $state.raw<{ [key: string]: Stock }>({});

  const subscribe = async (init?: boolean) => {
    let resp: Response;
    try {
      if (init || (await checkTradingTime())) resp = await fetch("/indices");
      else resp = new Response(null, { status: 400 });
    } catch (e) {
      console.error(e);
      resp = new Response(null, { status: 500 });
    }
    if (resp.ok) {
      indices = await resp.json();
      await new Promise((sleep) => setTimeout(sleep, 10000));
    } else if (resp.status == 400)
      await new Promise((sleep) => setTimeout(sleep, 1000));
    else await new Promise((sleep) => setTimeout(sleep, 30000));
    await subscribe();
  };

  onMount(async () => {
    await subscribe(true);
  });
</script>

{#if Object.keys(indices).length !== 0}
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <div class="indices">
    {#each Object.entries(names) as [key, val] (key)}
      {@const stock = indices[key]}
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <div id={key} onclick={() => mystocks.goto(stock)}>
        <span class="short">{key}</span>
        <span class="full">{val}</span>
        {#each fields as field (field)}
          <span style={addColor(stock, field)}>
            &nbsp;&nbsp;{stock[field]}
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
