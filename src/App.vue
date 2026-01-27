<template>
  <NConfigProvider :theme="theme">
    <MobileBanner />
    <NLayout size="giant" class="app-layout">
      <NLayoutHeader>
        <NCard>
          <NavBar
            :loading="loading"
            :projects="projects"
            v-model:params="params"
            @on-refresh="onRefresh"
            @on-project-select="onProjectSelect"
            @on-params-change="onParamsChange"
          />
        </NCard>
      </NLayoutHeader>
      <NLayoutContent class="app-content" scroll>
        <DataTable
          class="app-data-table"
          @on-sort="onSort"
          :articles="articles"
          :loading="loading"
          :projects="projects"
          v-model:params="params"
        />
      </NLayoutContent>
    </NLayout>
  </NConfigProvider>
</template>

<script>
import NavBar from "./components/NavBar/NavBar.vue";
import DataTable from "./components/DataTable/DataTable.vue";
import MobileBanner from "./components/MobileBanner/MobileBanner.vue";
import { NConfigProvider, NLayout, NLayoutHeader, NLayoutContent, NCard } from "naive-ui";
import { darkTheme } from "naive-ui";
import { query } from "./utils/api";

export default {
  name: "App",
  components: {
    NavBar,
    NConfigProvider,
    NCard,
    NLayout,
    NLayoutHeader,
    NLayoutContent,
    DataTable,
    MobileBanner
},
  methods: {
    async getArticles() {
      this.loading = true;
      localStorage.setItem("params", JSON.stringify(this.params));
      
      this.articles = await query("articles", this.params);
      this.loading = false;
    },
    async getLanguages() {
      const languages = await query("languages");
      const projects = {};

      languages.forEach(({ localname, site }) => {
        site.forEach(({ dbname, code }) => {
        if (code === "wiki") {
          code = "wikipedia"; 
          projects[dbname] = `${localname} ${code.charAt(0).toUpperCase()}${code.slice(
            1
          )}`;
        }
        });
      });

      this.projects = projects;
    },
    async onSort() {
      this.getArticles();
    },
    async onRefresh() {
      this.getArticles();
    },
    async onProjectSelect() {
      this.getArticles();
    },
    async onParamsChange() {
      this.getArticles();
    },
  },
  async mounted() {
    this.getLanguages();
    this.getArticles();
  },
  data() {
    const params = localStorage.getItem("params")
      ? JSON.parse(localStorage.getItem("params"))
      : null;

    return {
      loading: true,
      articles: [],
      projects: {},
      params: params || [
        ["sort", "indications|desc"],
        ["sort", "editors|desc"],
        ["editors_sign", ">"],
        ["editors", 0],
        ["indications_sign", "="],
        ["edits_sign", ">"],
        ["number_of_edits", 0]
      ],
    };
  },
  setup() {
    const primaryColor = "#4263eb";
    const theme = {
      ...darkTheme,
    };
    theme.common.primaryColor = primaryColor;
    theme.common.primaryColorHover = primaryColor;
    theme.common.fontSize = "18px";

    return {
      theme,
    };
  },
};
</script>

<style>
.app-data-table {
  height: calc(100vh - 111px);
}

.n-page-header__extra {
  width: 100%;
  display: flex;
  justify-content: flex-end;
  align-items: center;
}

body {
  overflow-y: hidden;
  height: 100vh;
}


</style>
