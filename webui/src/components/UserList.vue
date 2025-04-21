<template>
    
    <div class="userList">

        <div 
             @click="chooseUser"
             class="userContainer"
             v-for="user in this.userArray"
             v-key="user.username">
            
        <img 
        class="pfp" style="margin-top: 0px;"
        v-if="user.photo != null" 
        :src="'data:' + user.photoMimeType+';base64,'+ user.photo">

        <img class="pfp" v-else :src="'https://i.imgur.com/D95gXlb.png'">

        <h3>{{ user.username }}</h3>

        </div>

    </div>

</template>

<script>
import api from '@/api';

 export default {
    name: 'UserList',
  
    props: {
      currentUser: String,
    },
  
    data() {
      return {
        userArray: Array,
      };
    },
  
    methods: {

        chooseUser(x){
            console.log("You've chosen the user: ", x);
            this.$emit('selectedUserList', x)
        }

    },

    mounted:{

        fetchUsers(){
            this.userArray = api.listUsers()
        }

    }
}

</script>

<style>


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

</style>