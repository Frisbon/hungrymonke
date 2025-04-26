<template>
  
  <div v-if="this.userArray.length != 0">
    <h3>Seleziona utenti da aggiungere al gruppo</h3>

    <div class="user-selection-manual">
      <div
        v-for="user in userArray"
        :key="user.username"
        class="userContainer"
        @click="toggleUserSelection(user)"
        :class="{ 'selected': selectedUsernames.includes(user.username) }"
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
          src="https://i.imgur.com/D95gXlb.png"
          alt="Default PFP"
        />
        <h3>{{ user.username }}</h3>
      </div>

    </div>
    <button
      type="button"
      @click="confirmSelection"
      :disabled="selectedUsernames.length === 0"
      style="display:flex; align-self: center;" >
      Conferma Selezione
    </button>
  </div>
  <div v-else>
    <h3>Questo gruppo contiene tutti gli utenti possibili! 🤯</h3>
    <button @click="close">Close</button>
  </div>
  
</template>

<script>
import api from '@/api';

export default {
    name: 'GroupUserList', 

    props: {
      currentUser: String,
      convoID: String
    },

    data() {
      return {
        userArray: [], // Inizializziamo come array vuoto
        selectedUsernames: [], // Ora un array per gli username selezionati
      };
    },

    mounted(){
        this.fetchUsers();
    },

    methods: {

        close(){this.$emit("close-buttons")},

        async fetchUsers(){
            this.userArray = []; // Resetta la lista utenti prima di caricare
            this.selectedUsernames = []; // Resetta la selezione

            try {
                const allUsersResponse = await api.listUsers();
                const allUsers = allUsersResponse.data.Users || [];
                console.log("Lista utenti (raw) ricevuta:", allUsers);

                let existingGroupMembersUsernames = [];
              
                try {
                     const groupInfoResponse = await api.getConvoInfo(this.convoID); 
                     console.log(groupInfoResponse)
                     const groupUsers = groupInfoResponse.data.Group.users || []; 
                     existingGroupMembersUsernames = groupUsers.map(user => user.username);
                     console.log("Membri del gruppo esistente:", existingGroupMembersUsernames);
                } catch (error) {
                    console.error(`Errore nel fetch dei membri del gruppo ${this.convoID}:`, error);
                }

                this.userArray = allUsers.filter(user =>
                    user.username !== this.currentUser && // Filtro 1: utente corrente
                    !existingGroupMembersUsernames.includes(user.username) // Filtro 2: membri esistenti
                );

                console.log("Lista utenti filtrata (disponibili per aggiunta): ", this.userArray);

            } catch (error) {
                console.error("Errore nel fetch o nel filtro degli utenti:", error);
                 this.userArray = []; // Assicurati che sia un array vuoto in caso di errore
            } 
        },

        // Metodo chiamato quando si clicca su un userContainer
        toggleUserSelection(user) {
            const username = user.username;
            const index = this.selectedUsernames.indexOf(username);

            if (index > -1) {
                // Se l'utente è già selezionato, rimuovilo dall'array
                this.selectedUsernames.splice(index, 1);
                console.log("Deselezionato:", username);
            } else {
                // Se l'utente non è selezionato, aggiungilo all'array
                this.selectedUsernames.push(username);
                console.log("Selezionato:", username);
            }
            console.log("Utenti selezionati attuali:", this.selectedUsernames);
        },

        confirmSelection() {
            // Questo metodo viene chiamato quando si clicca il bottone "Conferma Selezione"
            if (this.selectedUsernames.length > 0) {
                console.log("Utenti confermati:", this.selectedUsernames);
                // Emetti un evento al componente genitore con l'array degli username selezionati
                this.$emit('users-selected', this.selectedUsernames);

                this.selectedUsernames = [];

            } else {
                console.warn('Nessun utente selezionato.');
            }
        },

    }
}

</script>

<style scoped>
/* Aggiungi o modifica questi stili per l'effetto di selezione multipla */

/* Stile base per rendere l'area cliccabile */
.userContainer {
    display: flex; /* Usa flexbox per allineare immagine e testo */
    align-items: center; /* Centra verticalmente gli elementi */
    gap: 10px; /* Spazio tra immagine e nome utente */
    cursor: pointer; /* Cambia il cursore per indicare che è cliccabile */
    border: 1px solid transparent; /* Bordo trasparente di default */
    padding: 10px; /* Aggiungi un po' di spazio interno */
    border-radius: 5px; /* Bordi arrotondati (opzionale) */
    transition: all 0.2s ease-in-out; /* Animazione fluida per il cambio di stato */
    margin: 5px; /* Spazio tra i contenitori utente */
    flex-direction: column; /* Metti immagine e nome in colonna */
    text-align: center; /* Centra il testo */
    width: 100px; /* Imposta una larghezza fissa (aggiusta se necessario) */
    justify-content: center; /* Centra contenuto verticalmente */
}

/* Stile per lo stato selezionato */
.userContainer.selected {
  border-color: teal; /* Bordo colorato quando selezionato */
  background-color: #e0f2f7; /* Colore di sfondo leggero (opzionale) */
}

/* Mantieni i tuoi stili esistenti e aggiustali se necessario */
.user-selection-manual {
    display: flex;
    justify-content: center;
    align-content: center;
    flex-wrap: wrap;
    align-items: center;
    overflow-y: auto; /* Permetti lo scroll se ci sono molti utenti */
    max-height: 300px; /* Imposta un'altezza massima con scroll (aggiusta se necessario) */
    padding: 10px;
    border: 1px solid #ccc; /* Aggiungi un bordo al contenitore principale */
    border-radius: 8px;
    margin-bottom: 20px;
}

img.pfp {
    margin: 0px !important; /* Rimuovi margine top/bottom se usi flex gap */
    width: 60px; /* Dimensioni fisse per l'immagine */
    height: 60px;
    border-radius: 50%;
    object-fit: cover;
    border: 1px solid #ccc;
}

h3 {
    margin-top: 5px; /* Spazio sopra il nome utente */
    font-size: 0.9em; /* Dimensione font più piccola per username */
    word-break: break-word; /* Spezza parole lunghe */
}


</style>

