<template>
  <div v-if="userArray.length > 0" class="group-user-list-container">
    <h3>Seleziona utenti da aggiungere al gruppo</h3>

    <div class="user-selection-area">
      <div
        v-for="user in userArray"
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
        <img
          class="pfp"
          v-else
          src="https://i.imgur.com/D95gXlb.png"
          alt="Default PFP"
        />
        <h3>{{ user.username }}</h3>
      </div>
    </div>

    <button
      type="button"
      @click="confirmSelection"
      :disabled="selectedUsernames.length === 0"
    >
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
      convoID: String
    },

    data() {
      return {
        userArray: [],
        selectedUsernames: [],
      };
    },

    mounted(){
        this.fetchUsers();
    },

    methods: {
        close(){this.$emit("close-buttons")},

        async fetchUsers(){
            this.userArray = [];
            this.selectedUsernames = [];

            try {
                const allUsersResponse = await api.listUsers();
                const allUsers = allUsersResponse.data.Users || [];
                
                let existingGroupMembersUsernames = [];
                if (this.convoID) {
                  try {
                       const groupInfoResponse = await api.getConvoInfo(this.convoID); 
                       const groupUsers = groupInfoResponse.data.Group.users || []; 
                       existingGroupMembersUsernames = groupUsers.map(user => user.username);
                  } catch (error) {
                      console.error(`Errore nel fetch dei membri del gruppo ${this.convoID}:`, error);
                  }
                }

                this.userArray = allUsers.filter(user =>
                    user.username !== this.currentUser &&
                    !existingGroupMembersUsernames.includes(user.username)
                );
            } catch (error) {
                console.error("Errore nel fetch o nel filtro degli utenti:", error);
                 this.userArray = [];
            } 
        },

        toggleUserSelection(user) {
            const username = user.username;
            const index = this.selectedUsernames.indexOf(username);
            if (index > -1) {
                this.selectedUsernames.splice(index, 1);
            } else {
                this.selectedUsernames.push(username);
            }
        },

        confirmSelection() {
            if (this.selectedUsernames.length > 0) {
                this.$emit('users-selected', this.selectedUsernames);
                this.selectedUsernames = [];
            }
        },
    }
}
</script>

<style scoped>
.group-user-list-container {
    padding: 20px;
    border: 1px solid #ccc;
    border-radius: 8px;
    text-align: center;
    width: 80%;
    max-width: 500px;
}

.user-selection-area {
    display: flex;
    justify-content: center;
    align-items: flex-start;
    flex-wrap: wrap;
    gap: 15px;
    padding: 10px 0;
    margin: 15px 0;
    max-height: 250px; /* Control height and enable scrolling */
    overflow-y: auto;
}

.user-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 10px;
    border: 2px solid transparent;
    border-radius: 10px;
    cursor: pointer;
    transition: all 0.2s ease-in-out;
    width: 90px;
}

.user-card:hover {
    background-color: #f9f9f9;
}

.user-card.selected {
  border-color: teal;
  background-color: #e0f7f7;
}

.pfp {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  border: 1px solid #ccc;
  object-fit: cover;
}

h3 {
    margin: 0 0 10px 0;
    font-size: 1.1em;
}

.user-card h3 {
    font-size: 0.8em;
    word-break: break-word;
}

button {
  margin-top: 10px;
}
</style>