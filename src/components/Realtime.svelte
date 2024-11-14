<script lang="ts">
  import { addColor, color, post } from "../misc";
  import { mystocks } from "../stock.svelte";

  let {
    stock = $bindable(),
  }: {
    stock: Stock;
  } = $props();

  let width = $derived(
    !stock.sell5.length && !stock.buy5.length ? "480px" : "360px",
  );

  const star = async () => {
    await post("/star", {
      action: mystocks.current.stared ? "unstar" : "star",
    });
    mystocks.current.stared = !mystocks.current.stared;
  };

  const open = () => {
    if (stock.index == "SSE")
      window.open(`https://quote.eastmoney.com/sh${stock.code}.html`);
    else if (stock.index == "SZSE")
      window.open(`https://quote.eastmoney.com/sz${stock.code}.html`);
    else if (stock.index == "BSE")
      window.open(`https://quote.eastmoney.com/bj/${stock.code}.html`);
  };
</script>

<div>
  <div style="display: flex; font-size: 2rem">
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <i
      class="material-icons star"
      class:stared={mystocks.current.stared}
      onclick={star}
    >
      {mystocks.current.stared ? "star" : "star_border"}
    </i>
    <span>{stock.name}</span>(<span>{stock.code}</span>)
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <i class="material-icons open" onclick={open}>open_in_new</i>
    &nbsp;&nbsp;&nbsp;
    <span style={addColor(stock, "now")}>{stock.now}</span>
    &nbsp;&nbsp;&nbsp;
    <span style={addColor(stock, "percent")}>{stock.percent}</span>
  </div>
  <div style="min-height: 52px">
    <table style="float: left; table-layout: fixed; width: {width}">
      <tbody>
        <tr>
          <td>昨收: <span>{stock.last}</span></td>
          <td>
            涨跌:
            <span style={addColor(stock, "change")}>{stock.change}</span>
          </td>
          <td>
            涨幅:
            <span style={addColor(stock, "percent")}>{stock.percent}</span>
          </td>
        </tr>
        <tr>
          <td>
            最高:
            <span style={addColor(stock, "high")}>{stock.high}</span>
          </td>
          <td>
            最低:
            <span style={addColor(stock, "low")}>{stock.low}</span>
          </td>
          <td>
            开盘:
            <span style={addColor(stock, "open")}>{stock.open}</span>
          </td>
        </tr>
      </tbody>
    </table>
    {#if stock.sell5.length || stock.buy5.length}
      <table>
        <tbody>
          <tr>
            <td>
              <span style="display: inline-flex">
                卖盘:&nbsp;
                {#each stock.sell5 as sell, index (index)}
                  <div class="sellbuy">
                    <span style={color(stock.last, sell.Price)}>
                      {sell.Price}
                    </span>
                    -
                    <span style={color(stock.last, sell.Price)}>
                      {sell.Volume}
                    </span>
                  </div>
                {/each}
              </span>
            </td>
          </tr>
          <tr>
            <td>
              <span style="display: inline-flex">
                买盘:&nbsp;
                {#each stock.buy5 as buy, index (index)}
                  <div class="sellbuy">
                    <span style={color(stock.last, buy.Price)}>
                      {buy.Price}
                    </span>
                    -
                    <span style={color(stock.last, buy.Price)}>
                      {buy.Volume}
                    </span>
                  </div>
                {/each}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    {/if}
  </div>
  <small>更新时间: <span class="update">{stock.update}</span></small>
</div>

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
