<script>
  import Cookies from "js-cookie";
  import { goto } from "@sapper/app";

  $: user = Cookies.get("Username");
</script>

<style>
  .topbar {
    position: fixed;
    top: 0;
    z-index: 2;
    width: 100%;
    height: 60px;
    background-color: #1a73e8;
  }

  .brand {
    padding-left: 20px;
    font-size: 25px;
    letter-spacing: 0.3px;
    color: white;
  }

  .brand:hover {
    color: white;
    text-decoration: none;
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
</style>

<nav class="navbar navbar-light topbar">
  <a class="brand" on:click={goto('/')}>My Stocks</a>
  <div class="navbar-nav flex-row">
    {#if user}
      <div class="navbar-nav flex-row">
        <a class="nav-link">{user}</a>
        <a class="nav-link link" on:click={goto('/setting')}>Setting</a>
        <a class="nav-link link" href="/logout">Log out</a>
      </div>
    {:else}
      <div class="navbar-nav flex-row">
        <a class="nav-link link" on:click={goto('/login')}>Log in</a>
      </div>
    {/if}
  </div>
</nav>
