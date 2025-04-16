<template>

<div class="message-options" 
    :class="{'isMine':isMessageMine}">
            <button>Forward</button>
            <button>Reply</button>
            <button v-if="!isReactedTo">React</button>
            <button v-else>Delete Reaction</button>
            <button v-if="isMessageMine">Delete</button>
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
}

.message-options.isMine{
    justify-content: flex-end;
}


</style>