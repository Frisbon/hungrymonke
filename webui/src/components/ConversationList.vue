<template>
    <div class="conversation-list">
      <h2>Conversations</h2>
      <ul>
        <li v-for="convo in conversations" :key="convo.ID" @click="selectConversation(convo.ID)">
          {{ convo.name }}
        </li>
      </ul>
    </div>
  </template>
  
  <script>
  import api from '../api';
  
  export default {
    name: 'ConversationList',
    data() {
      return {
        conversations: [],
      };
    },
    methods: {
      selectConversation(convoID) {
        this.$emit('select-conversation', convoID);
      },
      async fetchConversations() {
        try {
          const response = await api.getConversations();
          this.conversations = response.data;
        } catch (error) {
          console.error('Error fetching conversations:', error);
        }
      },
    },
    mounted() {
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
  </style>