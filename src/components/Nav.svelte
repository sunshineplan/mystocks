<script lang="ts">
  import { fire, post } from "../misc";
  import { username, component, flows } from "../stores";

  export let user: string;

  const logout = async () => {
    const resp = await post(window.universal + "/logout", undefined, true);
    if (resp.ok) {
      $username = "";
      window.history.pushState({}, "", "/");
      $component = "stocks";
    } else await fire("Error", "Unknow error", "error");
  };
</script>

<nav class="navbar navbar-light topbar">
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <span
    class="brand"
    on:click={() => {
      window.history.pushState({}, "", "/");
      $component = "stocks";
    }}
  >
    My Stocks
  </span>
  <!-- svelte-ignore a11y-click-events-have-key-events -->
  <!-- svelte-ignore a11y-no-static-element-interactions -->
  <span style="color:white" on:click={flows.toggle}>
    Switch to {$flows ? "Stocks" : "Flows"}
  </span>
  <div class="navbar-nav flex-row">
    {#if user}
      <div class="navbar-nav flex-row">
        <span class="nav-link">{user}</span>
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <!-- svelte-ignore a11y-no-static-element-interactions -->
        <span
          class="nav-link link"
          on:click={() => {
            window.history.pushState({}, "", "/setting");
            $component = "setting";
          }}
        >
          Setting
        </span>
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <!-- svelte-ignore a11y-no-static-element-interactions -->
        <span class="nav-link link" on:click={logout}>Log out</span>
      </div>
    {:else}
      <div class="navbar-nav flex-row">
        <!-- svelte-ignore a11y-click-events-have-key-events -->
        <!-- svelte-ignore a11y-no-static-element-interactions -->
        <span
          class="nav-link link"
          on:click={() => {
            window.history.pushState({}, "", "/login");
            $component = "login";
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
  }
</style>
