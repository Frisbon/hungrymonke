<template>
    <button @click="addUsers">Add Users</button>
    <button @click="leaveGroup">Leave Group</button>
    <button @click="renameGroup">Rename Group</button>
    <button @click="setGroupPic">Set Photo</button>


    <div>
        <GroupUserList v-if="showUserList"
            :currentUser="this.current-user"
            :groupName="this.selectedName"
            @users-selected="usersSelected"
        />
    </div>
</template>

<script>
import api from '@/api';
import GroupUserList from './GroupUserList.vue';


export default {
    name: 'GroupChatOptions',

    props: {
      currentUser: String,
      convoID: String,
    },

    data() {
      return {
        groupName: '',
        showUserList: false,
      };
    },
    methods: {
     
        async usersSelected(users){
            this.addUsers(users)
            const response = await api.addUsersToGroup(users) // invia post (AddToGroup)
            console.log("response: ", response.data)
            this.showUserList = false;
            this.$emit("closeButtons");
        },

        addUsers(){
            //apri finestra utenti e carica array
            this.showUserList = true;
        },

        async leaveGroup(){
            this.$emit("") // invia delete (LeaveGroup)
        },

        async renameGroup(){
            this.$emit("") // invia put (SetGroupName)
        },

        async setGroupPic(){
            this.$emit("") // invia put (SetGroupPhoto)
        },





        async fetchGroupName(){
            this.groupName = (await api.getConvoInfo(this.convoID)).data.Name
        }
    },

    mounted(){
        console.log("Montato i group chat options, recupero il nome del gruppo...")
        fetchGroupName()
        console.log(this.groupName)

    }
  };



// passa parametro che differenzia tra private e group se serve
</script>




<style>
</style>