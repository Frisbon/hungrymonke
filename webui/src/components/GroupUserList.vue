<template>
  <div v-if="userArray.length > 0" class="group-user-list-container">
    <h3>Seleziona utenti da aggiungere al gruppo</h3>

    <!-- SEARCH -->
    <input
      class="user-search"
      type="text"
      v-model="searchTerm"
      placeholder="Search users…"
      aria-label="Search users"
    />

    <div class="user-selection-area">
      <div
        v-for="user in filteredUsers"
        :key="user.username"
        class="user-card"
        @click="toggleUserSelection(user)"
        :class="{ 'selected': selectedUsernames.includes(user.username) }"
      >
        <img
          class="pfp"
          v-if="user.photo != null && user.photo !== ''"
          :src="'data:' + user.photoMimeType + ';base64,' + user.photo"
          alt="Profile Picture"
        />
        <img class="pfp" v-else src="https://i.imgur.com/D95gXlb.png" alt="Default PFP" />
        <h3>{{ user.username }}</h3>
      </div>
    </div>

    <!-- NO RESULTS -->
    <p v-if="searchTerm && filteredUsers.length === 0" class="no-results">
      no users found with your query brother
    </p>

    <button type="button" @click="confirmSelection" :disabled="selectedUsernames.length === 0">
      Conferma Selezione
    </button>
  </div>

  <div v-else class="group-user-list-container">
    <h3>Questo gruppo contiene tutti gli utenti possibili! 🤯</h3>
    <button @click="close">Close</button>
  </div>
</template>

<script>
import api from '@/api';

export default {
  name: 'GroupUserList',

  props: {
    currentUser: String,
    convoID: String, // se presente, esclude utenti già nel gruppo
  },

  data() {
    return {
      userArray: [],
      selectedUsernames: [],
      searchTerm: '',
    };
  },

  computed: {
    filteredUsers() {
      const q = this.searchTerm.trim().toLowerCase();
      if (!q) return this.userArray;
      return this.userArray.filter(u => (u.username || '').toLowerCase().startsWith(q));
    },
  },

  mounted() { this.fetchUsers(); },

  methods: {
    close() { this.$emit('close-buttons'); },

    async fetchUsers() {
      this.userArray = [];
      this.selectedUsernames = [];
      try {
        const allUsersResponse = await api.listUsers();
        const allUsers = allUsersResponse.data.Users || [];

        let existing = [];
        if (this.convoID) {
          try {
            const groupInfoResponse = await api.getConvoInfo(this.convoID);
            const groupUsers = groupInfoResponse.data.Group.users || [];
            existing = groupUsers.map(u => u.username);
          } catch (error) {
            console.error(`Errore nel fetch dei membri del gruppo ${this.convoID}:`, error);
          }
        }

        this.userArray = allUsers.filter(u =>
          u.username !== this.currentUser && !existing.includes(u.username)
        );
      } catch (error) {
        console.error('Errore nel fetch o nel filtro degli utenti:', error);
        this.userArray = [];
      }
    },

    toggleUserSelection(user) {
      const username = user.username;
      const i = this.selectedUsernames.indexOf(username);
      if (i > -1) this.selectedUsernames.splice(i, 1);
      else this.selectedUsernames.push(username);
    },

    confirmSelection() {
      if (this.selectedUsernames.length > 0) {
        this.$emit('users-selected', this.selectedUsernames);
        this.selectedUsernames = [];
      }
    },
  },
};
</script>

<style scoped>
.group-user-list-container{
  padding:20px;
  border:1px solid #ccc;
  border-radius:8px;
  text-align:center;
  width:80%;
  max-width:500px;
  margin:0 auto;
}
.user-search{
  display:block;
  width:100%;
  max-width:420px;
  margin:8px auto 12px;
  padding:8px 12px;
  border:1px solid #e2e8f0;
  border-radius:10px;
  outline:none;
}
.user-search:focus{
  box-shadow:0 0 0 3px #e5e7eb inset;
  border-color:#cbd5e1;
}

.user-selection-area{
  display:flex;
  gap:12px;
  flex-wrap:wrap;
  justify-content:center;
  align-items:flex-start;
}

.user-card{
  display:flex;
  flex-direction:column;
  align-items:center;
  gap:6px;
  width:150px;
  padding:10px 14px;
  border:1px solid #e2e8f0;
  border-radius:12px;
  background:#fff;
  cursor:pointer;
  transition:background-color .15s ease, box-shadow .15s ease, border-color .15s ease, transform .02s ease;
}

.user-card:hover,
.user-card.selected{
  background:#f3f4f6;
  border-color:#cbd5e1;
  box-shadow:0 0 0 3px #e5e7eb inset;
}
.user-card:active{ transform: translateY(1px); }

.pfp{
  width:60px; height:60px; border-radius:50%;
  overflow:hidden; display:inline-block; border:1px solid #ccc; object-fit:cover;
}

h3{ margin:0 0 10px 0; font-size:0.9em; word-break:break-word; }
.no-results{ text-align:center; margin-top:6px; color:#6b7280; }
button{ margin-top:10px; }
</style>
