<script lang="ts">
  import { addColor, checkTradingTime, color, post } from "../misc";
  import { mystocks } from "../stock.svelte";

  let timer: number;
  let controller: AbortController;

  $effect(() => {
    mystocks.current.index;
    mystocks.current.code;
    subscribe();
    return abort;
  });

  let width = $derived(
    !mystocks.current.stock.sell5?.length &&
      !mystocks.current.stock.buy5?.length
      ? "480px"
      : "360px",
  );

  const star = async () => {
    await post("/star", {
      action: mystocks.current.stared ? "unstar" : "star",
    });
    mystocks.current.stared = !mystocks.current.stared;
  };

  const open = () => {
    if (mystocks.current.stock.index == "SSE")
      window.open(
        `https://quote.eastmoney.com/sh${mystocks.current.stock.code}.html`,
      );
    else if (mystocks.current.stock.index == "SZSE")
      window.open(
        `https://quote.eastmoney.com/sz${mystocks.current.stock.code}.html`,
      );
    else if (mystocks.current.stock.index == "BSE")
      window.open(
        `https://quote.eastmoney.com/bj/${mystocks.current.stock.code}.html`,
      );
  };

  const subscribe = () => {
    controller = new AbortController();
    const fetchDate = async (force?: boolean) => {
      let resp: Response;
      try {
        if (force || (await checkTradingTime()))
          resp = await post("/realtime", {
            index: mystocks.current.index,
            code: mystocks.current.code,
          });
        else resp = new Response(null, { status: 400 });
      } catch (e) {
        console.error(e);
        resp = new Response(null, { status: 500 });
      }
      let timeout = 30000;
      if (resp.ok) {
        const res = await resp.json();
        if (res.name) {
          mystocks.current.update = res.update;
          delete res.update;
          mystocks.current.stock = res;
          document.title = `${res.name} ${res.now} ${res.percent}`;
        }
        timeout = 10000;
      } else if (resp.status == 400) timeout = 1000;
      timer = setTimeout(fetchDate, timeout);
    };
    fetchDate(true);
  };

  const abort = () => {
    if (timer) clearTimeout(timer);
    if (controller) controller.abort();
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
    <span>{mystocks.current.stock.name}({mystocks.current.stock.code})</span>
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <i class="material-icons open" onclick={open}>open_in_new</i>
    &nbsp;&nbsp;&nbsp;
    <span style={addColor(mystocks.current.stock, "now")}>
      {mystocks.current.stock.now}
    </span>
    &nbsp;&nbsp;&nbsp;
    <span style={addColor(mystocks.current.stock, "percent")}>
      {mystocks.current.stock.percent}
    </span>
  </div>
  <div style="min-height: 52px">
    <table style="float: left; table-layout: fixed; width: {width}">
      <tbody>
        <tr>
          <td>昨收: <span>{mystocks.current.stock.last}</span></td>
          <td>
            涨跌:
            <span style={addColor(mystocks.current.stock, "change")}>
              {mystocks.current.stock.change}
            </span>
          </td>
          <td>
            涨幅:
            <span style={addColor(mystocks.current.stock, "percent")}>
              {mystocks.current.stock.percent}
            </span>
          </td>
        </tr>
        <tr>
          <td>
            最高:
            <span style={addColor(mystocks.current.stock, "high")}>
              {mystocks.current.stock.high}
            </span>
          </td>
          <td>
            最低:
            <span style={addColor(mystocks.current.stock, "low")}>
              {mystocks.current.stock.low}
            </span>
          </td>
          <td>
            开盘:
            <span style={addColor(mystocks.current.stock, "open")}>
              {mystocks.current.stock.open}
            </span>
          </td>
        </tr>
      </tbody>
    </table>
    {#if mystocks.current.stock.sell5?.length || mystocks.current.stock.buy5?.length}
      <table>
        <tbody>
          <tr>
            <td>
              <span style="display: inline-flex">
                卖盘:&nbsp;
                {#each mystocks.current.stock.sell5 as sell, index (index)}
                  <div class="sellbuy">
                    <span
                      style={color(mystocks.current.stock.last, sell.Price)}
                    >
                      {sell.Price}
                    </span>
                    -
                    <span
                      style={color(mystocks.current.stock.last, sell.Price)}
                    >
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
                {#each mystocks.current.stock.buy5 as buy, index (index)}
                  <div class="sellbuy">
                    <span style={color(mystocks.current.stock.last, buy.Price)}>
                      {buy.Price}
                    </span>
                    -
                    <span style={color(mystocks.current.stock.last, buy.Price)}>
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
  <small>
    更新时间: <span class="update">{mystocks.current.update}</span>
  </small>
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
