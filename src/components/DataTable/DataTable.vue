<template>
  <NDataTable
    :columns="columns"
    :data="articles"
    max-height="calc(100vh - 200px)"
    :loading="loading"
    @update:sorter="onSortChange"
    :row-class-name="rowClassName"
  />
</template>

<script>
import { NDataTable } from "naive-ui";
import { h } from "vue";
import ReactionLike from '../Reaction/Reaction.vue';

export default {
  name: "DataTable",
  components: {
    NDataTable,
  },
  props: {
    params: Array,
    articles: Array,
    projects: Object,
    loading: Boolean,
  },
  methods: {
    onSortChange(sorts) {
      const params = Object.assign([], this.params);

      sorts.forEach(({ columnKey, order }) => {
        const idx = params.findIndex(
          ([param, name]) => param == "sort" && name.startsWith(columnKey)
        );
        const name = `${columnKey}|${order == "ascend" ? "asc" : "desc"}`;

        if (idx != -1) {
          if (order) {
            params[idx][1] = name;
          } else {
            params.splice(idx, 1);
          }
        } else if (order) {
          params.push(["sort", name]);
        }
      });

      this.$emit("update:params", params);
      this.$emit("onSort");
    },
  },
  data() {
    const breakingNewsClass = "potential-breaking-news";
    const rowClassName = ({ editors, indications }) => {
      if (editors && editors.length >= 5) {
        return breakingNewsClass;
      }

      if (indications && indications.length > 0) {
        return breakingNewsClass;
      }

      return "";
    };
    const sorter = {
      multiple: 1,
    };
    const nameWidth = 300;
    const columns = [
      {
        title: "Name",
        key: "name",
        width: nameWidth,
        sorter,
        render({ name, url }) {
          return h(
            "a",
            {
              href: url,
              target: "_blank",
              title: name,
              style: `
                max-width: ${nameWidth}px;
                overflow: hidden;
                white-space: nowrap;
                text-overflow: ellipsis;
                color: white;
                text-decoration: none;
                display: block;
              `,
              className: "table-column"
            },
            name
          );
        },
      },
      {
        title: "QID",
        render({ qid }) {
          return h(
            "a",
            {
              href: `https://www.wikidata.org/wiki/${qid}`,
              target: "_blank",
              title: qid,
              style: `
                color: white;
                text-decoration: none;
              `,
              className: "table-column"
            },
            qid
          );
        },
      },
      {
        title: "Unique Editors",
        key: "editors",
        sorter,
        render({ editors }) {
          return editors ? editors.length : 0;
        },
        className: "table-column"
      },
      {
        title: "Unique Editors Within 1st Hour",
        key: "editors_within_hour",
        sorter,
        className: "table-column"
      },
      {
        title: "Anonymous Editors Within 1st Hour",
        key: "anonymous_editors_within_hour",
        sorter,
      },
      {
        title: "Editors Ratio Within 1st Hour",
        key: "anonymous_editors_ratio_within_hour",
        sorter,
        className: "table-column"
      },
      {
        title: "Number Of Edits",
        key: "number_of_edits",
        sorter,
        className: "table-column"
      },
      {
        title: "Project",
        key: "project",
        className: "table-column",
        sorter,
        render: ({ project, project_url }) => {
          const title = this.projects[project];
          return h(
            "a",
            {
              href: project_url,
              target: "_blank",
              title: `${project} - ${title}`,
              style: `
                color: white;
                text-decoration: none;
              `,
            },
            `${title}`
          );
        },
      },
      {
        title: "Reactions",
        key: "reactions",
        className: "table-column reaction-column",
        render: (row) => {
          return h(ReactionLike, { article: row })
        }
      },
      {
        title: "Date Created",
        key: "date_created",
        sorter,
        className: "table-column"
      },
      {
        title: "Date Modified",
        key: "date_modified",
        sorter,
        className: "table-column"
      },
      {
        title: "Date Namespace Moved",
        key: "date_namespace_moved",
        sorter,
        className: "table-column"
      },
    ];
    const pagination = {
      pageSize: 100,
    };

    columns.forEach((column) => {
      const idx = this.params.findIndex(
        ([param, name]) => param === "sort" && name.startsWith(column.key)
      );

      if (idx != -1) {
        column.defaultSortOrder = this.params[idx][1].endsWith("asc")
          ? "ascend"
          : "descend";
      }
    });

    return {
      columns,
      pagination,
      rowClassName,
    };
  },
};
</script>

<style>
.n-data-table tr:not(.potential-breaking-news) .n-data-table-td {
  opacity: 0.8;
}

.table-column {
    width: min-contentt;
    box-sizing: border-box;
  }

.reaction-column {
  width: 200px;
}

@media screen and (max-width: 1200px){
  .reaction-column {
    width: 150px;
  }
}
</style>
