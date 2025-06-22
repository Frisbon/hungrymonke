<template>
  <div id="app">
    <!-- Form di login se l'utente non è loggato -->

    
    <div v-if="!isLoggedIn" class="login">


      <h2>WasaText 🌟</h2>
      <!-- ascolto per l'evento @submit e impedisco al browser la funzionalità di default (ricarica pagina) usando .prevent -->
      <form @submit.prevent="handleLogin">
        <div>
          <label>Insert your username:</label>
          <br>
          <input v-model="username" type="text" required />
        </div>
        <button type="submit">Login</button>
      </form>

      <p v-if="loginError" class="error">{{ loginError }}</p>
    </div>

    <!-- Interfaccia chat se l'utente è loggato -->
    <div v-else class="container">
      
      <ConversationList 
        ref="convoList"

        :username="username" 
        :userPfp="userPfp" 
        :userPfpType="userPfpType"
        :newUsernameTaken = "newUsernameTaken"
        @select-conversation="selectConversation"
        @changeUserPfp = "changeUserPfp"
        @changeUsername="changeUsername"
        @logout = "handleLogout"
        @resetNameError = "resetNameError"
        @updateConvertedConvos = "updateConvertedConvos"
        @newChat = "newChatSequence"
      />

      <div class="chat-area" v-if="showMessageWindow">

        <ChatMessages 
          :selectedConvoID="selectedConvoID" 
          :username="username" 
          :isGroup="isGroup"
          :selectedConvoRender = "selectedConvoRender"
          @forwarder = "forwarder"
          @reloadConvo = "reloadChatMessages"  
          @groupLeft = "reloadButNoChoosing"
        />

      </div>

      <div class="convo-list-area" v-if="showConvoListWindow">

        <ForwardingConvoList
        :currentUser="username"
        :convertedConvos="convertedConvos"
        @forwardToConvo="forwardingConvoListHandler"
        />

      </div>

      <div class="new-chat-area" v-if="showNewChatWindow">
        
        <NewChat
          :currentUser="username"
          @selected-user="startPrivateConvo"
          @selected-users="startGroupChat"
        />

      </div>

    </div>
  </div>
</template>

<script>
import ConversationList from './components/ConversationList.vue';
import ChatMessages from './components/ChatMessages.vue';
import api from './api';
import ForwardingConvoList from './components/ForwardingConvoList.vue';
import NewChat from './components/NewChat.vue';

export default {
  name: 'App',
  components: {
    ConversationList,
    ChatMessages,
    ForwardingConvoList,
    NewChat
  },
  data() {
    return {
      isLoggedIn: false,
      username: '',
      userPfp: 'https://i.imgur.com/D95gXlb.png', // setto pfp di default
      userPfpType: null,
      loginError: '',
      selectedConvoID: null,

      isGroup: null,
      selectedConvoRender: null,
      convertedConvos: null,
      newUsernameTaken: false,

      selectedMessage: null,
      
      showMessageWindow: true,
      showConvoListWindow: false,
      showNewChatWindow: false,

      reloadConvos: false,
    };
  },
  methods: {

    resetReload(){ this.reloadConvos = false;},
    reloadChatMessages(){
      this.showMessageWindow = false;

      const convoListComponent = this.$refs.convoList;

      convoListComponent.pollingFetcher();
  
      // nel frattempo devo ri-fetchare nella convo list
      setTimeout(() => {convoListComponent.selectConversation(this.selectedConvoID); this.showMessageWindow = true;}, 400);
    },

    reloadButNoChoosing(){
      this.showMessageWindow = false;

      const convoListComponent = this.$refs.convoList;

      convoListComponent.pollingFetcher();

      // nel frattempo devo ri-fetchare nella convo list
      setTimeout(() => {this.selectedConvoID = null; this.selectedConvoRender = null; this.showMessageWindow = true;}, 400);
    },

    newChatSequence(){
      this.showNewChatWindow = true
      this.showMessageWindow = false
    },

    async startPrivateConvo(user){
      console.log("Sono in App.vue dentro startPrivateConvo(), mi connetto al back-end...")
      const response = await api.startPrivateConvo(user)
      this.selectedConvoID = null
      this.showMessageWindow = false
      this.showNewChatWindow = false
      console.log("response: ", response.data)
    },

    async startGroupChat(users, name, picture, mime){
      console.log("Sono in App.vue dentro startGroupChat(), mi connetto al back-end...")
      const response = await api.startGroupChat(users, name, picture, mime)
      this.selectedConvoID = null
      this.showMessageWindow = false
      this.showNewChatWindow = false
      console.log("response: ", response.data)
    },


    async handleLogin() {
      try {
        const credentials = this.username;
        const response = await api.login(credentials);
        console.log('Login successful:', response.data);

        this.isLoggedIn = true;
        this.loginError = '';
        if (response.data.user.photo != null){this.userPfp = response.data.user.photo; this.userPfpType = response.data.user.photoMimeType;}

        //  salva anche lo username e foto profilo
        localStorage.setItem('username', this.username);
        localStorage.setItem('userPfp', this.userPfp);
        localStorage.setItem('userPfpType', this.userPfpType);
        

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
      this.recipientUsername = '';
      this.userPfp = 'https://i.imgur.com/D95gXlb.png';
      this.userPfpType = null;
      localStorage.removeItem('username');
      localStorage.removeItem('userPfp');
      localStorage.removeItem('userPfpType');
      this.isGroup = null;
      this.selectedConvoRender = null;
      this.newUsernameTaken = false

    },  

    async selectConversation(convoID, convoRender) {
      this.selectedConvoID = convoID;
      this.selectedConvoRender = convoRender;
      this.showMessageWindow = true
      this.showNewChatWindow = false
      this.showConvoListWindow = false
      console.log("Hai selezionato la chat, ecco il render:");
      console.log(convoRender)

      const response = await api.getConvoInfo(convoID);
      if (response?.data?.Group) {
        this.isGroup = true;
      } else {
        this.isGroup = false;
      }
    },

    resetNameError(){ this.newUsernameTaken = false},
    updateConvertedConvos(x){this.convertedConvos = x},

    async forwardingConvoListHandler(chosenConvoObject){ 
      this.showConvoListWindow = false; 
      this.showMessageWindow = true
      console.log("Sono in App.vue dentro forwardingConvoListHandler(), mi connetto al back-end...")
      console.log("prima dell'api.js, al momento mi hai passato l'oggetto: ", chosenConvoObject)
      const response = await api.forwardMessage(chosenConvoObject.convoid, this.selectedMessage)
      console.log("Response: ")
      console.log(response.data)
      this.selectedConvoID = chosenConvoObject.convoid
      console.log(this.showConvoListWindow, this.showMessageWindow)
    },
    
    forwarder(selectedMsg){
      this.showConvoListWindow = true; 
      this.showMessageWindow = false;
      this.selectedMessage = selectedMsg
      console.log(this.showConvoListWindow, this.showMessageWindow)

    },

    async changeUsername(newName){
      console.log("Sono in App.vue dentro changeUsername(), mi connetto al back-end...")
      const response = await api.changeUsername(newName)
      if (!response.error){

        console.log()
        this.username = response.user.username
        this.userPfp = response.user.photo
        this.userPfpType = response.user.photoMimeType

        // resetto tutto!
        localStorage.removeItem('username');
        localStorage.removeItem('userPfp');
        localStorage.removeItem('userPfpType');
        this.selectedConvoID = '';
        localStorage.setItem('username', this.username);
        localStorage.setItem('userPfp', this.userPfp);
        localStorage.setItem('userPfpType', this.userPfpType);
        this.newUsernameTaken = false

      }
      else{
        this.newUsernameTaken = true
      }
      

    },

    async changeUserPfp(newPfp){
      console.log("Sono in App.vue dentro changeUserPfp(), mi connetto al back-end...")
      const response = await api.changeUserPfp(newPfp)
      if (!response.error){
        this.userPfp = response.user.photo
        this.userPfpType = response.user.photoMimeType

        // resetto tutto!
        localStorage.removeItem('userPfp');
        localStorage.removeItem('userPfpType');
        this.selectedConvoID = '';
        localStorage.setItem('userPfp', this.userPfp);
        localStorage.setItem('userPfpType', this.userPfpType);
        console.log("La pfp dovrebbe essere stata cambiata!")
      }


    }

  
  },

  /*Appena carico il DOM, questo sarà il primo ad essere eseguito*/ 
  mounted() {
    //  Controlla se l'utente è già loggato (token e username salvato)
    if (localStorage.getItem('token')) {
      this.isLoggedIn = true;
      this.username = localStorage.getItem('username') || '';
      this.userPfp = localStorage.getItem('userPfp') || 'https://i.imgur.com/D95gXlb.png';
      this.userPfpType = localStorage.getItem('userPfpType') || null;
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
  
}

#app .login{ justify-items: center; }

.container {
  display: flex;
  max-width: 1200px;
  margin: 0 auto;
  height: fit-content;
}

.login {
  display: flex;
  flex-direction: column;
  background-color: white;
  font-weight: bold;
  border-radius: 10px; /* Adjust this value to control the roundness */
  border: #ccc 1px solid;
  max-height: fit-content;
  max-width: fit-content;
  justify-self: center;
  padding: 0px 20px;
}

.login h2{ margin: 10px 0px;}

.login form {

  
  border-top: #ccc 1px solid;
  padding: 0px 40px;
  padding-top: 10px;
}

#app button{

  background-color: white;
  font-weight: bold;
  border-radius: 10px; /* Adjust this value to control the roundness */
  border: #ccc 1px solid;
  margin-bottom: 10px;
  padding: 15px 20px 15px 20px;
}

#app button :hover{

  background-color: #f0f0f0;

}

#app button :disabled{
  background-color: #ccc;
  cursor: not-allowed;
}

.login input{

background-color: white;
font-weight: bold;
border-radius: 10px; /* Adjust this value to control the roundness */
border: #ccc 1px solid;
margin: 10px 0px;
text-align: center;
font-weight: normal;

}

.error {
  color: red;
}


.conversation-list {
  width: 30%; /* Adjust as needed */
  border-right: 1px solid #ccc;
  padding: 10px;
  padding-right: 20px;
  height: 100%;
  overflow-y: auto;
}
.chat-area {
  flex-grow: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
}

.convo-list-area {
    display: flex;
    flex-grow: 1;
    flex-direction: column;
    align-content: center;
    justify-content: center;
    overflow-y: hidden;
}

.new-chat-area {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;

}

</style>
