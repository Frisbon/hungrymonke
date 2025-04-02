<template>
    <div class="messages">

      <h2>Messages</h2>

      <!-- Se ancora non ho selezionato la chat-->
      <div v-if="!selectedConvoID"> 
        <p>Select a conversation to view messages.</p>
      </div>
      
      <!-- Se la chat è selezionata-->
      <div v-else>
        <ul>
          <li class="messageBubble" v-for="message in messages" :key="message.msgID">
            <strong>{{ message.author.username }}:</strong> {{ message.content.text }}
            <br v-if="message.content.photo && message.content.text">
            <img class="sent-img" v-if="message.content.photo" :src="'data:' + message.content.phototype + ';base64,' + message.content.photo" alt="Immagine allegata">
          </li>
        </ul>

        <form @submit.prevent="sendMessage">
         
          
        <form @submit.prevent="sendMessage">
         
          <input v-model="newMessage" placeholder="Type a message..." />
          <input type="file" @change="handleFileUpload" accept="image/*">
          <button type="submit" :disabled="isSubmitDisabled"> <strong>→</strong></button>

       </form>

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

    //TODO: similarmente posso calcolare se il messaggio è da renderizzare a sinistra o a destra con il css a seconda del current user.

    },

    watch: {
      /*
      ogni volta che selectedConvoID (nei props) cambia:
        - newConvoID sarà il nuovo valore
        - eseguo la funzione scritta sotto (recupero i nuovi messaggi)
      */
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

  .sent-img{
    max-height: 150px; /* Imposta l'altezza massima desiderata (puoi cambiarla) */
    width: auto; /* La larghezza si adatterà proporzionalmente */
  }

  .messages {
    width: 70%;
    padding: 10px;
  }
  ul {
    list-style: none;
    padding: 0;
  }
  li {
    padding: 5px 0;
  }
  form {
    margin-top: 10px;
  }
  input {
    padding: 5px;
    margin-right: 5px;
  }
  button {
    padding: 5px 10px;
  }
</style>