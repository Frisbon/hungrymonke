<template>
  <div id="app">
    <!-- Form di login se l'utente non è loggato -->

    
    <div v-if="!isLoggedIn" class="login">


      <h2>WasaText PORCODIO 9 🌟</h2>
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
        :reload="reloadConvos"
        @reloaded="listReloaded"
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
          ref="chat"
          v-if="showMessageWindow"
          :selectedConvoID="selectedConvoID"
          :selectedConvoRender="selectedConvoRender"
          :username="username"
          :userPfpType="userPfpType"
          :userPfp="userPfp"
          :isGroup="isGroup"
          @newProfilePicture="(file) => setPfp(file)"
          @newUsername="(newUsername) => setUsername(newUsername)"
          @closeChat="() => closeChat()"
          :key="selectedConvoID"
          @forwarder="forwarder"
        />




      </div>

      <div class="convo-list-area" v-if="showConvoListWindow">

        <ForwardingConvoList
          v-if="showForwardingConvoListWindow"
          :convertedConvos="convertedConvos"
          :currentUser="username"
          @forwardToConvo="forwardingConvoListHandler"
          @close="showForwardingConvoListWindow = false"
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
import defaultPfp from '@/components/blank_pfp.png';
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
      userPfp: defaultPfp, // setto pfp di default  
      userPfpType: null,
      loginError: '',
      selectedConvoID: '',

      isGroup: null,
      selectedConvoRender: null,
      convertedConvos: null,
      newUsernameTaken: false,
      currentConvoTime: null,

      selectedMessage: null,
      
      showMessageWindow: true,
      showConvoListWindow: false,
      showNewChatWindow: false,

      reloadConvos: false,

      showForwardingConvoListWindow: false,



    };
  },
  methods: {

    async syncSelectedChatTitle() {
      try {
        if (!this.selectedConvoID) return;

        const { data } = await api.getConvoInfo(this.selectedConvoID);
        const convo = data?.conversation;
        if (!convo) return;

        let newName = '';
        if (convo.group) {
          // chat di gruppo
          newName = convo.group.name;
        } else {
          // 1-to-1: prendi "l’altro" utente rispetto a me
          const other = (this.username !== convo.firstuser.username)
            ? convo.firstuser
            : convo.seconduser;
          newName = other?.username || '';
        }

        if (newName && this.selectedConvoRender && this.selectedConvoRender.chatName !== newName) {
          // aggiorna il titolo visibile
          this.selectedConvoRender.chatName = newName;

          // opzionale ma utile: rinfresca la sidebar
          this.reloadConvos = true;
          this.$nextTick(() => (this.reloadConvos = false));
        }
      } catch (e) {
        // no-op: non vogliamo sporcare la console con polling silenzioso
      }
    },

    fetchMessages(convoID) {
      if (!convoID) return;
      const cmp = this.$refs.chat;
      if (cmp && typeof cmp.fetchMessages === 'function') {
        return cmp.fetchMessages(convoID);
      }
    },


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
      setTimeout(() => {this.selectedConvoID = ''; this.selectedConvoRender = null; this.showMessageWindow = true;}, 400);
    },

    newChatSequence(){
      this.showNewChatWindow = true
      this.showMessageWindow = false
    },

    async startPrivateConvo(user){
      console.log("Sono in App.vue dentro startPrivateConvo(), mi connetto al back-end...")
      const response = await api.startPrivateConvo(user)
      this.selectedConvoID = ''
      this.showMessageWindow = false
      this.showNewChatWindow = false
      console.log("response: ", response.data)
    },

    async startGroupChat(users, name, picture, mime){
      console.log("Sono in App.vue dentro startGroupChat(), mi connetto al back-end...")
      const response = await api.startGroupChat(users, name, picture, mime)
      this.selectedConvoID = ''
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

      async forwardMessage(convoID, selectedMessage) {
    try {
      const res = await api.forwardMessage(convoID, selectedMessage);
      console.log('[Forward] OK:', res?.data);

      // chiudi overlay e ricarica i messaggi della chat corrente
      this.showForwardingConvoListWindow = false;
      await this.fetchMessages(this.selectedConvoID);
      // pulizia selezione
      this.selectedMessage = null;
    } catch (err) {
      console.error('[Forward] API error:', err);
      throw err;
    }
  },


    async selectConversation(convoID, convoRender) {
      this.selectedConvoID = convoID;
      this.selectedConvoRender = convoRender;

      this.currentConvoTime = convoRender?.chatTime || null;

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

       if (!this.isGroup && this.showGroupInfo) {
        this.showGroupInfo = false;
      }

    },

    resetNameError(){ this.newUsernameTaken = false},
    updateConvertedConvos(val) {
      const prevList = this.convertedConvos || [];
      this.convertedConvos = val || [];

      // Trova la conversazione attualmente aperta nella lista aggiornata
      const sel = this.selectedConvoID
        ? this.convertedConvos.find(c => c.convoid === this.selectedConvoID)
        : null;

      // 1) Mantieni il titolo e l'immagine in sync anche se è cambiato "solo" il nome/pfp
      if (
        sel &&
        this.selectedConvoRender &&
        this.selectedConvoRender.convoid === sel.convoid
      ) {
        const needHeaderUpdate =
          this.selectedConvoRender.chatName !== sel.chatName ||
          this.selectedConvoRender.chatPic !== sel.chatPic ||
          this.selectedConvoRender.chatPicType !== sel.chatPicType ||
          this.selectedConvoRender.chatTime !== sel.chatTime;

        if (needHeaderUpdate) {
          // aggiorna i metadati usati nel titolo/header
          this.selectedConvoRender = {
            ...this.selectedConvoRender,
            chatName: sel.chatName,
            chatPic: sel.chatPic,
            chatPicType: sel.chatPicType,
            chatTime: sel.chatTime,
          };
        }
      }

      // 2) Se ci sono nuovi messaggi (chatTime è cambiato), aggiorna la finestra messaggi
      if (sel && this.selectedConvoRender && this.selectedConvoRender.convoid === sel.convoid) {
        const prev = prevList.find(c => c.convoid === this.selectedConvoID);
        if (!prev || prev.chatTime !== sel.chatTime) {
          this.fetchMessages(this.selectedConvoID, false, false, true);
        }
      }
    },


    async forwardingConvoListHandler(chosenConvo) {
      try {
        // 1) inoltra
        await this.forwardMessage(chosenConvo.convoid, this.selectedMessage);

        // 2) imposta subito il render target (evita undefined in template)
        this.selectedConvoID = chosenConvo.convoid;
        this.selectedConvoRender = { ...chosenConvo };
        this.isGroup = !!(chosenConvo.convoid && chosenConvo.convoid.startsWith('group_'));

        // 3) chiudi overlay e riapri la chat centrale
        this.showForwardingConvoListWindow = false;
        this.showMessageWindow = true;

        // 4) carica i messaggi della nuova chat
        await this.fetchMessages(this.selectedConvoID);

        // 5) cleanup
        this.selectedMessage = null;
      } catch (err) {
        console.error('[Forward] switch to convo failed:', err);
      }
    }
    ,
    
    forwarder(message) {
      // Mem messaggio selezionato
      this.selectedMessage = message;

      // contenitore sia montato
      this.showConvoListWindow = true;
      this.showMessageWindow = false;    
      this.showNewChatWindow = false;    

      // overlay di inoltro
      this.showForwardingConvoListWindow = true;
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

        // Dopo aver aggiornato localStorage ecc.
        if (this.selectedConvoRender && this.selectedConvoRender.username === this.username) {
          this.selectedConvoRender.username = newName;
        }



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

    this._titleSync = setInterval(() => {
      this.syncSelectedChatTitle();
    }, 2500);
  },
  beforeUnmount() {
    if (this._titleSync) clearInterval(this._titleSync);
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
