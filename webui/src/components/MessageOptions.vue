<template>

  <div class="message-options"
      :class="{'isMine':isMessageMine}">
          <button @click="forwardButton">Forward</button>
          <button @click="replyButton">Reply</button>
          <button v-if="!isReactedTo" @click="reactButton">React</button>
          <button v-else @click="deleteReactionButton">Delete Reaction</button>
          <button v-if="isMessageMine" @click="deleteMessageButton">Delete</button>
          <button @click="closeButton">Close</button>
       </div>
  
  
  </template>
  
  <script>
  
  export default {
    name: 'MessageOptions',
  
    props: {
      username: String,
      selectedMessage: Object,
    },
  
    data() {
      return {
  
      };
    },
  
    methods: {

      forwardButton(){
      console.log("You've clicked the forwardButton()!")
      },
      replyButton(){
        console.log("You've clicked the replyButton()!")
      },
      deleteReactionButton(){
        console.log("You've clicked the deleteReactionButton()!")
      },
      deleteMessageButton(){
        console.log("You've clicked the deleteMessageButton()!")
      },
      closeButton(){
        console.log("You've clicked the closeButton()!")
        this.$emit("closeMessageOptions");
      },
      reactButton(){
        console.log("You've clicked the reactButton()!")
        this.$emit("openReactions");
      },
    },
  
  
  
  
    computed:{
      isMessageMine() {
        return this.selectedMessage && this.username === this.selectedMessage.author.username;
      },
      isReactedTo() {
        if (!this.selectedMessage || !this.selectedMessage.reactions) {
          return false;
        }
        for (let i = 0; i < this.selectedMessage.reactions.length; i++) {
          const reaction = this.selectedMessage.reactions[i];
          if (reaction.author.username === this.username) {
            return true;
          }
        }
        return false;
      }
    },
  }
  
  
  
  
  </script>
  
  <style>
  .message-options{
      display: flex;
      flex-direction: column;
      justify-content: flex-start;
      padding: 0px 10px;
  }
  
  .message-options.isMine{
      justify-content: flex-end;
  }
  
  button {
      background-color: white !important;
      font-weight: bold !important;
      border-radius: 10px !important;
      border: #ccc 1px solid !important;
      margin-bottom: 10px !important;
      padding: 10px 15px 10px 15px !important;
  }
  
  button:hover{ background-color: #f0f0f0 !important;}
  
  </style>