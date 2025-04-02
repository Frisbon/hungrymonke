<template>
  <div id="app">
    <!-- Form di login se l'utente non è loggato -->
    <div v-if="!isLoggedIn" class="login">
      <h2>Login</h2>
      <!-- ascolto per l'evento @submit e impedisco al browser la funzionalità di default (ricarica pagina) usando .prevent -->
      <form @submit.prevent="handleLogin">
        <div>
          <label>Username:</label>
          <input v-model="username" type="text" required />
        </div>
        <button type="submit">Login</button>
        <p v-if="loginError" class="error">{{ loginError }}</p>
      </form>
    </div>

    <!-- Interfaccia chat se l'utente è loggato -->
    <div v-else class="container">
      <button @click="handleLogout" class="logout-btn">Logout</button>
      
      <ConversationList :username="username" @select-conversation="selectConversation"/>

      <div class="chat-area">
      <ChatMessages :selectedConvoID="selectedConvoID" :recipientUsername="recipientUsername"/>
      </div>
      
    </div>
  </div>
</template>

<script>
import ConversationList from './components/ConversationList.vue';
import ChatMessages from './components/ChatMessages.vue';
import api from './api';

export default {
  name: 'App',
  components: {
    ConversationList,
    ChatMessages,
  },
  data() {
    return {
      isLoggedIn: false,
      username: '',
      loginError: '',
      selectedConvoID: null,
      recipientUsername: "",
    };
  },
  methods: {
    async handleLogin() {
      try {
        const credentials = this.username;
        const response = await api.login(credentials);
        console.log('Login successful:', response.data);
        
        this.isLoggedIn = true;
        this.loginError = '';
      } catch (error) {
        console.error('Login failed:', error);
        this.loginError = error.response?.data?.error || 'Login failed';
      }
    },

    handleLogout() {
      api.logout();
      this.isLoggedIn = false;
      this.username = '';
      this.selectedConvoID = '';
      this.recipientUsername= '';
    },

    selectConversation(convoID) {
      this.selectedConvoID = convoID;
      console.log("Hai selezionato la chat con ID: "+convoID)
      // ora renderizzo i messaggi?
    },

    //TODO: NewConversation(), con recipientUsername...

  },

  /*Appena carico il DOM, questo sarà il primo ad essere eseguito*/ 
  mounted() {
    // Controlla se l'utente è già loggato (token salvato)
    if (localStorage.getItem('token')) {
      this.isLoggedIn = true;
    }
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  text-align: center;
  color: #2c3e50;
  margin-top: 60px;
  height: 100vh; /* Make the app container take full viewport height */
  overflow: hidden; /* Prevent scrolling on the main container */
}
.container {
  display: flex;
  max-width: 1200px;
  margin: 0 auto;
  height: 100%; /* Make the chat container take full height */
}
.login {
  max-width: 400px;
  margin: 0 auto;
  padding: 20px;
}
.login div {
  margin-bottom: 10px;
}
.login label {
  display: inline-block;
  width: 100px;
}
.login input {
  padding: 5px;
  width: 200px;
}
.login button {
  padding: 5px 10px;
}
.error {
  color: red;
}
.logout-btn {
  position: absolute;
  top: 10px;
  right: 10px;
  background-color: #e74c3c;
  color: white;
  padding: 8px 12px;
  border: none;
  cursor: pointer;
}
.logout-btn:hover {
  background-color: #c0392b;
}
.conversation-list {
  width: 30%; /* Adjust as needed */
  border-right: 1px solid #ccc;
  padding: 10px;
  height: 100%;
  overflow-y: auto;
}
.chat-area {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
}
</style>
