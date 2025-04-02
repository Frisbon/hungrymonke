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
          <li v-for="message in messages" :key="message.msgID">
            <strong>{{ message.author.username }}:</strong> {{ message.content.text }}
            <br v-if="message.content.photo && message.content.text">
            <img class="sent-img" v-if="message.content.photo" :src="'data:' + message.content.phototype + ';base64,' + message.content.photo" alt="Immagine allegata">
          </li>
        </ul>

        <form @submit.prevent="sendMessage">
         
          <input v-model="newMessage" placeholder="Type a message..." required />
          <button type="submit"> <strong>→</strong></button>

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
      };
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
        if (!this.newMessage || !this.selectedConvoID) return; // se non seleziono convo oppure non invio messaggio...

        try {
          const content = { text: this.newMessage };

          const messageToSend = {
            message: content, // Il backend si aspetta un campo "message" che contiene l'oggetto Content
            recipientUsername: "" // do per scontato che ho l'ID
          };

          await api.sendMessage(this.selectedConvoID, messageToSend);
          
          // resetto il messaggio e ricarico quelli nuovi
          this.newMessage = '';
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