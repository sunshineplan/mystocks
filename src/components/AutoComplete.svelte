<script>
  import { onMount } from "svelte";
  import { goto } from "@sapper/app";
  import { post } from "../misc.js";

  let suggest;

  async function load() {
    if (suggest.length >= 2) {
      let source = await post("/suggest", { keyword: suggest });
      let data = await source.json();
      return data.map((i) => `${i.Index}:${i.Code} ${i.Name} ${i.Type}`);
    }
    return [];
  }

  function hide() {
    document.querySelector("#suggestsList").style.display = "none";
  }

  onMount(async () => {
    const autoComplete = await import("@tarekraafat/autocomplete.js");
    new autoComplete({
      selector: "#suggest",
      data: { src: load, cache: false },
      trigger: { event: ["input", "focus"] },
      searchEngine: (query, record) => {
        return record;
      },
      placeHolder: "Search Stock",
      threshold: 1,
      debounce: 200,
      maxResults: 10,
      resultsList: {
        render: true,
        container: (source) => {
          source.setAttribute("id", "suggestsList");
          source.setAttribute("class", "suggestsList");
        },
      },
      resultItem: {
        content: (data, src) => {
          src.innerHTML = data.match;
        },
      },
      noResults: () => {
        let result = document.createElement("li");
        result.innerHTML = "No Results";
        document.querySelector("#suggestsList").appendChild(result);
      },
      onSelection: async (feedback) => {
        let stock = feedback.selection.value.split(" ")[0].split(":");
        await goto(`/${stock[0]}/${stock[1]}`);
        suggest = "";
      },
    });
    document.querySelector("#suggest").addEventListener("blur", hide);
    document
      .querySelector("#suggest")
      .addEventListener(
        "focus",
        () => (document.querySelector("#suggestsList").style.display = "block")
      );
    document.querySelector("#suggest").addEventListener("keyup", (evt) => {
      if (evt.key == "Escape") {
        suggest = "";
        hide();
      }
    });
    hide();
  });
</script>

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

<div class="search">
  <div class="icon"><i class="material-icons">search</i></div>
  <input bind:value={suggest} id="suggest" />
</div>
