<script setup>
import { ref } from "vue";

const props = defineProps({
  busy: { type: Boolean, default: false },
  authed: { type: Boolean, default: false },
  teams: { type: Array, default: () => [] },
});

const emit = defineEmits(["create-team", "join-team"]);

const createTeamForm = ref({ name: "", description: "" });
const joinTeamId = ref("");

function createTeam() {
  emit("create-team", { ...createTeamForm.value }, () => {
    createTeamForm.value = { name: "", description: "" };
  });
}

function joinTeam() {
  emit("join-team", joinTeamId.value.trim());
}
</script>

<template>
  <section class="card">
    <h3>团队操作</h3>
    <div class="form-grid">
      <input v-model="createTeamForm.name" placeholder="团队名称" />
      <textarea
        v-model="createTeamForm.description"
        rows="3"
        placeholder="团队描述"
      />
      <button :disabled="props.busy || !props.authed" @click="createTeam">
        创建团队
      </button>
    </div>

    <div class="inline-tools">
      <input v-model="joinTeamId" placeholder="输入 team_id 加入" />
      <button :disabled="props.busy || !props.authed" @click="joinTeam">
        加入团队
      </button>
    </div>

    <div class="mini-list" v-if="props.teams.length">
      <p class="muted">最近团队</p>
      <button
        v-for="team in props.teams.slice(0, 6)"
        :key="team.team_id"
        class="chip"
        type="button"
      >
        {{ team.name }} ({{ team.team_id }})
      </button>
    </div>
  </section>
</template>
