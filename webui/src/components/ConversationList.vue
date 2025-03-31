<template>
  <div class="conversation-list">
    <h2>Conversations</h2>
    <ul>
      <li
        v-for="(convo, index) in validConversations"
        :key="convo.convoid || index"
        @click="selectConversation(convo.convoid)"
      >
        {{ convo.convoid }} 
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
  computed: {
    validConversations() {
      return this.conversations.filter(convo => typeof convo === 'object' && convo.convoid);
    }
  },
  methods: {
    selectConversation(convoID) {
      if (convoID) {
        this.$emit('select-conversation', convoID);
      }
    },
    async fetchConversations() {
      try {
        const response = await api.getConversations();
        const responseData = typeof response === 'string' ? JSON.parse(response) : response;
        console.log("Parsed response data:", responseData);
        
        if (responseData && Array.isArray(responseData['User Conversations'])) {
          this.conversations = responseData['User Conversations'];
        } else {
          console.error('Unexpected response format:', responseData);
          this.conversations = [];
        }
      } catch (error) {
        console.error('Error fetching conversations:', error);
        this.conversations = [];
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