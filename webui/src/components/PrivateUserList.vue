<template>
  <h3>Choose a user to start a conversation with:</h3>
  
  <div class="user-selection-container">
    <div
      v-for="user in userArray"
      :key="user.username"
      class="user-card"
      @click="selectUser(user)"
      :class="{ 'selected': this.selectedUsername && this.selectedUsername === user.username }"
    >
      
            <img
              class="pfp"
              v-if="user.photo != null && user.photo !== ''"
              :src="'data:' + user.photoMimeType + ';base64,' + user.photo"
              alt="Profile Picture" />
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
              :disabled="!this.selectedUsername"
              style="display:flex; align-self: center;" >
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
      };
    },

    mounted(){this.fetchUsers()},

    methods: {
        async fetchUsers(){
            try {
                const response = await api.listUsers();
                // Filter out the current user from the list
                this.userArray = response.data.Users.filter(user => user.username !== this.currentUser);
            } catch (error) {
                console.error("Error fetching users:", error);
            }
        },

        selectUser(user) {
          this.selectedUsername = user.username;
        },

        confirmSelection() {
          if (this.selectedUsername) {
            this.$emit('user-selected', this.selectedUsername);
          }
        },
    }
}
</script>

<style scoped>
.user-selection-container {
    display: flex;
    justify-content: center;
    align-items: flex-start;
    flex-wrap: wrap;
    gap: 15px;
    padding: 20px;
}

.user-card {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 8px;
    padding: 15px;
    border: 2px solid #eee;
    border-radius: 10px;
    cursor: pointer;
    transition: all 0.2s ease-in-out;
    width: 100px;
}

.user-card:hover {
    border-color: #ccc;
    background-color: #f9f9f9;
}

.user-card.selected {
  border-color: teal;
  background-color: #e0f7f7;
}

.pfp {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  border: 1px solid #ccc;
  object-fit: cover;
}

h3 {
    margin: 0;
    font-size: 0.9em;
    word-break: break-word;
    text-align: center;
}
</style>