<template>
  <h3>Choose a user to start a conversation with:</h3>
  
  <div class="user-selection-manual">
    <div
      v-for="user in userArray"
      :key="user.username"
      class="userContainer"
      @click="selectUser(user)"
      :class="{ 'selected': this.selectedUsername && this.selectedUsername === user.username }"
    >
      
            <img
              class="pfp"
              style="margin-top: 0px;"
              v-if="user.photo != null && user.photo !== ''"
              :src="'data:' + user.photoMimeType + ';base64,' + user.photo"
              alt="Profile Picture" />
            <img
              class="pfp"
              v-else
              src="https:// i.imgur.com/D95gXlb.png"
              alt="Default PFP"
            />
              <h3>{{ user.username }}</h3>  
    </div>
            
  </div>
  <button
              type="button"
              @click="confirmSelection"
              :disabled="!this.selectedUsername"
              style="display:flex; align-self: center;" >
              Conferma Utente
            </button>
</template>

<script>
import api from '@/api';

export default {
    name: 'PrivateUserList',
  
    props: {
      currentUser: String,
    },
  
    data() {
      return {
        userArray: Array,
        selectedUsername: null,
      };
    },

    mounted(){this.fetchUsers()},

    methods: {

        async fetchUsers(){
            this.userArray = (await api.listUsers()).data.Users
            console.log("MI È ARRIVATO: ", this.userArray)
        },

        selectUser(user) {
          this.selectedUsername = user.username; //  Aggiorna la variabile con l'utente cliccato
          console.log("Utente cliccato:", this.selectedUsername)
        },

        confirmSelection() {
          //  Questo metodo viene chiamato quando il form viene sottomesso (o il bottone cliccato se @click)
          if (this.selectedUsername) {
            console.log("Utente selezionato:", this.selectedUsername);
            //  Emetti un evento al componente genitore con lo username selezionato
            this.$emit('user-selected', this.selectedUsername);

          }
        },

       
    }
}

</script>

<style>

img.pfp {margin: 10px 0px !important;}



.user-selection-manual {
    display: flex;
    justify-content: center;
    align-content: center;
    flex-wrap: wrap;
    align-items: center;
    overflow-y: hidden;
}

.selected {
  border-color: teal; /* Bordo colorato quando selezionato */
  background-color: #e0f2f7; /* Colore di sfondo leggero (opzionale) */
}

.userContainer {width: 25%}

.userContainer:hover{ background-color: #f0f0f0}
</style>