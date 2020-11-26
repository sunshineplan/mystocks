<template>
  <div class="search">
    <div class="icon">
      <i class="material-icons">search</i>
    </div>
    <input v-model.trim="suggest" id="suggest" />
  </div>
</template>

<script>
import autoComplete from "@tarekraafat/autocomplete.js";
import { post } from "@/misc.js";

export default {
  name: "AutoComplete",
  data() {
    return {
      suggest: "",
      autoComplete: "",
    };
  },
  mounted() {
    this.autoComplete = new autoComplete({
      selector: "#suggest",
      data: { src: this.load, cache: false },
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
      onSelection: (feedback) => {
        let stock = feedback.selection.value.split(" ")[0].split(":");
        this.$router.push(`/stock/${stock[0]}/${stock[1]}`);
        this.suggest = "";
      },
    });
    document.querySelector("#suggest").addEventListener("blur", this.hide);
    document
      .querySelector("#suggest")
      .addEventListener(
        "focus",
        () => (document.querySelector("#suggestsList").style.display = "block")
      );
    document.querySelector("#suggest").addEventListener("keyup", (evt) => {
      if (evt.key == "Escape") {
        this.suggest = "";
        this.hide();
      }
    });
    this.hide();
  },
  methods: {
    async load() {
      if (this.suggest.length >= 2) {
        let source = await post("/suggest", { keyword: this.suggest });
        let data = await source.json();
        return data.map((i) => `${i.Index}:${i.Code} ${i.Name} ${i.Type}`);
      }
      return [];
    },
    hide() {
      document.querySelector("#suggestsList").style.display = "none";
    },
  },
};
</script>
