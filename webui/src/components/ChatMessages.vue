<template>
    <div class="messages">
      <h2>Messages</h2>
      <div v-if="!selectedConvoID">
        <p>Select a conversation to view messages.</p>
      </div>
      <div v-else>
        <ul>
          <li v-for="message in messages" :key="message.msgID">
            <strong>{{ message.sender }}:</strong> {{ message.content }}
          </li>
        </ul>
        <form @submit.prevent="sendMessage">
          <input v-model="newMessage" placeholder="Type a message..." required />
          <button type="submit">Send</button>
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
          this.messages = response.data.messages || [];
        } catch (error) {
          console.error('Error fetching messages:', error);
          this.messages = [];
        }
      },
      async sendMessage() {
        if (!this.newMessage || !this.selectedConvoID) return;
        try {
          const message = {
            sender: 'You', // Replace with actual user data after login
            content: this.newMessage,
          };
          await api.sendMessage(this.selectedConvoID, message);
          this.newMessage = '';
          this.fetchMessages(this.selectedConvoID); // Refresh messages
        } catch (error) {
          console.error('Error sending message:', error);
        }
      },
    },
  };
  </script>
  
  <style scoped>
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