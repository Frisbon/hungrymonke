<template>
  <div class="conversation-list">
    <h2>Hello {{ this.username }}!</h2>

    <ul>
      <li class="convoBubble"
        v-for="(c, index) in this.convertedConvos"
        :key="c.convoid || index"
        @click="selectConversation(c.convoid)"
        
      >
          <p>{{ c.chatName }}</p>
          <p>{{ c.chatStatus +" | "+ c.chatPreview }}</p>
          <p>{{ c.chatTime }}</p>


      </li>
    </ul>
  </div>
</template>

<script>
import api from '../api';

export default {
  name: 'ConversationList',
  
  props: {
    username: String,
  },


  data() {
    return {
      conversations: [],
      convertedConvos: [],
    };
  },

/*

TODO => SORT CONVOS BY TIME + RENDER CONVO PROPERTIES 

*/
  methods: {

    selectConversation(convoID) {
      if (convoID) {
        this.$emit('select-conversation', convoID); // Emit a Select-conversation signal
      }
    },

    //
    async fetchConversations() {
      try {
        const response = await api.getConversations(); // tecnicamente mi ritorna un JSON

        // se dovesser ritorare una stringa, trasformala in json, altrimenti lascia così
        const responseData = typeof response === 'string' ? JSON.parse(response) : response;

        console.log("Parsed response data:", responseData);
        
        // Se esistono dati ricevuti e ricevo effettivamente un array di conversazioni...
        if (responseData && Array.isArray(responseData['User Conversations'])) {
          this.conversations = responseData['User Conversations'];  
          this.convoDataRenderer();
        } else {
          console.error('Unexpected response format:', responseData);
          this.conversations = [];
        }
      
      } catch (error) {
        console.error('Error fetching conversations:', error);
        this.conversations = [];
      }

      
    },

    async convoDataRenderer(){
      console.log("Starting convoDataRenderer");
      this.convertedConvos = []; // Reset the array before populating it

      console.log("Contenuto di this.conversations prima del ciclo:", this.conversations);

      for (let i = 0; i < this.conversations.length; i++) {
        var c = this.conversations[i];
        console.log("Elaborazione conversazione:", i, c);

        // prima imposto i parametri che posso ricavare da convo
        var toRender = {
          convoid: c.convoid,
          chatPreview: c.preview,
          chatTime: c.datelastmessage,
          chatStatus: c.messages?.[c.messages.length - 1]?.status, 

          chatPic: null,
          chatName: null,
          chatPicType: null
        };

        // poi vedo se è di gruppo o privata per il render corretto...
        var response = await api.getConvoInfo(c.convoid);

        console.log(`Risposta per convo ${c.convoid}:`, response);

        if (response?.data?.Group){
          toRender.chatPic   = response.data.Group.groupphoto;
          toRender.chatPicType = response.data.Group.photoMimeType;
          toRender.chatName = response.data.Group.name;
        }
        else if(response?.data?.Private){
          var x = response.data.Private;
          var otherDude = x.firstuser ? (this.username != x.firstuser.username ? x.firstuser : x.seconduser) : x.seconduser;
          toRender.chatPic   = otherDude?.photo;
          toRender.chatPicType = otherDude?.photoMimeType;
          toRender.chatName = otherDude?.username;
        }

        console.log("Oggetto toRender prima del push:", toRender);
        this.convertedConvos.push(toRender);
        console.log("Lunghezza di convertedConvos dopo il push:", this.convertedConvos.length);
      }

      // Ordina l'array 'this.convertedConvos' in base a 'chatTime' (dal più recente al più vecchio)
      this.convertedConvos.sort((a, b) => new Date(b.chatTime).getTime() - new Date(a.chatTime).getTime());
      console.log("Contenuto finale di convertedConvos:", this.convertedConvos);
    },
},
 /* Appena carico la pagina recupera le conversazioni */
 mounted() {
    console.log("ConversationList component mounted!");
    this.fetchConversations();
    console.log("now trying to render chats!");
    
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
</style>