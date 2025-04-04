<template>
  <div class="conversation-list">


    <div class="logged_user_menu">
      <div class="user-info">

        <img 
          @click="changePfp" class="mainpfp" style="margin-top: 0px; "
          v-if="userPfp != 'https://i.imgur.com/D95gXlb.png' || userPfp != '' " 
          :src="'data:' + this.userPfpType + ';base64,' + this.userPfp">

        <img @click="changePfp" class="mainpfp" v-else :src="'https://i.imgur.com/D95gXlb.png'">
        
        <div style="display: flex; align-items: center; justify-items: flex-start;">

          <h1 style="">Hello,</h1>
          <h1 @click="changeName" class="changeName">{{ this.username }}</h1>

        </div>
        
      </div>

      <div class="logged_menu_buttons">

        <button>New Chat</button>
        <button @click="handleLogout" class="logout-btn">Logout</button>

      </div>
    </div>

    <ul>
      <li class="convoBubble"
        v-for="(c, index) in this.convertedConvos"
        :key="c.convoid || index"
        @click="selectConversation(c.convoid)"
        >

          <div class="convoBubbleLeft">
            
            <div style="display: flex; align-items: center; justify-items: flex-start;">

              <img class="pfp" v-if="c.chatPic != null && c.chatPic != 'https://i.imgur.com/D95gXlb.png' && c.chatPic != ''"    
              :src="'data:' + c.chatPicType + ';base64,' + c.chatPic">                 
              
              <img class="pfp" v-else :src="'https://i.imgur.com/D95gXlb.png'">
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
            <p class="notificationDot" v-if="c.chatStatus == 'delivered' && c.lastSender != this.username"></p>
            <br v-else>

            <p>{{ c.chatTime.substring(11, 16) }}</p>
          
          </div>

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
    userPfp: String,
    userPfpType: String,
  },


  data() {
    return {
      conversations: [],
      convertedConvos: [],
    };
  },

  methods: {

    handleLogout() {
      this.$emit('logout');
    },

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
          lastSender: c.messages?.[c.messages.length - 1]?.author?.username,

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
        else if(response?.data?.PrivateConvo){
          var x = response.data.PrivateConvo;
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


    //todo
    changeName(){
      console.log("Trying to change name...")
    },

    changePfp(){
      console.log("Trying to change the profile picture...")
    }

},
 /* Appena carico la pagina recupera le conversazioni */
 mounted() {
    console.log("ConversationList component mounted!");
    this.fetchConversations();
    
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
  background-color: #f0f0f0;
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
  padding-bottom: 10px;
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
}

.logged_menu_buttons button {

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
