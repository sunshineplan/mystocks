<script>
  import Nav from "../components/Nav.svelte";
  import Indices from "../components/Indices.svelte";
</script>

<style>
  :global(body) {
    font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto,
      "Helvetica Neue", Arial, "Noto Sans", "Microsoft YaHei New",
      "Microsoft Yahei", 微软雅黑, 宋体, SimSun, STXihei, 华文细黑, sans-serif;
    white-space: nowrap;
  }

  :global(header) {
    padding: 10px 20px;
  }

  :global(small) {
    color: gray;
    padding-left: 5px;
  }

  :global(button + button) {
    margin-left: 0.3em;
  }

  :global(.content) {
    position: fixed;
    top: 0;
    padding: 60px 0 80px;
    width: 100%;
    height: 100%;
  }

  :global(.icon) {
    flex-direction: column;
    display: flex;
    justify-content: center;
    padding-left: 20px;
  }
</style>

<Nav />

<main>
  <slot />
</main>

<Indices />
