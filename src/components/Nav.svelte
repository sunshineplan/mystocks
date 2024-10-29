<script lang="ts">
  import { fire, post } from "../misc";
  import { isFlows, mystocks } from "../stock.svelte";

  const logout = async () => {
    const resp = await post(window.universal + "/logout", undefined, true);
    if (resp.ok) {
      mystocks.username = "";
      window.history.pushState({}, "", "/");
      mystocks.component = "stocks";
    } else await fire("Error", "Unknow error", "error");
  };
</script>

<nav class="navbar navbar-light topbar">
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <span
    class="brand"
    onclick={() => {
      window.history.pushState({}, "", "/");
      mystocks.component = "stocks";
    }}
  >
    My Stocks
  </span>
  <!-- svelte-ignore a11y_click_events_have_key_events -->
  <!-- svelte-ignore a11y_no_static_element_interactions -->
  <span style="color:white" onclick={() => isFlows.toggle()}>
    Switch to {isFlows.status ? "Stocks" : "Flows"}
  </span>
  <div class="navbar-nav flex-row">
    {#if mystocks.username}
      <div class="navbar-nav flex-row">
        <span class="nav-link">{mystocks.username}</span>
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <span
          class="nav-link link"
          onclick={() => {
            window.history.pushState({}, "", "/setting");
            mystocks.component = "setting";
          }}
        >
          Setting
        </span>
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <span class="nav-link link" onclick={logout}>Log out</span>
      </div>
    {:else}
      <div class="navbar-nav flex-row">
        <!-- svelte-ignore a11y_click_events_have_key_events -->
        <!-- svelte-ignore a11y_no_static_element_interactions -->
        <span
          class="nav-link link"
          onclick={() => {
            window.history.pushState({}, "", "/login");
            mystocks.component = "login";
          }}
        >
          Log in
        </span>
      </div>
    {/if}
  </div>
</nav>

<style>
  .topbar {
    position: fixed;
    top: 0;
    z-index: 2;
    width: 100%;
    height: 60px;
    background-color: #1a73e8;
    padding: 0.5rem 1rem;
  }

  .brand {
    padding-left: 20px;
    font-size: 25px;
    letter-spacing: 0.3px;
    color: white;
  }

  .topbar .nav-link {
    padding-left: 8px;
    padding-right: 8px;
    color: white !important;
  }

  .topbar .link:hover {
    background: rgba(255, 255, 255, 0.2);
    border-radius: 5px;
  }

  span {
    cursor: pointer;
    user-select: none;
  }
</style>
