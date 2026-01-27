<template>
  <NPageHeader>
      <template #avatar>
          <div class="logo-container">
            <NImage preview-disabled :src="require('./../../assets/header__logo.png')" />
            <div class="beta"><span>beta</span></div>
          </div>
      </template>
    <template #extra>
      <div class="nav-bar-extra">
        <NGrid :cols="5" :x-gap="12" :y-gap="12">
          <NGi :span="3">
            <NGrid :cols="3" :x-gap="12" :y-gap="12">
              <NGi>
                <NCard title="Unique Editors" size="small">
                  <template #header-extra>
                    <NButton
                      secondary
                      round
                      bordered
                      type="error"
                      @click="resetEditors"
                      :disabled="loading"
                    >
                      Reset
                    </NButton>
                  </template>
                  <NGrid :cols="2">
                    <NGi>
                      <NSelect
                        class="nav-bar-short-select"
                        placeholder="Sign"
                        clearable
                        :disabled="loading"
                        :options="signs"
                        v-model:value="editorsSign"
                      />
                    </NGi>
                    <NGi>
                      <NInputNumber
                        clearable
                        class="nav-bar-short-select"
                        placeholder="#"
                        v-model:value="editors"
                        :disabled="loading"
                      />
                    </NGi>
                  </NGrid>
                </NCard>
              </NGi>
              <NGi>
                <NCard title="Number of Edits" size="small">
                  <template #header-extra>
                    <NButton
                      secondary
                      round
                      bordered
                      type="error"
                      @click="resetEdits"
                      :disabled="loading"
                    >
                      Reset
                    </NButton>
                  </template>
                  <NGrid :cols="2">
                    <NGi>
                      <NSelect
                        class="nav-bar-short-select"
                        placeholder="Sign"
                        :disabled="loading"
                        :options="signs"
                        v-model:value="editsSign"
                      />
                    </NGi>
                    <NGi>
                      <NInputNumber
                        class="nav-bar-short-select"
                        placeholder="#"
                        v-model:value="edits"
                        :disabled="loading"
                      />
                    </NGi>
                  </NGrid>
                </NCard>
              </NGi>
              <NGi>
                <NCard title="Indications Count" size="small">
                  <template #header-extra>
                    <NButton
                      secondary
                      round
                      bordered
                      type="error"
                      @click="resetIndications"
                      :disabled="loading"
                    >
                      Reset
                    </NButton>
                  </template>
                  <NGrid :cols="2">
                    <NGi>
                      <NSelect
                        class="nav-bar-short-select"
                        placeholder="Sign"
                        :disabled="loading"
                        :options="indicationsSigns"
                        v-model:value="indicationsSign"
                      />
                    </NGi>
                    <NGi>
                      <NInputNumber
                        class="nav-bar-short-select"
                        placeholder="#"
                        v-model:value="indicationsCount"
                        :disabled="loading"
                      />
                    </NGi>
                  </NGrid>
                </NCard>
              </NGi>
            </NGrid>
          </NGi>
          <NGi :span="2">
            <NGrid cols="1 2xl:2" rows="2" :x-gap="12" :y-gap="12" responsive="screen">
              <NGi>
                <NInput
                  class="nav-bar-select"
                  placeholder="QID"
                  @change="onChange"
                  clearable
                />
              </NGi>
              <NGi>
                <NSelect
                  class="nav-bar-select"
                  placeholder="All Projects"
                  filterable
                  multiple
                  :filter="onFilter"
                  :disabled="loading"
                  :options="options"
                  :on-update:value="onSelect"
                  :default-value="defaultOptions"
                />
              </NGi>
              <NGi>
              </NGi>
              <NGi>
                <NGrid :cols="1" :x-gap="12" :y-gap="12">
                  <NGi>
                  </NGi>
                  <NGi>
                    <NButton
                      strong
                      round
                      secondary
                      bordered
                      type="info"
                      class="nav-bar-button"
                      :disabled="loading"
                      @click="onRefreshClick"
                      >Refresh</NButton
                    >
                  </NGi>
                </NGrid>
              </NGi>
            </NGrid>
          </NGi>
        </NGrid>
      </div>
    </template>
  </NPageHeader>
</template>

<script>
import { ref } from "vue";
import {
  NPageHeader,
  NButton,
  NImage,
  NSelect,
  NInput,
  NInputNumber,
  NCard,
  NGrid,
  NGi,
} from "naive-ui";

export default {
  name: "NavBar",
  props: {
    loading: Boolean,
    projects: Object,
    params: Array,
  },
  components: {
    NPageHeader,
    NButton,
    NImage,
    NSelect,
    NInput,
    NInputNumber,
    NCard,
    NGrid,
    NGi,
  },
  methods: {
    onFilter(text, { label }) {
      return label.toLowerCase().startsWith(text.toLowerCase());
    },
    onRefreshClick() {
      this.$emit("onRefresh");
    },
    onSelect(projects) {
      const params = Object.assign([], this.params).filter(
        ([param]) => param != "project"
      );
      projects.forEach((project) => {
        params.push(["project", project]);
      });

      this.$emit("update:params", params);
      this.$emit("onProjectSelect");
    },
    onChange(qid) {
      const params = Object.assign([], this.params).filter(([param]) => param != "qid");

      if (qid) {
        params.push(["qid", qid]);
      }

      this.$emit("update:params", params);
      this.$emit("onParamsChange");
    },
    resetEditors() {
      this.editors = 0;
      this.editorsSign = ">";
      const params = Object.assign([], this.params).filter(
        ([param]) => param != "editors" && param != "editors_sign"
      );
      params.push(["editors", 0], ["editors_sign", ">"]);
      this.$emit("update:params", params);
      this.$emit("onParamsChange");
    },
    resetEdits() {
      this.edits = 0;
      this.editsSign = ">";
      const params = Object.assign([], this.params).filter(
        ([param]) => param != "number_of_edits" && param != "edits_sign"
      );
      params.push(["number_of_edits", 0], ["edits_sign", ">"]);
      this.$emit("update:params", params);
      this.$emit("onParamsChange");
    },
    resetIndications() {
      this.indicationsCount = null;
      this.indicationsSign = "=";
      const params = Object.assign([], this.params).filter(
        ([param]) => param != "indications_count"
      );
      this.$emit("update:params", params);
      this.$emit("onParamsChange");
    },
    updateParams(paramName, paramValue) {
      const params = Object.assign([], this.params).filter(
        ([param]) => param != paramName
      );

      if (paramValue !== null) {
        params.push([paramName, paramValue]);
      }

      this.$emit("update:params", params);
      this.$emit("onParamsChange");
    },
  },
  watch: {
    projects() {
      this.options = [];

      for (const value in this.projects) {
        this.options.push({
          value,
          label: this.projects[value],
        });
      }
    },
    editorsSign(newSign) {
      this.updateParams("editors_sign", newSign);
    },
    editors(newEditors) {
      if(newEditors == -1){
        this.editors = 0;
        return
      }
      this.updateParams("editors", newEditors);
    },
    editsSign(newSign) {
      this.updateParams("edits_sign", newSign);
    },
    edits(newEdits) {
      if(newEdits == -1){
        this.edits = 0;
        return
      }
      this.updateParams("number_of_edits", newEdits);
    },
    indicationsSign(newSign) {
      this.updateParams("indications_sign", newSign);
    },
    indicationsCount(newIndications, oldIndications) {
      if(oldIndications == null && this.indicationsCount == null){
        this.indicationsCount = 1;
      }

      if(newIndications == 0){
        this.indicationsCount = 1;
        return
      }

      this.updateParams("indications_count", newIndications);
    },
    indications(newIndications) {
      this.updateParams("indications", newIndications);
    },
  },
  data() {
    const defaultOptions = [];
    const editors = ref(null);
    const editorsSign = ref(null);
    const edits = ref(0);
    const editsSign = ref(null);
    const indicationsSign = ref(null);
    const indicationsCount = ref(null);
    const indications = ref([]);

    this.params.forEach(([param, value]) => {
      if (param == "project") {
        defaultOptions.push(value);
      }
      switch (param) {
        case "editors_sign":
          editorsSign.value = value;
          break;
        case "editors":
          editors.value = value;
          break;
        case "edits_sign":
          editsSign.value = value;
          break;
        case "number_of_edits":
          edits.value = value;
          break;
        case "indications_sign":
          indicationsSign.value = value;
          break;
        case "indications_count":
          indicationsCount.value = value;
          break;
        case "indications":
          indications.value = value;
          break;
      }
    });

    return {
      defaultOptions,
      options: [],
      editors,
      editorsSign,
      edits,
      editsSign,
      indicationsSign,
      indicationsSigns: [
        {
          label: "=",
          value: "=",
        },
      ],
      indicationsCount,
      indications,
      signs: [
        {
          label: "<",
          value: "<",
        },
        {
          label: ">",
          value: ">",
        },
      ]
    };
  },
};
</script>

<style scoped>
.nav-bar-button {
  padding: 20px;
  font-size: 20px;
  border-width: 2px;
  border-style: solid;
  align-self: flex-end;
}

.nav-bar-extra {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: flex-end;
}

.nav-bar-extra > * {
  margin-left: 10px;
}

.nav-bar-select {
  min-width: 300px;
  width: auto;
  margin-right: 8px;
}
.nav-bar-short-select {
  min-width: 70px;
  width: auto;
  margin-right: 8px;
}
.beta {
  background-color: #1CAC78;
  font-weight: bold;
  font-size: 12px;
  padding: 5px 15px 5px 15px;
  border-radius: 5px;
  border-style: none;
  text-align: center;
  width: fit-content;
  margin-top: 5px;
}
.logo-container {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style>
