<template>
  <div class="chat-messages-container">
    <!-- If no chat is selected -->
    <div v-if="!selectedConvoID" class="no-conversation">
      <p>Select a conversation to view messages.</p>
    </div>

    <!-- If a chat is selected -->
    <div v-else class="conversation-active">

      <div class="chatMenu">
          
        <img 
        class="chatpfp" style="margin-top: 0px;"
        v-if="this.selectedConvoRender.chatPic != null" 
        :src="'data:' + this.selectedConvoRender.chatPicType + ';base64,' + this.selectedConvoRender.chatPic">

        <img @click="changePfp" class="chatpfp" v-else :src="'https://i.imgur.com/D95gXlb.png'">
        
        <h1>{{ this.selectedConvoRender.chatName }}</h1>

        
        <button :disabled="!this.isGroup" @click="groupchatOptions">
          ...
        </button>

      </div>



      <div class="message-list"> 

          <div class="message-bubble"
          v-for="message in messages"
          :key="message.msgID"
          :class="{'other-message': message.author.username !== username,
                   'my-message': message.author.username === username}">
            
            <div v-if="message.author.username !== this.username && this.isGroup">
            <strong class="author-name">{{ message.author.username }}:</strong> 
            <!-- ADD FORWARDED OR REPLIED CONTENT HERE-->
            <br>
            </div>

            {{ message.content.text }}
            <br v-if="message.content.photo && message.content.text">
            
            <img class="sent-img" v-if="message.content.photo" 
            :src="'data:' + message.content.photoMimeType + ';base64,' 
            + message.content.photo" alt="Immagine allegata">
          
            <!-- ADD REACTIONS, TIME AND SEEN STATUS HERE-->
            <div class='messageStats'>
              <!-- la reaction bubble avrà foto in minuscolo e reaction accanto tutto accerchiato alla telegram-->
              <div class="reactionBubble" v-for="reaction in message.reactions" v-bind:key="reaction.author.username">
                <img class="reactionImage" :src="'data:undefined;base64,'+ reaction.author.photo" >
                <div class="reactionEmoticon">{{ reaction.emoticon }}</div>
              </div>

              <div class ="timestampAndStatus">
                
                <p class="noParagraph">{{ message.timestamp.substring(11, 16) }}</p>

                <p class="noParagraph" style="color: teal; margin: 0px 5px" v-if="message.status == `seen` && message.author.username == this.username">🗸🗸</p>
                <p class="noParagraph" style="color: gray; margin: 0px 5px" v-if="message.status == `delivered` && message.author.username == this.username">🗸</p>
                
                
             </div>
             <!-- Se sono io, vedo status time e reactions, altrimenti solo time e reactions-->
            </div>
          </div>
      </div>

      <form @submit.prevent="sendMessage" class="message-input-form">
        <div style="display: flex">
          <input type="file" @change="handleFileUpload" accept="image/*" />
          <button type="submit" :disabled="isSubmitDisabled">
            <strong>→</strong>
          </button>
        </div>  
        <input class="textInput" v-model="newMessage" placeholder="Type a message..." />
        
      </form>

    </div>

    
  </div>
</template>

<script>
import api from '../api';

export default {
  name: 'ChatMessages',

  props: {
    selectedConvoID: String,
    username: String,
    isGroup: Boolean,
    selectedConvoRender: Object,
  },

  data() {
    return {
      messages: [],
      newMessage: '',
      selectedFile: null,
      base64Image: '', // Per memorizzare la foto
    };
  },

  computed: {
    isSubmitDisabled() {
      return !this.newMessage && !this.base64Image;
    }
  },

  watch: {
    selectedConvoID(newConvoID) {
      if (newConvoID) {
        this.fetchMessages(newConvoID);
      }
    },
  },


  methods: {

    //NB: permette di inviare una foto per volta, e non in bulk (TODO)
    handleFileUpload(event) {
      this.selectedFile = event.target.files[0];
      if (this.selectedFile) {
        this.convertToBase64();
      } else {
        this.base64Image = '';
      }
    },

    //usato per convertire immagine in formato leggibile dal back-end (base64)
    convertToBase64() {
      const reader = new FileReader();
      reader.onload = (e) => {
        this.base64Image = e.target.result;
        // Questo conterrà il prefisso "data:image/jpeg;base64," ecc.
      };
      reader.readAsDataURL(this.selectedFile);
    },

    async fetchMessages(convoID) {
      try {
        const response = await api.getMessages(convoID);
        console.log("(RESPONSE DATA) Ho cercato di fetchare i messages:", response.data);
        this.messages = response.data.conversation.messages || [];
      } catch (error) {
        console.error('Error fetching messages:', error);
        this.messages = [];
      }
    },

    // TODO, Ridichiara sendMessage in App.Vue per chat nuova
    // questo send message da per scontato che la convo esiste
    async sendMessage() {
      if (!this.selectedConvoID) return; // se non seleziono convo...

      try {
        const content = {};
        if (this.newMessage) {
          content.text = this.newMessage;
        }
        if (this.base64Image) {
          // Estrai il tipo MIME dalla stringa data URL
          const mimeType = this.base64Image.substring(this.base64Image.indexOf(":") + 1, this.base64Image.indexOf(";"));
          // Rimuovi il prefisso data URL per ottenere solo la stringa Base64
          const base64Data = this.base64Image.substring(this.base64Image.indexOf(",") + 1);
          content.photo = base64Data;
          content.photoMimeType = mimeType;
        }

        const messageToSend = {
          message: content, // Il backend si aspetta un campo "message" che contiene l'oggetto Content
          recipientUsername: "" // do per scontato che ho l'ID
        };

        await api.sendMessage(this.selectedConvoID, messageToSend);

        // resetto le variabili e ricarico i messaggi
        this.newMessage = '';
        this.selectedFile = null;
        this.base64Image = '';
        this.fetchMessages(this.selectedConvoID);


      }
      catch(error){
        console.log("Error: ",error)
      }
    },
  },
};
</script>

<style scoped>
.chat-messages-container {
  display: flex;
  flex-direction: column;
  height: 100%; /* Fill the height of the chat-area in App.vue */
}

.chatMenu {
  display: flex;
  align-items: center;
  justify-content: space-evenly;
  margin: 0px 10px;
  border-bottom: #ccc 1px solid;
  border-top: #ccc 1px solid;
  
}

.chatMenu button {
  background-color: white;
  padding: 15px 20px 15px 20px;
  font-weight: bold;
  border-radius: 10px; /* Adjust this value to control the roundness */
  border: #ccc 1px solid
}

.conversation-list {
  border-bottom: #ccc 1px solid;
  border-top: #ccc 1px solid;
}
.no-conversation {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #999;
}

.conversation-active {
  display: flex;
  flex-direction: column;
  height: 100%;

}

.chatpfp {
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


.message-list {
  overflow-y: scroll; /* Enable scrolling */
  max-height: 50vh;
  padding: 10px;
  display: flex;
  flex-direction: column;
}


.other-message {
  background-color: #f0f0f0;
  align-self: flex-start;

  text-align: left;
}

.my-message {
  background-color: rgb(235, 232, 255); 

  align-self: flex-end;

}

.author-name {
  font-weight: bold;
  font-size: small;
}

.message-bubble {
  padding: 8px 12px;
  margin-bottom: 5px;
  border-radius: 8px;
  word-break: break-word;
  width: fit-content;
}


.sent-img {
  padding-top: 5px;
  max-height: 200px;
  width: auto;
}

.sent-img:hover{
  border: rgb(112, 92, 111) 1px solid;
  max-height: 500px;
  width: auto;
}

.message-input-form {
  padding: 10px;
    border-top: 1px solid #ccc;
    display: flex;
    gap: 10px;
    justify-content: center;
    flex-direction: column-reverse;
    flex-wrap: wrap;
    align-content: center;
    margin-right: 10px;
    margin-left: 10px;
}



.message-input-form input[type="file"] {
  width: auto;
}

.message-input-form button {
  padding: 8px 15px !important; 
  background-color: #007bff !important; 
  color: white !important; 
  border: none !important; 
  border-radius: 5px !important; 
  cursor: pointer !important; 
}

.message-input-form button:disabled {
  background-color: #ccc !important;
  cursor: not-allowed !important;
}

.chatMenu button:disabled {
  cursor: not-allowed;
}

.chatMenu button:hover {
  background-color: #f0f0f0; 
}


.textInput{

background-color: white;
font-weight: normal;
border-radius: 10px; /* Adjust this value to control the roundness */
border: #ccc 1px solid;
margin: 10px 0px;
text-align: center;
font-weight: normal;

}

.reactionImage{

  width: 20px; /* Adjust the size as needed */
  height: 20px; /* Should be the same as the width */
  border-radius: 50%; /* This makes the element circular */
  overflow: hidden; /* Clips content that goes outside the circle */
  display: inline-block; /* Allows multiple profile pictures to sit on the same line */
  border: 1px solid #ccc;
  object-fit: cover

}

.reactionBubble{

  border: 1px solid #ccc;
  border-radius: 5px;
  display: flex;
  align-items: center;
  padding: 2px;
  background-color: #fae5f0;
  font-style: normal;
  font-size: larger;
}

.messageStats{
  display: flex; 
  justify-content: space-around;
  align-items:center;
  font-style: italic;
  color: #808080;
}

.noParagraph{margin: 0px;}

.timestampAndStatus{
  display: flex;
}

</style>