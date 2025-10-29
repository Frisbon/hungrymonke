<template>
  <h3>Choose a user to start a conversation with:</h3>

  <!-- SEARCH -->
  <input
    class="user-search"
    type="text"
    v-model="searchTerm"
    placeholder="Search users…"
    aria-label="Search users"
  />

  <div class="user-selection-container">
    <div
      v-for="user in filteredUsers"
      :key="user.username"
      class="user-card"
      @click="selectUser(user)"
      :class="{ 'selected': selectedUsername && selectedUsername === user.username }"
    >
      <img
        class="pfp"
        :src="(user.photo && user.photoMimeType) ? ('data:' + user.photoMimeType + ';base64,' + user.photo) : 'https://i.imgur.com/D95gXlb.png'"
        @error="e => (e.target.src = 'https://i.imgur.com/D95gXlb.png')"
      />
      <h3>{{ user.username }}</h3>
    </div>
  </div>

  <!-- NO RESULTS -->
  <p v-if="searchTerm && filteredUsers.length === 0" class="no-results">
    no users found with your query brother
  </p>

  <button
    type="button"
    @click="confirmSelection"
    :disabled="!selectedUsername"
    style="display:flex; align-self: center;"
  >
    Conferma Utente
  </button>
</template>

<script>
import api from '@/api';

export default {
  name: 'PrivateUserList',

  props: {
    currentUser: String,
  },

  data() {
    return {
      userArray: [],
      selectedUsername: null,
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
    async fetchUsers() {
      try {
        const response = await api.listUsers();
        this.userArray = (response.data.Users || []).filter(u => u.username !== this.currentUser);
      } catch (error) {
        console.error('Error fetching users:', error);
      }
    },

    selectUser(user) {
      this.selectedUsername = user.username;
    },

    confirmSelection() {
      if (this.selectedUsername) this.$emit('user-selected', this.selectedUsername);
    },
  },
};
</script>

<style scoped>
.user-search{
  display:block;
  width: 100%;
  max-width: 420px;
  margin: 8px auto 12px;
  padding: 8px 12px;
  border: 1px solid #e2e8f0;
  border-radius: 10px;
  outline: none;
}
.user-search:focus{
  box-shadow: 0 0 0 3px #e5e7eb inset;
  border-color: #cbd5e1;
}

.user-selection-container{
  display:flex;
  justify-content:center;
  align-items:flex-start;
  gap:12px;
  flex-wrap:wrap;
}

.user-card{
  width:150px;
  padding:10px 14px;
  border:1px solid #e2e8f0;
  border-radius:12px;
  text-align:center;
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
  border:1px solid #ccc; object-fit:cover; margin-bottom:6px;
}

h3{
  margin:0; font-size:.9em; word-break:break-word; text-align:center;
}

.no-results{
  text-align:center; margin-top:6px; color:#6b7280; /* gray-500 */
}
</style>
