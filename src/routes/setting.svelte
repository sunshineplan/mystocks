<script>
  import { BootstrapButtons, post, valid } from "../misc.js";

  let password, password1, password2, validated;

  function setting() {
    if (valid()) {
      this.validated = false;
      post("/setting", {
        password,
        password1,
        password2,
      }).then((resp) => {
        if (!resp.ok)
          resp
            .text()
            .then((err) => BootstrapButtons.fire("Error", err, "error"));
        else
          resp.json().then((json) => {
            if (json.status == 1)
              BootstrapButtons.fire(
                "Success",
                "Your password has changed. Please Re-login!",
                "success"
              ).then(() => (window.location = "/"));
            else
              BootstrapButtons.fire("Error", json.message, "error").then(() => {
                if (json.error == 1) password = "";
                else {
                  password1 = "";
                  password2 = "";
                }
              });
          });
      });
    } else this.validated = true;
  }
  function goback() {
    this.$router.go(-1);
  }
  function cancel(event) {
    if (event.key == "Escape") this.goback();
  }
</script>

<svelte:window on:keydown={cancel} />

<svelte:head>
  <title>Setting</title>
</svelte:head>

<div
  class="content"
  on:keyup={(e) => {
    if (e.key == 'Enter') setting();
  }}>
  <header style="padding-left: 20px">
    <h3>Setting</h3>
    <hr />
  </header>
  <div
    style="margin-left: 120px; width: 250px"
    class="was-validated: {validated}">
    <div class="form-group">
      <label for="password">Current Password</label>
      <input
        class="form-control"
        type="password"
        bind:value={password}
        id="password"
        maxlength="20"
        required />
      <div class="invalid-feedback">This field is required.</div>
    </div>
    <div class="form-group">
      <label for="password1">New Password</label>
      <input
        class="form-control"
        type="password"
        bind:value={password1}
        id="password1"
        maxlength="20"
        required />
      <div class="invalid-feedback">This field is required.</div>
    </div>
    <div class="form-group">
      <label for="password2">Confirm Password</label>
      <input
        class="form-control"
        type="password"
        bind:value={password2}
        id="password2"
        maxlength="20"
        required />
      <div class="invalid-feedback">This field is required.</div>
      <small class="form-text text-muted">
        Max password length: 20 characters.
      </small>
    </div>
    <button class="btn btn-primary" on:click={setting}>Change</button>
    <button class="btn btn-primary" on:click={goback}>Cancel</button>
  </div>
</div>
