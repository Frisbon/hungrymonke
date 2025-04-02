<template>
  <div class="chat-messages-container">
    <!-- Se ancora non ho selezionato la chat-->
    <div v-if="!selectedConvoID" class="no-conversation">
      <p>Select a conversation to view messages.</p>
    </div>

    <!-- Se la chat è selezionata-->
    <div v-else class="conversation-active">
      <ul class="message-list">
        <li class="messageBubble" v-for="message in messages" :key="message.msgID">
          <strong>{{ message.author.username }}:</strong> {{ message.content.text }}
          <br v-if="message.content.photo && message.content.text">
          <img class="sent-img" v-if="message.content.photo" :src="'data:' + message.content.photoMimeType + ';base64,' + message.content.photo" alt="Immagine allegata">
        </li>
      </ul>

      <form @submit.prevent="sendMessage" class="message-input-form">
        <input v-model="newMessage" placeholder="Type a message..." />
        <input type="file" @change="handleFileUpload" accept="image/*">
        <button type="submit" :disabled="isSubmitDisabled"> <strong>→</strong></button>
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

.message-list {
  flex-grow: 1; /* Make the message list take up available vertical space */
  overflow-y: auto; /* Enable scrolling for messages */
  padding: 10px;
}

.messageBubble {
  padding: 8px 12px;
  margin-bottom: 5px;
  background-color: #f0f0f0;
  border-radius: 8px;
  word-break: break-word;
}

.sent-img{
  max-height: 150px; /* Imposta l'altezza massima desiderata (puoi cambiarla) */
  width: auto; /* La larghezza si adatterà proporzionalmente */
}

.message-input-form {
  padding: 10px;
  border-top: 1px solid #ccc;
  display: flex; /* Arrange input and button horizontally */
  gap: 10px; /* Space between input and button */
}

.message-input-form input[type="text"] {
  flex-grow: 1; /* Make the input take up remaining width */
  padding: 8px;
  border: 1px solid #ddd;
  border-radius: 5px;
}

.message-input-form input[type="file"] {
  /* Basic styling, adjust as needed */
  width: auto;
}

.message-input-form button {
  padding: 8px 15px;
  background-color: #007bff;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.message-input-form button:disabled {
  background-color: #ccc;
  cursor: not-allowed;
}
</style>