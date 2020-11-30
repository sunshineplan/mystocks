<script>
  import { onMount } from "svelte";
  import { stock } from "../stores.js";
  import { post, addColor } from "../misc.js";

  let stared = false;
  $: width = !$stock.sell5 && !$stock.buy5 ? "480px" : "360px";

  function star() {
    if (stared)
      post("/star", { action: "unstar" }).then(() => (stared = false));
    else post("/star").then(() => (stared = true));
  }

  function open() {
    window.open("http://stockpage.10jqka.com.cn/" + $stock.code);
  }

  onMount(() => {
    fetch("/star")
      .then((response) => response.text())
      .then((text) => {
        if (text == "1") stared = true;
      });
  });
</script>

<style>
  .star {
    color: #f4b400;
    width: 50px;
    height: 50px;
    font-size: 40px;
    cursor: default;
    display: flex;
    justify-content: center;
    align-items: center;
  }

  .star:hover {
    background-color: #fcebbaa1;
    border-radius: 50%;
  }

  .open {
    margin-left: 0.5rem;
    color: gray;
    cursor: pointer;
    display: flex;
    align-items: center;
  }

  .sellbuy {
    min-width: 115px;
    padding-right: 6px;
  }
</style>

<div>
  <div style="display: flex; font-size: 2rem">
    <i class="material-icons star {stared ? 'stared' : ''}" on:click={star}>
      {stared ? 'star' : 'star_border'}
    </i>
    <span>{$stock.name}</span>(<span>{$stock.code}</span>)
    <i class="material-icons open" on:click={open}> open_in_new </i>
    &nbsp;&nbsp;&nbsp;
    <span style={addColor($stock, 'now')}>{$stock.now}</span>
    &nbsp;&nbsp;&nbsp;
    <span style={addColor($stock, 'percent')}>{$stock.percent}</span>
  </div>
  <div style="min-height: 52px">
    <table style="float: left; table-layout: fixed; width: {width}">
      <tbody>
        <tr>
          <td>昨收: <span>{$stock.last}</span></td>
          <td>
            涨跌:
            <span style={addColor($stock, 'change')}>{$stock.change}</span>
          </td>
          <td>
            涨幅:
            <span style={addColor($stock, 'percent')}>{$stock.percent}</span>
          </td>
        </tr>
        <tr>
          <td>
            最高:
            <span style={addColor($stock, 'high')}>{$stock.high}</span>
          </td>
          <td>
            最低:
            <span style={addColor($stock, 'low')}>{$stock.low}</span>
          </td>
          <td>
            开盘:
            <span style={addColor($stock, 'open')}>{$stock.open}</span>
          </td>
        </tr>
      </tbody>
    </table>
    {#if $stock.sell5 || $stock.buy5}
      <table>
        <tbody>
          <tr>
            <td>
              <span style="display: inline-flex">
                卖盘:&nbsp;
                {#each $stock.sell5 as sell, index (index)}
                  <div class="sellbuy" style="color: red">
                    {sell.Price}-{sell.Volume}
                  </div>
                {/each}
              </span>
            </td>
          </tr>
          <tr>
            <td>
              <span style="display: inline-flex">
                买盘:&nbsp;
                {#each $stock.buy5 as buy, index (index)}
                  <div class="sellbuy" style="color: green">
                    {buy.Price}-{buy.Volume}
                  </div>
                {/each}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    {/if}
  </div>
  <small>更新时间: <span class="update">{$stock.update}</span></small>
</div>
