<template>
  <div class="chat-messages-container">
    <!-- If no chat is selected -->
    <div v-if="!selectedConvoID" class="no-conversation">
      <p>Select a conversation to view messages.</p>
    </div>

    <!-- If a chat is selected -->
    <div v-else class="conversation-active">

      <div class="chatMenu">
          
        <img @click="changePfp" class="chatpfp" :src="chatPicSrc" style="margin-top: 0px;">

        <h1>{{ this.selectedConvoRender.chatName }}</h1>

        
        <button :disabled="!this.isGroup" @click="groupChatOptions">
          ...
        </button>

      </div>

      <GroupChatOptions
            v-if="!this.closeButtons"
            :currentUser="this.username"
            :convoID="this.selectedConvoID"
            @closeButtons = "closeGroupChatOptions"
            @groupDataUpdated = "reloadConvo"
            @groupLeft = "groupLeft"
      />

      <div v-if="this.closeButtons">
          

          <p 
            v-if="messages.length == 0"
            style="font-size: large; color: gray; font-weight: lighter;"
          > Digita il tuo primo messaggio! 😄</p>

          <div class="message-list"> 

            <div class= "message-container" v-for="message in messages" :key="message.msgID"
            
            :class="{'other-message': message.author.username !== username,
                        'my-message': message.author.username === username}"
            >
              
              <!-- Fallo comparire  a sx di un messaggio guardando username-->
              <MessageOptions 
              :username = "this.username"
              :selectedMessage = "this.selectedMessage"
              v-if="!this.reactionsOpen && this.messageOptionsOpen && this.selectedMessage && this.selectedMessage.msgid == message.msgid && this.selectedMessage.author.username == this.username"

              @closeMessageOptions = 'closeMessageOptions'
              @openReactions = 'switchReactions'
              @reloadMsgs = 'reloadMessages'
              @reply = 'replyProcedure'
              @handleForward = 'handleForward'

              />
              <EmojiButtons
                @emojiSelected = "reactionWindow"
                v-if="this.reactionsOpen && this.messageOptionsOpen && this.selectedMessage && this.selectedMessage.msgid == message.msgid && this.selectedMessage.author.username == this.username"
                @closeReactions = 'switchReactions'
              />
            
              
              
              
              <div class="message-bubble"
                @click="handleMessageClick(message)"
                :class="{'other-message-bubble': message.author.username !== username,
                        'my-message-bubble': message.author.username === username}"
                >
                  




                  <div v-if="message.author.username !== this.username && this.isGroup">
                    <div style="display: flex;">
                    <strong >{{ message.author.username }}:</strong><br> 
                    <!-- ADD FORWARDED OR REPLIED CONTENT HERE-->
                    <i class="message-bubble-heading" v-if="message.isforwarded" style="color: gray;">↪ (Forwarded)</i>
                    </div>
                  </div>
                  <div v-else>
                    <div style="display: flex;">
                    <!-- ADD FORWARDED OR REPLIED CONTENT HERE-->
                    <i class="message-bubble-heading" v-if="message.isforwarded" style="color: gray;">↪ (Forwarded)</i>
                    </div>
                  </div>
                  <div class="message-reply" v-if="message.replyingto != null">
                      <i >↩ {{ message.replyingto.author.username }}:</i> 
                      <p class="noParagraph" v-if="message.replyingto.content.text != null">{{ message.replyingto.content.text.slice(0,16)}}...</p>
                      <p class="noParagraph" v-else>[📸 Photo]</p>
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

              <!-- Fallo comparire a destra o a sx di un messaggio guardando username-->
              <MessageOptions 
              :username = "this.username"
              :selectedMessage = "this.selectedMessage"
              v-if="!this.reactionsOpen && this.messageOptionsOpen && this.selectedMessage && this.selectedMessage.msgid == message.msgid && this.selectedMessage.author.username != this.username"

              @closeMessageOptions = 'closeMessageOptions'
              @openReactions = 'switchReactions'
              @reloadMsgs = 'reloadMessages'
              @reply = 'replyProcedure'
              @handleForward = 'handleForward'

              
              />
              <EmojiButtons
                @emojiSelected = "reactionWindow"
                v-if="this.reactionsOpen && this.messageOptionsOpen && this.selectedMessage && this.selectedMessage.msgid == message.msgid && this.selectedMessage.author.username != this.username"
                @closeReactions = 'switchReactions'
              />

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
          <div class="reply-div" v-if="replyingMsg != ''">
              <!-- Replying shit + button-->
              <strong v-if="replyingMsg.content.text != null">↪ Replying To: {{replyingMsg.content.text.substring(0,16)}}...</strong> 
              <strong v-else> ↪ Replying To: [📸 Photo] </strong>
              <button @click="resetReply">⨉</button>
          </div>

      </div>


    </div>

    
  </div>
</template>

<script>
import api from '../api';
import defaultPfp from '../assets/blank_pfp.png'; // relative import to the project's assets
import EmojiButtons from './EmojiButtons.vue';
import MessageOptions from './MessageOptions.vue';
import GroupChatOptions from './GroupChatOptions.vue';

export default {
  name: 'ChatMessages',

  components:{MessageOptions, EmojiButtons, GroupChatOptions},

  // typed/defaulted props
  props: {
    selectedConvoID: { type: String, required: true },
    isGroup:        { type: Boolean, default: null },
    username:       { type: String, default: '' },
    selectedConvoRender: { type: Object, default: null },
  },

  data() {
    return {
      messages: [],
      newMessage: '',
      selectedFile: null,
      base64Image: '', // Per memorizzare la foto
      selectedMessage: null,
      selectedEmoji: '',
      reactionsOpen: false,
      messageOptionsOpen: false,
      replyingMsg: "",
      closeButtons: true,
      // added polling fields
      _msgPoll: null,
      _lastMsgId: null,
    };
  },

  computed: {
    // existing computed
    isSubmitDisabled() {
      return !this.newMessage && !this.base64Image;
    },
    // new computed for pfp fallback
    chatPicSrc() {
      const pic  = this.selectedConvoRender?.chatPic;
      const mime = this.selectedConvoRender?.chatPicType;
      if (pic && mime) return `data:${mime};base64,${pic}`;
      return defaultPfp;
    }
  },

  watch: {
    // start polling when convo changes
    selectedConvoID () { this.startMsgPolling() },
    // also refresh silently when username updates
    username () { this.fetchMessagesSilent() },
  },

  mounted() { this.startMsgPolling() },
  beforeUnmount() { this.stopMsgPolling() },

  methods: {
    // BEGIN: polling + silent fetch methods
    async fetchMessagesSilent () {
      if (!this.selectedConvoID) return;
      try {
        const { data } = await api.getMessages(this.selectedConvoID);
        const list = data?.Messages || data?.messages || [];
        const last = list.length ? (list[list.length - 1].MsgID || list[list.length - 1].msgid) : null;
        if (list.length !== this.messages.length || last !== this._lastMsgId) {
          this.messages   = list;
          this._lastMsgId = last;
        }
      } catch (_) {}
    },

    startMsgPolling () {
      this.stopMsgPolling();
      this.fetchMessagesSilent();
      this._msgPoll = setInterval(() => this.fetchMessagesSilent(), 5000);
    },

    stopMsgPolling () {
      if (this._msgPoll) {
        clearInterval(this._msgPoll);
        this._msgPoll = null;
      }
    },
    // END: polling + silent fetch methods

    reloadConvo(){
      this.$emit("reloadConvo")
    },

    groupLeft(){this.$emit("groupLeft")},

    closeGroupChatOptions(){this.closeButtons = true},

    groupChatOptions(){  console.log('Metodo groupChatOptions chiamato!'); this.closeButtons = !this.closeButtons},
    

    resetReply(){this.replyingMsg = ""},

    reactionWindow(emoji){
      console.log("Mi hai passato: "+emoji)
      api.sendReaction(this.selectedMessage.msgid,emoji)
      this.reloadMessages()
    },

    reloadMessages(){
      //  per 2 secondi faccio polling molto veloce,.
      let seconds = 2;
        const timer = setInterval(() => {
          console.log("fetcher messaggi spammato!");
          seconds--;
          this.fetchMessages(this.selectedConvoID)
          if (seconds <= 0) clearInterval(timer);
        }, 500);
      this.closeMessageOptions()
      this.reactionsOpen = false
    },

    switchReactions(){

      console.log("reactionsOpen=="+this.reactionsOpen)
      this.reactionsOpen = !this.reactionsOpen
    },

    closeMessageOptions(){
      this.messageOptionsOpen = false; 
    },

    replyProcedure(){
      this.closeMessageOptions()
      console.log("bruh log")
      this.replyingMsg = this.selectedMessage;
    },
    // NB: permette di inviare una foto per volta, e non in bulk (TODO)
    handleFileUpload(event) {
      this.selectedFile = event.target.files[0];
      if (this.selectedFile) {
        this.convertToBase64();
      } else {
        this.base64Image = '';
      }
    },

    // usato per convertire immagine in formato leggibile dal back-end (base64)
    convertToBase64() {
      const reader = new FileReader();
      reader.onload = (e) => {
        this.base64Image = e.target.result;
        //  Questo conterrà il prefisso "data:image/jpeg;base64," ecc.
      };
      reader.readAsDataURL(this.selectedFile);
    },

    async fetchMessages(convoID) {
      try {
        const response = await api.getMessages(convoID);
        console.log("(RESPONSE DATA) Ho cercato di fetchare i messages:", response.data);
        this.messages = response.data.conversation?.messages || response.data?.messages || [];
      } catch (error) {
        console.error('Error fetching messages:', error);
        this.messages = [];
      }
    },

    //  TODO, Ridichiara sendMessage in App.Vue per chat nuova
    //  questo send message da per scontato che la convo esiste
    async sendMessage() {
      if (!this.selectedConvoID) return; //  se non seleziono convo...

      try {
        const content = {};
        if (this.newMessage) {
          content.text = this.newMessage;
        }
        if (this.base64Image) {
          //  Estrai il tipo MIME dalla stringa data URL
          const mimeType = this.base64Image.substring(this.base64Image.indexOf(":") + 1, this.base64Image.indexOf(";"));
          //  Rimuovi il prefisso data URL per ottenere solo la stringa Base64
          const base64Data = this.base64Image.substring(this.base64Image.indexOf(",") + 1);
          content.photo = base64Data;
          content.photoMimeType = mimeType;
        }

        
        var messageToSend = {};

        if (this.replyingMsg != ''){
            messageToSend = {
              message: content, //  Il backend si aspetta un campo "message" che contiene l'oggetto Content
              recipientUsername: "", //  do per scontato che ho l'ID
              replyingto: this.replyingMsg,
          };
        }else{
          messageToSend = {
            message: content, 
            recipientUsername: "", 
          };
        }

        await api.sendMessage(this.selectedConvoID, messageToSend);

        //  resetto le variabili e ricarico i messaggi
        this.newMessage = '';
        this.selectedFile = null;
        this.base64Image = '';
        // trigger a silent immediate refresh, polling will keep things up-to-date
        this.fetchMessagesSilent();
        this.replyingMsg = "";


      }
      catch(error){
        console.log("Error: ",error)
      }
    },

    handleMessageClick(message){

      console.log("Youve clicked on a message!")
      if (message.author.username == this.username){
        console.log("And it is yours! :)")
      }

      this.selectedMessage = message
      this.messageOptionsOpen = true
      this.reactionsOpen = false

    },

    handleForward(){
      this.closeMessageOptions()
      console.log("faccio l'emit con il messaggio selezionato")
      this.$emit("forwarder", this.selectedMessage)
      this.selectedMessage = null; //  resetto

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
  justify-content: flex-start;
  text-align: left;
}

.my-message {
  justify-content: flex-end;
}

.message-bubble-heading {
  font-size: small;
}

.message-reply{
  display: flex; color: gray; 
  justify-content: space-around; font-size: small;
  flex-direction: column;
}

.reply-div{
    padding-top: 10px;
    border-top: 1px gray solid;
    font-size: small;
    display: flex;
    width: 38%;
    /* align-content: stretch; */
    align-items: baseline;
    justify-content: space-between;
    align-self: center;
}

.message-bubble {
  padding: 8px 12px;
  margin-bottom: 5px;
  border-radius: 8px;
  word-break: break-word;
  width: fit-content;
  height: fit-content;
}

.message-bubble:hover{
  border: 1px solid teal
}

.message-bubble.my-message-bubble{background-color: rgb(235, 232, 255); }
.message-bubble.other-message-bubble{background-color:#f0f0f0; }

.sent-img {
  padding-top: 5px;
  max-height: 200px;
  width: auto;
  display: flex;
  justify-self: center;
}

.sent-img:hover{
  border: rgb(112, 92, 111) 1px solid;
  max-height: 500px;
  width: auto;
}

.message-input-form {
    padding: 10px 10px 0px 10px;
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

.chatMenu button{margin-bottom: 0px !important; } 

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

.message-container{
  display: flex;
}

.message-options{
  display: flex;
    width: fit-content;
    height: fit-content;
    flex-direction: row;
    gap: 5px;
    max-width: 40%;
    flex-wrap: wrap;
}
</style>