<template>
  <div class="conversation-list">


    <div class="logged_user_menu">
      <div class="user-info">

        <img
          @click="changePfp"
          class="mainpfp"
          style="margin-top: 0px;"
          :src="(userPfp && userPfpType) ? ('data:' + userPfpType + ';base64,' + userPfp) : 'https://i.imgur.com/D95gXlb.png'"
          @error="e => (e.target.src = 'https://i.imgur.com/D95gXlb.png')"
        />

        <input type="file" @change="handleFileUpload" accept="image/*" style="display: none;" ref="fileInput">

        <div style="display: flex; align-items: center; justify-items: flex-start;">

          <h2 style="">Hello,</h2>
          <h2 @click="changeUsernameInputField" class="changeName">{{ this.username }}</h2>

        </div>

      </div>

      <div class="logged_menu_buttons" v-if="!showUsernameInput">

        <button @click="newChat">New Chat</button>
        <button @click="handleLogout" class="logout-btn">Logout</button>


      </div>

      <div class="form" v-if="showUsernameInput">
      <form @submit.prevent="changeUsername">
        <div>
          <label>Insert a new username:</label>
          <br>
          <input v-model="newUsername" required type="text"/>
          <p style="color: red; font-weight: bold; margin-top: 0px;" v-if="newUsernameTaken">
            Username already taken!
          </p>
        </div>
        <button type="submit">Set New Name</button>
      </form>
      <button v-if="showUsernameInput" @click="cancelChangeName">Cancel</button>
      </div>


    </div>

    <ul>
      <p v-if="noMessages" style="text-align: center; color: #999999;">You have no conversations!</p>
      <li class="convoBubble"
        v-for="c in this.convertedConvos"
        :key="c.convoid"
        @click="selectConversation(c)"
        >

          <div class="convoBubbleLeft">

            <div style="display: flex; align-items: center; justify-items: flex-start;">

              <img class="pfp" v-if="c.chatPic != null && c.chatPic != 'https://i.imgur.com/D95gXlb.png' && c.chatPic != ''"
              :src="'data:' + c.chatPicType + ';base64,' + c.chatPic">


              <img
                class="pfp"
                :src="(c.chatPic && c.chatPicType) ? ('data:' + c.chatPicType + ';base64,' + c.chatPic) : 'https://i.imgur.com/D95gXlb.png'"
                @error="e => (e.target.src = 'https://i.imgur.com/D95gXlb.png')"
              />

              <h3 class='chatName'>{{ c.chatName }}</h3>

            </div>
            <div v-if="c.lastSender == this.username" style="display: flex; align-items: center;">


                <p style="color: teal" v-if="c.chatStatus == `seen`">🗸🗸</p>
                <p style="color: gray" v-if="c.chatStatus == `delivered`">🗸</p>
                <p style="padding-left: 10px">{{c.chatPreview}}</p>

            </div>
            <p v-else style="padding-left: 5px">
              {{c.chatPreview}}
            </p>
          </div>

          <div class="convoBubbleRight">
            <p class="notificationDot" v-if="!hasSeenLastMessage(c)"></p>
            <br v-else>

            <p>{{ c.chatTime.substring(11, 16) }}</p>

          </div>

      </li>
    </ul>  
  </div>
</template>

<script>
import api from '../api';
import { resizeImageIfNeeded } from '@/utils/resizeImage';

export default {
  name: 'ConversationList',

  props: {
    username: String,
    userPfp: String,
    userPfpType: String,
    newUsernameTaken: Boolean,
    reload: Boolean,
  },


  data() {
    return {
      conversations: [], //  usato per ricevere dati dall'api
      convertedConvos: [], //  usato dal ciclo for per renderizzare le chat
      updatedConvertedConvos: [], //  usata per auto-fetchare e confontare in live

      showUsernameInput: false,
      newUsername:'',

      noMessages: false,
      selectedFile: null,

      fileInput: null, //  Add a ref to the file input
    };
  },

  methods: {


    // NB: permette di inviare una foto per volta, e non in bulk (TODO)
    async handleFileUpload(event) {
      const file = event.target.files && event.target.files[0];
      if (!file) return;

      // Ridimensiona/Comprimi PFP
      const { file: resized } = await resizeImageIfNeeded(file, {
        maxWidth: 512,
        maxHeight: 512,
        maxBytes: 220 * 1024,   // ~220 KB
        outputMime: 'image/jpeg',
        quality: 0.9
      });

      this.$emit('changeUserPfp', resized);
      console.log("Profile picture selected (resized):", resized);
      // pulisci l'input per poter ricaricare lo stesso file se serve
      if (event.target) event.target.value = '';
    },

    cancelChangeName(){ this.showUsernameInput = false; this.$emit("resetNameError")},

    handleLogout() {
      this.$emit('logout');
    },

    selectConversation(convo) {
 
      // faccio sparire subito il notification dot
      if (convo.lastSender !== this.username) {
        convo.chatStatus = 'seen';
      }

      // emetto ad app.vue ID e convo struct 
      this.$emit('select-conversation', convo.convoid, convo);
    },


    //  Controlla se l'utente ha visto l'ultimo messaggio
    hasSeenLastMessage(convo) {
      // Se l'ultimo messaggio è stato inviato dall'utente corrente, è stato "visto" .
      if (convo.lastSender === this.username) {
      return true;
      }
      
      // Se non c'è un ultimo mittente => la chat è nuova e vuota. Nessun messaggio non letto.
      if (!convo.lastSender) {
      return true;
      }

      // Se lo stato del messaggio è 'seen', è stato visto.
      if (convo.chatStatus === 'seen') {
      return true;
      }
      
      // Se l'utente corrente è in 'seenBy', è stato visto.
      if (convo.seenBy && convo.seenBy.some(user => user.username === this.username)) {
      return true;
      }
      
      // else, l'ultimo messaggio non è stato visto.
      return false;
    },


    changeUsernameInputField(){
      console.log("Opening Username Field");
      this.showUsernameInput = true;
    },

    changeUsername(){
        console.log("Trying to change username to "+this.newUsername+"...");
        this.$emit("changeUsername", this.newUsername)
        
        //  chiudi interfaccia solo se non ho errore
        if(this.newUsernameTaken == true){

          console.log("Hai cambiato nome, chiudo interfaccia e ricarico chat...")
          this.showUsernameInput = false
          this.convertedConvos = []


        //  per 2 secondi faccio polling molto veloce,.
        let seconds = 2;
        const timer = setInterval(() => {
          console.log("Eseguito!");
          seconds--;
          this.pollingFetcher();
          if (seconds <= 0) clearInterval(timer);
        }, 500);

        }else{
          console.log("Errore con il nome, l'interfaccia rimane aperta")
        }
        
        
        


    },

    changePfp(){
      console.log("Trying to change the profile picture...")
      //  Programmatically click the file input
      this.$refs.fileInput.click();
    },

    newChat(){
      console.log("Trying to start a new chat..")
      this.$emit('newChat');
    },

    //  helper function for fetching
    arrayEquals(a, b) {
          if (a === b) return true; //  Se sono lo stesso oggetto, sono identici
          if (a == null || b == null) return false;
          if (a.length !== b.length) return false;

          for (let i = 0; i < a.length; i++) {
            if (a[i] !== b[i]) return false;
          }
          return true;
        },

    async pollingFetcher(){

      try {
        const response = await api.getConversations();
        const responseData = typeof response === 'string' ? JSON.parse(response) : response;

        //  Se esistono dati ricevuti e ricevo effettivamente un array di conversazioni...
        if (responseData && Array.isArray(responseData['User Conversations'])) {

          const uniqueConversations = [];
          const seenConvoIds = new Set(); //  Uso un Set per tenere traccia degli ID già visti

          if (responseData['User Conversations']) { //  Aggiungi un controllo null/undefined
            for (const convo of responseData['User Conversations']) {
              if (convo && convo.convoid && !seenConvoIds.has(convo.convoid)) {
                uniqueConversations.push(convo);
                seenConvoIds.add(convo.convoid);
              }
            }
          }

          this.conversations = uniqueConversations; //  Usa l'array filtrato

  

          this.updatedConvertedConvos = []; //  Reset the array before populating it

          for (let i = 0; i < this.conversations.length; i++) {
            var c = this.conversations[i];

            //  prima imposto i parametri che posso ricavare da convo
            var toRender = {
              convoid: c.convoid,
              chatPreview: c.preview,
              chatTime: c.datelastmessage,
              chatStatus: c.messages?.[c.messages.length - 1]?.status,
              lastSender: c.messages?.[c.messages.length - 1]?.author?.username,
              seenBy: c.messages?.[c.messages.length - 1]?.seenby || [],

              chatPic: null,
              chatName: null,
              chatPicType: null
            };

            //  poi vedo se è di gruppo o privata per il render corretto...
            var response2 = await api.getConvoInfo(c.convoid);

            if (response2?.data?.Group){
              toRender.chatPic   = response2.data.Group.groupphoto;
              toRender.chatPicType = response2.data.Group.photoMimeType;
              toRender.chatName = response2.data.Group.name;
            }
            else if(response2?.data?.PrivateConvo){
              var x = response2.data.PrivateConvo;

              var otherDude = x.firstuser ? (this.username != x.firstuser.username ? x.firstuser : x.seconduser) : x.seconduser;
              toRender.chatPic   = otherDude?.photo;
              toRender.chatPicType = otherDude?.photoMimeType;
              toRender.chatName = otherDude?.username;
            }

            this.updatedConvertedConvos.push(toRender);

          }

          //  Ordina l'array 'this.updatedConvertedConvos' in base a 'chatTime' (dal più recente al più vecchio)
          this.updatedConvertedConvos.sort((a, b) => new Date(b.chatTime).getTime() - new Date(a.chatTime).getTime());

          if (!this.arrayEquals(this.updatedConvertedConvos, this.convertedConvos)) // se ho trovato qualcosa di nuovo...
          {
            this.convertedConvos = this.updatedConvertedConvos
          }
          this.noMessages = false

        } else {
          console.error('Unexpected response format:', responseData);
          this.conversations = [];
          this.noMessages = true
        }

      } catch (error) {
        console.error('Error fetching conversations:', error);
        this.noMessages = true
        this.conversations = [];
      }

      //  ogni volta che ricalcolo i converted convos faccio un emit in modo che App.vue aggiorni la sua variabile.
      this.$emit("updateConvertedConvos", this.convertedConvos)

    },


    /*
    
      IMPORTANTE
      PERCHÈ
      AIUTA
      A
      SETTARE
      IL
      TEMPO
      DI
      RISPOSTA
      DEL
      SITO
    
    */
    startPolling() {
      //  Fetch conversations every 3 seconds (adjust the interval as needed)
      this.pollingInterval = setInterval(() => {
        console.log("Polling for new conversations...");
        this.pollingFetcher();
      }, 5000); //  3000ms = 3 second
    },

    stopPolling() {
      if (this.pollingInterval) {
        clearInterval(this.pollingInterval);
        this.pollingInterval = null;
        console.log("Polling stopped.");
      }
    },
  },

  watch: {

    reload(){
      this.pollingFetcher()
      this.$emit("reloaded")
    }
  },
 /* Appena carico la pagina recupera le conversazioni */
  mounted() {
    console.log("ConversationList component mounted!");
    this.pollingFetcher();
    this.startPolling();
  },
  beforeUnmount() {
    this.stopPolling(); //  Stop polling when the component is destroyed
  },

};
</script>

<style scoped>
.conversation-list {
  width: 30%;
  border-right: 1px solid #ccc;
  padding: 10px;
}
ul {
  list-style: none;
  padding: 0;
}
li {
  padding: 10px;
  cursor: pointer;
}
li:hover {
  background-color: #f0f0f0;
}

.logged_menu_buttons button:hover {
  background-color: #f0f0f0 !important;
}

.mainpfp {
  /* Make the element a square to ensure a perfect circle */
  width: 60px; /* Adjust the size as needed */
  height: 60px; /* Should be the same as the width */
  border-radius: 50%; /* This makes the element circular */
  overflow: hidden; /* Clips content that goes outside the circle */
  display: inline-block; /* Allows multiple profile pictures to sit on the same line */
  margin-top: 10px;
  border: 1px solid #ccc;
  object-fit: cover

}

.mainpfp:hover{

  border-color: lightseagreen;

}


.pfp {
  /* Make the element a square to ensure a perfect circle */
  width: 50px; /* Adjust the size as needed */
  height: 50px; /* Should be the same as the width */
  border-radius: 50%; /* This makes the element circular */
  overflow: hidden; /* Clips content that goes outside the circle */
  display: inline-block; /* Allows multiple profile pictures to sit on the same line */
  margin-top: 10px;
  border: 1px solid #ccc;
  object-fit: cover
}


.logged_user_menu {
  display: flex;
  flex-direction: column;
  align-items: center;
  border-bottom: #ccc 1px solid;
}

.user-info {
  display: flex;
  width: 100%;
  justify-content: space-evenly;
  align-items: center; /* Vertically align items in the user info section */
}

.logged_menu_buttons {
  display: flex;
  justify-content: space-evenly;
  width: 100%;
  padding: 10px;
  background-color: white;
  border-radius: 10px; /* Adjust this value to control the roundness */
  margin-left: 5px;
  margin-right: 5px;
}

.form {

  font-weight: bold;
  margin-top: 10px;
  border-top: #ccc 1px solid;
  padding: 0px 40px;
  padding-top: 10px;
  width: 100%;
}



input{

background-color: white;
font-weight: normal;
border-radius: 10px; /* Adjust this value to control the roundness */
border: #ccc 1px solid;
margin: 10px 0px;
text-align: center;
font-weight: normal;

}

button {

  background-color: white;
  padding: 15px 20px 15px 20px;
  font-weight: bold;
  border-radius: 10px; /* Adjust this value to control the roundness */
  border: #ccc 1px solid
}


.convoBubble {
  display: flex;
  border: 1px solid #ccc;
  margin: 0px 10px 10px;

  border-radius: 10px;
  padding: 0 15px; /* Add some padding inside the bubble */
  justify-content: space-between; /* Puts space between the left and right sections */
  align-items: center; /* Vertically aligns items in the center */
}


.convoBubbleRight {
  margin-left: 20px;
  display: flex;
  flex-direction: column; /* Stack the notification and time vertically */
  align-items: flex-end; /* Align items to the right */
}

.notificationDot {
  border: 3px;
  color: green;
  background-color: lightgreen;
  border-color: green;
  border-radius: 50%;
  padding: 5px; /* Adjust padding as needed */
  margin-bottom: 5px; /* Add some space between the dot and the time */
  width: 10px; /* Set a fixed width */
  height: 10px; /* Make height equal to width for a circle */

}

h3.chatName {
    justify-self: left;
    padding-left: 15px;
    margin-bottom: 0px;
    margin-top: 15px;

}

.chatTime {
  font-size: 0.8em; /* Adjust font size as needed */
  color: #777; /* Adjust color as needed */
}

/* Remove the default styling that might interfere */
.convoBubbleLeft {
  justify-content: left; /* Remove this */
}

.convoBubbleRight {
  justify-content: right; /* Remove this */
}

br[v-else] {
  display: none; /* Hide the <br> tag when the notification dot is visible */
}

.changeName:hover{
  color: lightseagreen;
}
.changeName{
  padding-left: 8px;
  text-decoration: underline;
  color: teal;
}


</style>
