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
    };
  },
  mounted() {
    const autoCompletejs = new autoComplete({
      selector: "#suggest",
      data: {
        src: async () => {
          const resp = await post("/suggest", { keyword: this.suggest });
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
      maxResults: 10,
      noResults: (dataFeedback, generateList) => {
        generateList(autoCompletejs, dataFeedback, dataFeedback.results);
        const result = document.createElement("li");
        result.innerHTML = "No Results";
        document.querySelector("#autoComplete_list").appendChild(result);
      },
      onSelection: (feedback) => {
        const stock = feedback.selection.value.split(" ")[0].split(":");
        this.$router.push(`/stock/${stock[0]}/${stock[1]}`);
        this.suggest = "";
      },
    });
  },
};
</script>

<style scoped>
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

<style>
.autoComplete_list {
  position: absolute;
  background-color: white;
  box-shadow: 0px 5px 4px rgba(101, 119, 134, 0.2),
    5px 2px 4px rgba(101, 119, 134, 0.2), -5px 2px 4px rgba(101, 119, 134, 0.2);
  border-radius: 5px;
  list-style-type: none;
  padding: 0;
  width: 250px;
  cursor: default;
  text-indent: 20px;
  top: 100%;
  z-index: 99;
}

.autoComplete_result:hover {
  color: white;
  background-color: #008eff;
  border-radius: 5px;
}

.autoComplete_selected {
  color: white;
  background-color: #008eff;
  border-radius: 5px;
}
</style>
