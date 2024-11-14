<script lang="ts">
  import autoComplete from "@tarekraafat/autocomplete.js";
  import { onMount } from "svelte";
  import { post } from "../misc";
  import { mystocks } from "../stock.svelte";

  let suggest = $state("");

  onMount(() => {
    const autoCompleteJS = new autoComplete({
      selector: "#suggest",
      data: {
        src: async () => {
          const resp = await post("/suggest", { keyword: suggest });
          const data = await resp.json();
          return data.map((i) => `${i.Index}:${i.Code} ${i.Name} ${i.Type}`);
        },
      },
      searchEngine: (query, record) => {
        return record;
      },
      placeHolder: "Search Stock",
      threshold: 2,
      debounce: 300,
      resultsList: {
        maxResults: 10,
        noResults: true,
        element: (list, data) => {
          if (!data.results.length) {
            const result = document.createElement("span");
            result.innerText = "No Results";
            list.appendChild(result);
          }
        },
      },
      events: {
        input: {
          focus() {
            if (suggest) autoCompleteJS.start();
          },
          selection(event) {
            const stock = event.detail.selection.value.split(" ")[0].split(":");
            mystocks.goto(stock[0], stock[1]);
            suggest = "";
          },
        },
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
</style>
