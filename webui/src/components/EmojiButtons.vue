<template>
    <div class="emoji-button-container">
      <button
        v-for="emoji in emojis"
        :key="emoji"
        class="emoji-button"
        @click="emitEmoji(emoji)"
      >
        {{ emoji }}
      </button>
      <button class="emoji-button" @click="closeMenu">X</button>
    </div>
  </template>
  
  <script>
  export default {
    name: 'EmojiButtons',
    data() {
      return {
        emojis: ['😊', '👍', '❤️', '😂', '👏', '🎉', '🤔', '🔥', '💀', '🍆', '👀'],
      };
    },
    methods: {
      emitEmoji(emoji) {
        this.$emit('emojiSelected', emoji);
      },

      closeMenu(){this.$emit("closeReactions")}
    },
  };
  </script>
  
  <style scoped>
  .emoji-button-container {
    display: flex;
    flex-wrap: wrap;
    gap: 5px; 
    max-width: 25%;
    align-content: flex-start;
  }
  
  .emoji-button {
    font-size: 1.2em;
    border: 1px solid #ccc;
    border-radius: 5px;
    background-color: white;
    cursor: pointer;
    line-height: 1;
    width: fit-content;
    height: fit-content;
    padding: 5px !important;
    margin: 0px;
  }
  
  .emoji-button:hover {
    background-color: #f0f0f0;
  }
  </style>