<script>
  import { onMount } from "svelte";
  import { goto } from "@sapper/app";
  import { BootstrapButtons, post } from "../misc.js";

  let username, password, rememberme;

  async function login() {
    if (!document.querySelector("#username").checkValidity())
      await BootstrapButtons.fire(
        "Error",
        "Username cannot be empty.",
        "error"
      );
    else if (!document.querySelector("#password").checkValidity())
      await BootstrapButtons.fire(
        "Error",
        "Password cannot be empty.",
        "error"
      );
    else {
      const resp = await post("/login", { username, password, rememberme });
      if (!resp.ok)
        await BootstrapButtons.fire("Error", await resp.text(), "error");
      else {
        localStorage.setItem("username", username);
        await goto("/");
      }
    }
  }

  onMount(() => {
    username = localStorage.getItem("username");
    document.querySelector("#username").focus();
  });
</script>

<style>
  .login {
    width: 250px;
    margin: 0 auto;
  }
</style>

<svelte:head>
  <title>Log In</title>
</svelte:head>

<div
  class="content"
  on:keyup={async (e) => {
    if (e.key == 'Enter') await login();
  }}>
  <header>
    <h3
      class="d-flex justify-content-center align-items-center"
      style="height: 100%">
      Log In
    </h3>
  </header>
  <div class="login">
    <div class="form-group">
      <label for="username">Username</label>
      <input
        class="form-control"
        bind:value={username}
        id="username"
        maxlength="20"
        placeholder="Username"
        required />
    </div>
    <div class="form-group">
      <label for="password">Password</label>
      <input
        class="form-control"
        type="password"
        bind:value={password}
        id="password"
        maxlength="20"
        placeholder="Password"
        required />
    </div>
    <div class="form-group form-check">
      <input
        type="checkbox"
        class="form-check-input"
        bind:checked={rememberme}
        id="rememberme" />
      <label class="form-check-label" for="rememberme">Remember Me</label>
    </div>
    <hr />
    <button
      class="btn btn-primary login"
      on:click={async () => await login()}>Log In</button>
  </div>
</div>
