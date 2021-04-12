<script>
  import autoComplete from "@tarekraafat/autocomplete.js";
  import { onMount } from "svelte";
  import { current, component } from "../stores";
  import { post } from "../misc";

  let suggest = "";

  onMount(async () => {
    const autoCompletejs = new autoComplete({
      selector: "#suggest",
      data: {
        src: async () => {
          const resp = await post("/suggest", { keyword: suggest });
          const data = await resp.json();
          return data.map((i) => `${i.Index}:${i.Code} ${i.Name} ${i.Type}`);
        },
      },
      trigger: { event: ["input", "focus"] },
      searchEngine: (query, record) => {
        return record;
      },
      placeHolder: "Search Stock",
      threshold: 2,
      debounce: 300,
      resultsList: {
        maxResults: 10,
        noResults: (list) => {
          const result = document.createElement("li");
          result.innerHTML = "No Results";
          list.appendChild(result);
        },
      },
      onSelection: (feedback) => {
        const stock = feedback.selection.value.split(" ")[0].split(":");
        $current = { index: stock[0], code: stock[1] };
        window.history.pushState({}, "", `/stock/${stock[0]}/${stock[1]}`);
        $component = "stock";
        suggest = "";
      },
    });
  });
</script>

<div class="search">
  <div class="icon"><i class="material-icons">search</i></div>
  <input bind:value={suggest} id="suggest" />
</div>

<style>
  .search {
    position: relative;
    width: 250px;
    display: flex;
    float: right;
    margin-bottom: 10px;
    margin-right: 150px;
    background-color: #e6ecf0;
    border-radius: 9999px;
  }

  .search:hover {
    box-shadow: 0 1px 6px 0 rgba(32, 33, 36, 0.28);
  }

  .icon {
    flex-direction: column;
    display: flex;
    justify-content: center;
    padding-left: 20px;
  }

  #suggest {
    background-color: transparent;
    padding: 10px;
    border: 0;
  }

  #suggest:focus {
    outline: none;
  }

  @media (max-width: 1360px) {
    .search {
      margin-right: 0;
    }
  }

  :global(.autoComplete_list) {
    position: absolute;
    background-color: white;
    box-shadow: 0px 5px 4px rgba(101, 119, 134, 0.2),
      5px 2px 4px rgba(101, 119, 134, 0.2),
      -5px 2px 4px rgba(101, 119, 134, 0.2);
    border-radius: 5px;
    list-style-type: none;
    padding: 0;
    width: 250px;
    cursor: default;
    text-indent: 20px;
    top: 100%;
    z-index: 99;
  }

  :global(.autoComplete_result:hover) {
    color: white;
    background-color: #008eff;
    border-radius: 5px;
  }

  :global(.autoComplete_selected) {
    color: white;
    background-color: #008eff;
    border-radius: 5px;
  }
</style>
