<template>
    
    <div class="typeSelectorDiv" v-if="isTypeToChoose">
        
        <h3>Select the type of conversation:</h3>
        <div class="divButtons">
            <button class="privateButton" @click="privateChatProcedure">Private Chat</button>
            <button class="groupButton" @click="groupChatProcedure">Group Chat</button>
        </div>

        <form @submit.prevent="chooseUsersGroupProcedure" class="groupInputArea" v-if="isGroupChosen">

            <label>Group Name</label>
            <input v-model="selectedName" type="text" required>

            <label>Group Pictre</label>
            <input type="file" @change="handleFileUpload" accept="image/*" />

            <button :disabled="isImageConverting || (selectedFile && base64Image === '')">Seleziona Utenti</button>

        </form>

    </div>

    <PrivateUserList v-if="openPrivateList"
        :current-user="this.currentUser"
        @user-selected="userSelected"
    />
    
    <GroupUserList v-if="openGroupList"
        :current-user="this.currentUser"
        :groupName="this.selectedName"
        @users-selected="usersSelected"
    />

</template>

<script>
import GroupUserList from './GroupUserList.vue';
import PrivateUserList from './PrivateUserList.vue';


 export default {
   name: 'NewChat', 

   components:{PrivateUserList, GroupUserList},

   props: {
     currentUser: String,
   },

   data() {
     return {

        selectedName: "",
        selectedFile: null,

        // questi due da fare emit dopo che ho creato convo per aprire chat subito dopo
        chosenPrivateConvoID: null, 
        chosenGroupConvoID: null,   

        isGroupChosen: false,
        isTypeToChoose: true,
        
        openGroupList: false,
        openPrivateList: false,
        
        base64Image: '',
        isImageConverting: false,
     };
   },

  methods: {

    groupChatProcedure(){this.isGroupChosen = true;},
    privateChatProcedure(){
        this.openPrivateList = true; 
        this.openGroupList = false;
        this.isTypeToChoose = false;
    },

    chooseUsersGroupProcedure(){
        this.isTypeToChoose = false,
        this.openGroupList = true
    },

    //private convo
    userSelected(selectedUser){
        console.log("sono in newchat vue e mi hai passato: ", selectedUser)
        this.$emit("selected-user", selectedUser)
    },
    //group convo
    usersSelected(selectedUserArray){
        console.log("sono in newchat vue")

        let base64Data = '';
        let mimeType = '';

        if (this.base64Image) {
            const parts = this.base64Image.split(','); // Dividi la stringa alla prima virgola
            if (parts.length === 2) {
                base64Data = parts[1];
                // Estraggo il MIME type dal prefisso
                const meta = parts[0].split(';'); // Dividi il prefisso al punto e virgola
                mimeType = meta[0].substring(5); // Rimuovi "data:"
            } else {
                 console.error("Formato stringa Base64 inatteso (nessuna o troppe virgole):", this.base64Image);
                 // Potresti voler impostare base64Data e mimeType a stringa vuota o gestire l'errore
                 base64Data = '';
                 mimeType = '';
            }
        } 



        this.$emit("selected-users", selectedUserArray, this.selectedName, base64Data, mimeType)
    },

    handleFileUpload(event) {
      this.selectedFile = event.target.files[0];
      if (this.selectedFile) {
          this.isImageConverting = true; // Inizia la conversione
          this.convertToBase64();
      } else {
          this.base64Image = '';
          this.isImageConverting = false; // Nessun file, nessuna conversione
      }
    },

    convertToBase64() {
      const reader = new FileReader();
      reader.onload = (e) => {
          this.base64Image = e.target.result;
          this.isImageConverting = false; // Conversione completata
          console.log("Immagine convertita in Base64. Lunghezza:", this.base64Image.length); // Debugging
      };
      reader.onerror = (error) => { // Aggiungi gestione errori
            console.error("Errore durante la conversione dell'immagine:", error);
            this.base64Image = '';
            this.isImageConverting = false;
      };
      reader.readAsDataURL(this.selectedFile);
    },

  },
   mounted() {
      console.log('NewChat mounted.');
   }
}
</script>

<style>


.pfp {
  /* Make the element a square to ensure a perfect circle */
  width: 60px; /* Adjust the size as needed */
  height: 60px; /* Should be the same as the width */
  border-radius: 50%; /* This makes the element circular */
  overflow: hidden; /* Clips content that goes outside the circle */
  display: inline-block; /* Allows multiple profile pictures to sit on the same line */
  margin: 10px 0px;
  border: 1px solid #ccc;
  object-fit: cover

}

.userContainer{
    align-self: center;
    width: 40%;
    display: flex;
    border: 1px solid #ccc;
    margin: 0px 10px 10px;
    border-radius: 10px;
    padding: 0 15px;
    justify-content: space-between;
    align-items: center;
}

.typeSelectorDiv{
    display: flex; 
    justify-content: center;
    flex-direction: column; flex-grow: 1; align-content: center;
}

.divButtons button {
    margin: 0px 20px;
}

</style>