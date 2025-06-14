<template>
    <button @click="addUsers">Add Users</button>
    <button @click="leaveGroup">Leave Group</button>
    <button @click="handleRenaming">Rename Group</button>
    <button @click="setGroupPic">Set Photo</button>
    <input type="file" @change="handleFileUpload" accept="image/*" style="display: none;" ref="fileInput">


    <div>
        <GroupUserList v-if="showUserList"
            :currentUser="this.current-user"
            :convoID="this.convoID"
            @users-selected="usersSelected"
            @close-buttons="closeButtons"
        />
    </div>

    <div class="form" v-if="showNameInput">
      <form @submit.prevent="renameGroup">
        <div>
          <label>Insert a new group name:</label>
          <br>
          <input v-model="newName" required type="text"/>
        </div>
        <button type="submit">Set New Name</button>
      </form>
      <button @click="handleRenaming">Cancel</button>
    </div>



</template>

<script>
import api from '@/api';
import GroupUserList from './GroupUserList.vue';


export default {
    name: 'GroupChatOptions',

    components: { 
        GroupUserList 
    },

    props: {
      currentUser: String,
      convoID: String,
    },

    data() {
      return {       
        showUserList: false,
        newName: '',

        showNameInput: false, //  Aggiungi questa data property per mostrare/nascondere il form rinomina
        selectedFile: null, //  Aggiungi data property per il file selezionato
      };
    },
    methods: {
     
        closeButtons(){this.$emit("closeButtons")},

        async usersSelected(users){
            this.addUsers(users)
            const response = await api.addUsersToGroup(users, this.convoID) //  invia post (AddToGroup)
            console.log("response: ", response.data)
            this.showUserList = false;
            this.$emit("closeButtons")
        },

        addUsers(){
            // apri finestra utenti e carica array
            this.showUserList = !this.showUserList;
        },

        async leaveGroup(){
            const response = await api.leaveGroup(this.convoID)
            console.log("response: ", response)
            this.$emit("closeButtons")
            this.$emit("groupLeft");
        },

        handleRenaming(){
            this.showNameInput = !this.showNameInput
            this.newName = '';
        },

        async renameGroup(){
            const response = await api.setGroupName(this.convoID, this.newName)
            console.log("response: ", response.data)
            this.$emit("closeButtons")
            this.$emit("groupDataUpdated");
        },

        setGroupPic() {
            console.log("Tasto Set Photo cliccato, apro selezione file");
            //  Trova l'input file usando il suo ref
            const fileInput = this.$refs.fileInput;
            //  Simula un click sull'input file
            if (fileInput) {
                fileInput.click();
            }
        },


        handleFileUpload(event) {
        this.selectedFile = event.target.files[0];
        if (this.selectedFile) {
            let response = api.setGroupPhoto(this.convoID, this.selectedFile)
            console.log("fileUploadLog: ", response.data)
            this.$emit("groupDataUpdated")
        } 
        this.selectedFile = null;
        },


    },

    mounted(){
        console.log("Montato i group chat options")
       
    },

  };



//  passa parametro che differenzia tra private e group se serve
</script>




<style>
</style>